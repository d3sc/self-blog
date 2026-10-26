package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/personal-blog/cms/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type Service struct {
	DB             *pgxpool.Pool
	Codec          *securecookie.SecureCookie
	Limiter        interface{ Allow(string) bool }
	CookieSecure   bool
	TrustedOrigins []string
}

func New(db *pgxpool.Pool, secret string, secure bool, lim interface{ Allow(string) bool }, origins []string) *Service {
	key := sha256.Sum256([]byte(secret))
	return &Service{DB: db, Codec: securecookie.New(key[:], key[:]), Limiter: lim, CookieSecure: secure, TrustedOrigins: origins}
}
func (s *Service) originAllowed(c *gin.Context) bool {
	o := c.GetHeader("Origin")
	if o == "" {
		return true
	}
	for _, v := range s.TrustedOrigins {
		if o == v {
			return true
		}
	}
	return false
}

func (s *Service) Login(c *gin.Context) {
	if !s.originAllowed(c) {
		apiErr(c, 403, "origin_denied", "origin is not trusted")
		return
	}
	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if c.ShouldBindJSON(&in) != nil || in.Username == "" || in.Password == "" {
		apiErr(c, 400, "validation_error", "username and password are required")
		return
	}
	if !s.Limiter.Allow(c.ClientIP()) {
		apiErr(c, 429, "rate_limited", "too many login attempts")
		return
	}
	var a models.Admin
	var hash string
	if e := s.DB.QueryRow(c, "SELECT id, username, password_hash FROM admins WHERE username=$1", in.Username).Scan(&a.ID, &a.Username, &hash); e != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Password)) != nil {
		apiErr(c, 401, "invalid_credentials", "invalid credentials")
		return
	}
	sid := make([]byte, 32)
	_, _ = rand.Read(sid)
	val, err := s.Codec.Encode("session", a.ID)
	if err != nil { apiErr(c, 500, "internal_error", "internal server error"); return }
	csrf := make([]byte, 32)
	_, _ = rand.Read(csrf)
	value := base64.RawURLEncoding.EncodeToString(csrf)
	csrfValue := value
	http.SetCookie(c.Writer, &http.Cookie{Name: "blog_session", Value: val, Path: "/", HttpOnly: true, Secure: s.CookieSecure, SameSite: http.SameSiteLaxMode, MaxAge: 86400})
	http.SetCookie(c.Writer, &http.Cookie{Name: "blog_csrf", Value: csrfValue, Path: "/", Secure: s.CookieSecure, SameSite: http.SameSiteLaxMode, MaxAge: 86400})
	c.JSON(200, gin.H{"admin": a})
}
func apiErr(c *gin.Context, s int, code, msg string) {
	c.AbortWithStatusJSON(s, gin.H{"error": gin.H{"code": code, "message": msg}})
}
func (s *Service) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		co, e := c.Request.Cookie("blog_session")
		if e != nil {
			c.Next()
			return
		}
		var id string
		if s.Codec.Decode("session", co.Value, &id) != nil {
			c.Next()
			return
		}
		var a models.Admin
		if s.DB.QueryRow(c, "SELECT id, username FROM admins WHERE id=$1", id).Scan(&a.ID, &a.Username) == nil {
			c.Set("admin", a)
		}
		c.Next()
	}
}
func Require(c *gin.Context) {
	if _, ok := c.Get("admin"); !ok {
		apiErr(c, 401, "unauthorized", "authentication required")
		return
	}
	c.Next()
}
func (s *Service) CSRF(c *gin.Context) {
	if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
		c.Next()
		return
	}
	if o := c.GetHeader("Origin"); o != "" && !s.originAllowed(c) {
		apiErr(c, 403, "origin_denied", "origin is not trusted")
		return
	}
	a, e := c.Request.Cookie("blog_csrf")
	if e != nil || a.Value == "" || a.Value != c.GetHeader("X-CSRF-Token") {
		apiErr(c, 403, "csrf_failed", "CSRF token required")
		return
	}
	c.Next()
}
func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{Name: "blog_session", Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	c.Status(204)
}
func Me(c *gin.Context) { a, _ := c.Get("admin"); c.JSON(200, gin.H{"admin": a}) }
func HashPassword(p string) (string, error) {
	b, e := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b), e
}
func Expired(t time.Time) bool { return time.Now().After(t) }
