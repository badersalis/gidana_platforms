package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv            string
	Port              string
	DatabaseURL       string
	DBPath            string
	JWTSecret         string
	JWTExpiryHours    int
	UploadDir         string
	MaxUploadSizeMB   int64
	SupabaseURL       string
	SupabaseKey       string
	SupabaseBucket    string
	AllowedOrigins    string
	CinetPayAPISecret string
}

var App *Config

func Load() {
	_ = godotenv.Load()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "72"))
	maxUpload, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "5"), 10, 64)
	App = &Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		DBPath:            getEnv("DB_PATH", "gidana_dev.db"),
		JWTSecret:         getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiryHours:    jwtExpiry,
		UploadDir:         getEnv("UPLOAD_DIR", "./uploads/properties"),
		MaxUploadSizeMB:   maxUpload,
		SupabaseURL:       getEnv("SUPABASE_URL", ""),
		SupabaseKey:       getEnv("SUPABASE_KEY", ""),
		SupabaseBucket:    getEnv("SUPABASE_BUCKET", ""),
		AllowedOrigins:    getEnv("ALLOWED_ORIGINS", "*"),
		CinetPayAPISecret: getEnv("CINETPAY_API_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
