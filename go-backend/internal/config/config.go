package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env         string `yaml:"env"`
	Addr        string `yaml:"addr"`
	DBDSN       string `yaml:"db_dsn"`
	JWTSecret   string `yaml:"jwt_secret"`
	JWTIssuer   string `yaml:"jwt_issuer"`
	JWTTTL      int64  `yaml:"jwt_ttl"`
	UploadDir   string `yaml:"upload_dir"`
	BaseURL     string `yaml:"base_url"`
	AllowOrigin string `yaml:"allow_origin"`
}

// Load 会优先尝试从 config.yaml 读取配置,找不到或解析失败时回退到环境变量/默认值
func Load() Config {
	// 默认值
	cfg := Config{
		Env:         "dev",
		Addr:        ":8080",
		DBDSN:       "root:cg123456@tcp(127.0.0.1:3306)/task_system?charset=utf8mb4&parseTime=True&loc=Local",
		JWTSecret:   "change-me",
		JWTIssuer:   "task-system",
		JWTTTL:      7200,
		UploadDir:   "./public/uploads",
		BaseURL:     "http://127.0.0.1:8080",
		AllowOrigin: "*",
	}

	// 1. 先尝试读取 YAML 文件,路径可通过 CONFIG_FILE 环境变量覆盖
	configFile := getEnv("CONFIG_FILE", "config.yaml")
	if b, err := os.ReadFile(configFile); err == nil {
		var fileCfg Config
		if err := yaml.Unmarshal(b, &fileCfg); err == nil {
			mergeConfig(&cfg, &fileCfg)
		}
	}

	// 2. 再用环境变量覆盖(如果存在)
	cfg.Env = getEnv("APP_ENV", cfg.Env)
	cfg.Addr = getEnv("APP_ADDR", cfg.Addr)
	cfg.DBDSN = getEnv("DB_DSN", cfg.DBDSN)
	cfg.JWTSecret = getEnv("JWT_SECRET", cfg.JWTSecret)
	cfg.JWTIssuer = getEnv("JWT_ISSUER", cfg.JWTIssuer)
	cfg.JWTTTL = getEnvInt64("JWT_TTL", cfg.JWTTTL)
	cfg.UploadDir = getEnv("UPLOAD_DIR", cfg.UploadDir)
	cfg.BaseURL = getEnv("BASE_URL", cfg.BaseURL)
	cfg.AllowOrigin = getEnv("ALLOW_ORIGIN", cfg.AllowOrigin)

	return cfg
}

// mergeConfig 用非零值覆盖默认配置
func mergeConfig(base, override *Config) {
	// string: 非空则覆盖
	if override.Env != "" {
		base.Env = override.Env
	}
	if override.Addr != "" {
		base.Addr = override.Addr
	}
	if override.DBDSN != "" {
		base.DBDSN = override.DBDSN
	}
	if override.JWTSecret != "" {
		base.JWTSecret = override.JWTSecret
	}
	if override.JWTIssuer != "" {
		base.JWTIssuer = override.JWTIssuer
	}
	if override.JWTTTL != 0 {
		base.JWTTTL = override.JWTTTL
	}
	if override.UploadDir != "" {
		base.UploadDir = override.UploadDir
	}
	if override.BaseURL != "" {
		base.BaseURL = override.BaseURL
	}
	if override.AllowOrigin != "" {
		base.AllowOrigin = override.AllowOrigin
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
