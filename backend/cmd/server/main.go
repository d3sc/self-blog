package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/personal-blog/cms/backend/internal/auth"
	"github.com/personal-blog/cms/backend/internal/config"
	"github.com/personal-blog/cms/backend/internal/posts"
	"github.com/personal-blog/cms/backend/internal/ratelimit"
	"log"
	"net/http"
	"os"
)

func migrateDatabase(cfg config.Config) {
	m, err := migrate.New(cfg.MigrationsPath, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		migrateDatabase(cfg)
		return
	}
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	migrateDatabase(cfg)
	lim := ratelimit.New(10, 15*60e9)
	a := auth.New(db, cfg.SessionSecret, cfg.CookieSecure, lim, cfg.TrustedOrigins)
	p := &posts.Store{DB: db}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), a.Middleware())
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	v := r.Group("/api/v1")
	v.POST("/auth/login", a.Login)
	v.POST("/auth/logout", auth.Require, a.CSRF, auth.Logout)
	v.GET("/auth/me", auth.Require, auth.Me)
	v.GET("/posts", p.PublicList)
	v.GET("/posts/:slug", p.GetPublic)
	ad := v.Group("/admin", auth.Require)
	ad.Use(a.CSRF)
	ad.GET("/posts", p.AdminList)
	ad.GET("/posts/:id", p.GetAdmin)
	ad.POST("/posts", p.Save)
	ad.PUT("/posts/:id", p.Update)
	ad.DELETE("/posts/:id", p.Delete)
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
