// backend/internal/config/config_test.go
package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/config"
)

func TestLoad(t *testing.T) {
	// 環境変数を一時的に設定
	os.Setenv("APP_ENV", "test")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("REDIS_DB", "1")
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("REDIS_DB")
	}()

	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.Equal(t, "test", cfg.AppEnv)
	assert.Equal(t, "9090", cfg.ServerPort)
	assert.Equal(t, 1, cfg.RedisDB)
	assert.Equal(t, "localhost", cfg.DBHost) // デフォルト値
	assert.False(t, cfg.UseSecretsManager)
}

func TestLoad_WithSecretsManager(t *testing.T) {
	os.Setenv("USE_SECRETS_MANAGER", "true")
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("AWS_SECRET_ID", "my-secret")
	defer func() {
		os.Unsetenv("USE_SECRETS_MANAGER")
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_SECRET_ID")
	}()

	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.True(t, cfg.UseSecretsManager)
	assert.Equal(t, "us-west-2", cfg.AWSRegion)
	assert.Equal(t, "my-secret", cfg.AWSSecretID)
}

func TestConfig_DSN(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "db",
		DBPort:     "5432",
		DBUser:     "user",
		DBPassword: "pass",
		DBName:     "mydb",
		DBSSLMode:  "disable",
	}
	expected := "host=db port=5432 user=user password=pass dbname=mydb sslmode=disable"
	assert.Equal(t, expected, cfg.DSN())
}

func TestConfig_RedisAddr(t *testing.T) {
	cfg := &config.Config{
		RedisHost: "redis",
		RedisPort: "6379",
	}
	expected := "redis:6379"
	assert.Equal(t, expected, cfg.RedisAddr())
}

func TestGetDurationEnv(t *testing.T) {
	os.Setenv("TEST_DURATION", "10s")
	defer os.Unsetenv("TEST_DURATION")

	val := config.GetDurationEnv("TEST_DURATION", 5*time.Second)
	assert.Equal(t, 10*time.Second, val)

	valDefault := config.GetDurationEnv("NON_EXISTENT", 5*time.Second)
	assert.Equal(t, 5*time.Second, valDefault)
}
