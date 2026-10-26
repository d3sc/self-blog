package posts

import (
	"context"
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/personal-blog/cms/backend/internal/models"
	"github.com/personal-blog/cms/backend/internal/slug"
	"math"
	"strings"
	"time"
)

type Store struct{ DB *pgxpool.Pool }

func scan(r interface{ Scan(...any) error }) (models.Post, error) {
	var p models.Post
	e := r.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.PublishedAt)
	return p, e
}
func (s *Store) list(c context.Context, public bool, page, limit int, status string) (models.Page, error) {
	where := ""
	args := []any{}
	if public {
		where = "WHERE status='published'"
	} else if status != "" && status != "all" {
		where = "WHERE status=$1"
		args = append(args, status)
	}
	var total int
	if e := s.DB.QueryRow(c, "SELECT count(*) FROM posts "+where, args...).Scan(&total); e != nil {
		return models.Page{}, e
	}
	args = append(args, limit, (page-1)*limit)
	order := "updated_at DESC"
	if public {
		order = "published_at DESC"
	}
	rows, e := s.DB.Query(c, "SELECT id,title,slug,excerpt,content,status,created_at,updated_at,published_at FROM posts "+where+" ORDER BY "+order+" LIMIT $"+fmt.Sprint(len(args)-1)+" OFFSET $"+fmt.Sprint(len(args)), args...)
	if e != nil {
		return models.Page{}, e
	}
	defer rows.Close()
	items := []models.Post{}
	for rows.Next() {
		p, e := scan(rows)
		if e != nil {
			return models.Page{}, e
		}
		items = append(items, p)
	}
	return models.Page{Items: items, Pagination: models.Pagination{Page: page, Limit: limit, Total: total, TotalPages: int(math.Ceil(float64(total) / float64(limit)))}}, rows.Err()
}
func page(c *gin.Context, def int) (int, int) {
	p := 1
	l := def
	fmt.Sscanf(c.Query("page"), "%d", &p)
	fmt.Sscanf(c.Query("limit"), "%d", &l)
	if p < 1 {
		p = 1
	}
	if l < 1 {
		l = def
	}
	if l > 100 {
		l = 100
	}
	return p, l
}
func (s *Store) PublicList(c *gin.Context) {
	p, l := page(c, 10)
	v, e := s.list(c, true, p, l, "")
	if e != nil {
		err(c, e)
		return
	}
	c.JSON(200, v)
}
func (s *Store) AdminList(c *gin.Context) {
	p, l := page(c, 20)
	v, e := s.list(c, false, p, l, c.Query("status"))
	if e != nil {
		err(c, e)
		return
	}
	c.JSON(200, v)
}
func (s *Store) GetPublic(c *gin.Context) {
	p, e := scan(s.DB.QueryRow(c, "SELECT id,title,slug,excerpt,content,status,created_at,updated_at,published_at FROM posts WHERE slug=$1 AND status='published'", c.Param("slug")))
	if e != nil {
		c.Status(404)
		return
	}
	c.JSON(200, p)
}
func (s *Store) GetAdmin(c *gin.Context) {
	p, e := scan(s.DB.QueryRow(c, "SELECT id,title,slug,excerpt,content,status,created_at,updated_at,published_at FROM posts WHERE id=$1", c.Param("id")))
	if e != nil {
		c.Status(404)
		return
	}
	c.JSON(200, p)
}
func input(c *gin.Context) (models.PostInput, error) {
	var in models.PostInput
	e := c.ShouldBindJSON(&in)
	if e == nil {
		in.Title = strings.TrimSpace(in.Title)
		if in.Slug != nil {
			v := strings.TrimSpace(*in.Slug)
			in.Slug = &v
		}
		if in.Title == "" || in.Content == "" || (in.Status != "draft" && in.Status != "published") {
			e = fmt.Errorf("invalid post")
		}
	}
	return in, e
}
func (s *Store) Save(c *gin.Context) {
	in, e := input(c)
	if e != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "validation_error", "message": "invalid post"}})
		return
	}
	sl := ""
	if in.Slug != nil {
		sl = *in.Slug
	}
	if sl == "" {
		sl = slug.Make(in.Title)
	}
	if sl == "" {
		c.JSON(400, gin.H{"error": gin.H{"code": "validation_error", "message": "slug cannot be empty"}})
		return
	}
	now := time.Now()
	id := uuid.New()
	_, e = s.DB.Exec(c, "INSERT INTO posts(id,title,slug,excerpt,content,status,created_at,updated_at,published_at) VALUES($1,$2,$3,$4,$5,$6,$7::timestamptz,$7::timestamptz,CASE WHEN $6='published' THEN $7::timestamptz ELSE NULL END)", id, in.Title, sl, in.Excerpt, in.Content, in.Status, now)
	if e != nil {
		log.Printf("post insert failed: %v", e)
		if strings.Contains(e.Error(), "duplicate") || strings.Contains(e.Error(), "unique") {
			c.JSON(409, gin.H{"error": gin.H{"code": "slug_conflict", "message": "slug already exists"}})
			return
		}
		err(c, e)
		return
	}
	s.GetAdminByID(c, id.String(), 201)
}
func (s *Store) Update(c *gin.Context) {
	in, e := input(c)
	if e != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "validation_error", "message": "invalid post"}})
		return
	}
	sl := ""
	if in.Slug != nil {
		sl = *in.Slug
	}
	if sl == "" {
		sl = slug.Make(in.Title)
	}
	_, e = s.DB.Exec(c, "UPDATE posts SET title=$1,slug=$2,excerpt=$3,content=$4,status=$5,updated_at=now(),published_at=CASE WHEN $5='published' THEN COALESCE(published_at,now()) ELSE NULL END WHERE id=$6", in.Title, sl, in.Excerpt, in.Content, in.Status, c.Param("id"))
	if e != nil {
		if strings.Contains(e.Error(), "duplicate") || strings.Contains(e.Error(), "unique") { c.JSON(409, gin.H{"error": gin.H{"code":"slug_conflict","message":"slug already exists"}}); return }
		err(c, e)
		return
	}
	s.GetAdminByID(c, c.Param("id"), 200)
}
func (s *Store) GetAdminByID(c *gin.Context, id string, status int) {
	p, e := scan(s.DB.QueryRow(c, "SELECT id,title,slug,excerpt,content,status,created_at,updated_at,published_at FROM posts WHERE id=$1", id))
	if e != nil {
		c.Status(404)
		return
	}
	c.JSON(status, p)
}
func (s *Store) Delete(c *gin.Context) {
	r, e := s.DB.Exec(c, "DELETE FROM posts WHERE id=$1", c.Param("id"))
	if e != nil {
		err(c, e)
		return
	}
	if r.RowsAffected() == 0 {
		c.Status(404)
		return
	}
	c.Status(204)
}
func err(c *gin.Context, e error) {
	c.JSON(500, gin.H{"error": gin.H{"code": "internal_error", "message": "internal server error"}})
}
