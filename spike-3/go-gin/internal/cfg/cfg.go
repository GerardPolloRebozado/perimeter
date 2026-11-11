package cfg

import (
	"fmt"
	"os"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RedisAddr     string
	OIDCIssuerURL string
	ServiceName   string
	Version       string
}

func Load() (*Config, error) {
	c := &Config{
		Port:          getenv("PORT", "8080"),
		DatabaseURL:   getenv("DATABASE_URL", "postgres://yugabyte:yugabyte@localhost:5433/yugabyte?sslmode=disable"),
		RedisAddr:     getenv("REDIS_ADDR", "localhost:6379"),
		OIDCIssuerURL: os.Getenv("OIDC_ISSUER_URL"),
		ServiceName:   "go-gin",
		Version:       "0.0.1",
	}
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	return c, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
