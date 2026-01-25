package config

import "os"

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	R2AccountID    string
	R2AccessKey    string
	R2SecretKey    string
	R2Bucket       string
	R2PublicDomain string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/lovelion?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-key"),
		Port:           getEnv("PORT", "8080"),
		R2AccountID:    getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKey:    getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretKey:    getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2Bucket:       getEnv("R2_BUCKET_NAME", ""),
		R2PublicDomain: getEnv("R2_PUBLIC_DOMAIN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
