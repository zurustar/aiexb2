// backend/internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config はアプリケーション設定を保持する構造体
type Config struct {
	AppEnv        string
	ServerPort    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	OIDCProvider  string
	OIDCClientID  string
	OIDCSecret    string
	OIDCRedirect  string
	AuditSecret   string // 監査ログ署名用シークレット
}

// Load は環境変数から設定を読み込みます
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "esms"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		OIDCProvider:  getEnv("OIDC_PROVIDER", ""),
		OIDCClientID:  getEnv("OIDC_CLIENT_ID", ""),
		OIDCSecret:    getEnv("OIDC_CLIENT_SECRET", ""),
		OIDCRedirect:  getEnv("OIDC_REDIRECT_URL", "http://localhost:8080/auth/callback"),
		AuditSecret:   getEnv("AUDIT_SECRET", "default-audit-secret-key"),
	}

	redisDBStr := getEnv("REDIS_DB", "0")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}
	cfg.RedisDB = redisDB

	return cfg, nil
}

// DSN はデータベース接続文字列を返します
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// RedisAddr はRedis接続アドレスを返します
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// IsDev は開発環境かどうかを返します
func (c *Config) IsDev() bool {
	return c.AppEnv == "development"
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetDurationEnv は環境変数をtime.Durationとして取得します
func GetDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}
