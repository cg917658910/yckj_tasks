package config

import (
	"fmt"
	"os"
)

type Config struct {
	Env         string
	Addr        string
	DBDSN       string
	JWTSecret   string
	JWTIssuer   string
	JWTTTL      int64
	UploadDir   string
	BaseURL     string
	AllowOrigin string
}

func Load() Config {
	return Config{
		Env:         getEnv("APP_ENV", "dev"),
		Addr:        getEnv("APP_ADDR", ":8080"),
		DBDSN:       getEnv("DB_DSN", "root:cg123456@tcp(127.0.0.1:3306)/task_system?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me"),
		JWTIssuer:   getEnv("JWT_ISSUER", "task-system"),
		JWTTTL:      getEnvInt64("JWT_TTL", 7200),
		UploadDir:   getEnv("UPLOAD_DIR", "./public/uploads"),
		BaseURL:     getEnv("BASE_URL", "http://127.0.0.1:8080"),
		AllowOrigin: getEnv("ALLOW_ORIGIN", "*"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		var out int64
		_, err := fmt.Sscanf(v, "%d", &out)
		if err == nil {
			return out
		}
	}
	return def
}
