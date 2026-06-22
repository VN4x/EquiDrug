package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	APIHost     string
	APIPort     string
	DatabaseURL string
	CORSOrigins []string
}

func Load() Config {
	origins := strings.Split(getEnv("CORS_ORIGINS", "http://localhost:16126"), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return Config{
		APIHost:     getEnv("API_HOST", "0.0.0.0"),
		APIPort:     getEnv("API_PORT", "16125"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://equidrug:equidrug@localhost:16127/equidrug?sslmode=disable"),
		CORSOrigins: origins,
	}
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.APIHost, c.APIPort)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
