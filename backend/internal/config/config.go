package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	JWTExpiry      time.Duration
	Port           string
	CORSOrigins    []string
	R2AccountID    string
	R2AccessKey    string
	R2SecretKey    string
	R2Bucket       string
	R2PublicDomain string

	AuthRateLimit int

	// AI receipt extraction
	GeminiAPIKey           string
	GeminiModel            string
	GeminiBaseURL          string
	ReceiptExtractEnabled  bool
	ReceiptRateLimitPerDay int
}

func Load() *Config {
	isRelease := os.Getenv("GIN_MODE") == "release"

	cfg := &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/lovelion?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-key"),
		JWTExpiry:      parseDurationDays(getEnv("JWT_EXPIRY_DAYS", "30")),
		Port:           getEnv("PORT", "8080"),
		CORSOrigins:    parseCORSOrigins(os.Getenv("CORS_ORIGINS")),
		R2AccountID:    getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKey:    getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretKey:    getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2Bucket:       getEnv("R2_BUCKET_NAME", ""),
		R2PublicDomain: getEnv("R2_PUBLIC_DOMAIN", ""),

		AuthRateLimit:          parsePositiveInt(getEnv("AUTH_RATE_LIMIT", "30"), 30),

		GeminiAPIKey:           getEnv("GEMINI_API_KEY", ""),
		GeminiModel:            getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
		GeminiBaseURL:          getEnv("GEMINI_BASE_URL", "https://generativelanguage.googleapis.com"),
		ReceiptExtractEnabled:  getEnv("RECEIPT_EXTRACT_ENABLED", "false") == "true",
		ReceiptRateLimitPerDay: parsePositiveInt(getEnv("RECEIPT_EXTRACT_RATE_LIMIT_PER_DAY", "20"), 20),
	}

	if isRelease {
		if cfg.JWTSecret == "dev-secret-key" {
			slog.Error("JWT_SECRET must be set in production (GIN_MODE=release)")
			os.Exit(1)
		}
		if cfg.DatabaseURL == "" {
			slog.Error("DATABASE_URL must be set in production (GIN_MODE=release)")
			os.Exit(1)
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

func parsePositiveInt(s string, defaultValue int) int {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return defaultValue
	}
	return n
}

func parseDurationDays(s string) time.Duration {
	days, err := strconv.Atoi(s)
	if err != nil || days <= 0 {
		return 7 * 24 * time.Hour
	}
	return time.Duration(days) * 24 * time.Hour
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
