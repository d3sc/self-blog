package config

import (
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL, SessionSecret, Port, MigrationsPath string
	CookieSecure                                     bool
	TrustedOrigins                                   []string
}

func Load() (Config, error) {
	c := Config{DatabaseURL: os.Getenv("DATABASE_URL"), SessionSecret: os.Getenv("SESSION_SECRET"), Port: os.Getenv("PORT"), MigrationsPath: os.Getenv("MIGRATIONS_PATH")}
	if c.Port == "" {
		c.Port = "8080"
	}
	if c.MigrationsPath == "" {
		c.MigrationsPath = "file://migrations"
	}
	if c.DatabaseURL == "" || len(c.SessionSecret) < 32 {
		return c, errors.New("DATABASE_URL and SESSION_SECRET (32+ characters) are required")
	}
	if b, e := base64.StdEncoding.DecodeString(c.SessionSecret); e == nil && len(b) >= 32 {
		c.SessionSecret = string(b)
	}
	c.CookieSecure = true
	if v := os.Getenv("COOKIE_SECURE"); v != "" {
		c.CookieSecure, _ = strconv.ParseBool(v)
	}
	for _, v := range strings.Split(os.Getenv("TRUSTED_ORIGINS"), ",") {
		if strings.TrimSpace(v) != "" {
			c.TrustedOrigins = append(c.TrustedOrigins, strings.TrimSpace(v))
		}
	}
	return c, nil
}
