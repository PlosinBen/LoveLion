package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	CORSOrigins    []string
	R2AccountID    string
	R2AccessKey    string
	R2SecretKey    string
	R2Bucket       string
	R2PublicDomain string
}

func Load() *Config {
	isRelease := os.Getenv("GIN_MODE") == "release"

	cfg := &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/lovelion?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-key"),
		Port:           getEnv("PORT", "8080"),
		CORSOrigins:    parseCORSOrigins(os.Getenv("CORS_ORIGINS")),
		R2AccountID:    getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKey:    getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretKey:    getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2Bucket:       getEnv("R2_BUCKET_NAME", ""),
		R2PublicDomain: getEnv("R2_PUBLIC_DOMAIN", ""),
	}

	if isRelease {
		if cfg.JWTSecret == "dev-secret-key" {
			log.Fatal("JWT_SECRET must be set in production (GIN_MODE=release)")
		}
		if cfg.DatabaseURL == "" {
			log.Fatal("DATABASE_URL must be set in production (GIN_MODE=release)")
		}
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseCORSOrigins parses comma-separated origins, e.g. "https://lovelion.app,https://www.lovelion.app"
func parseCORSOrigins(raw string) []string {
	if raw == "" {
		return nil
	}
	var origins []string
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}
