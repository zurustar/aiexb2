// backend/cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/your-org/esms/internal/handler"
	"github.com/your-org/esms/internal/repository"
	"github.com/your-org/esms/internal/service"
	"github.com/your-org/esms/pkg/oidc"
)

// Config はアプリケーション設定
type Config struct {
	Port         string
	DatabaseURL  string
	RedisURL     string
	OIDCIssuer   string
	OIDCClientID string
	OIDCSecret   string
}

func main() {
	// 設定読み込み
	config := loadConfig()

	// ロガー初期化
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting ESMS API Server...")

	// データベース接続
	dbPool, err := initDatabase(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbPool.Close()
	log.Println("Database connection established")

	// 依存サービスのヘルスチェック
	if err := healthCheck(dbPool); err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	log.Println("All dependency health checks passed")

	// Redis接続（将来の実装用）
	// redisClient := initRedis(config.RedisURL)
	// defer redisClient.Close()

	// OIDC クライアント初期化
	oidcClient, err := initOIDCClient(config)
	if err != nil {
		log.Fatalf("Failed to initialize OIDC client: %v", err)
	}
	log.Println("OIDC client initialized")

	// pgxpool.Pool を *sql.DB に変換
	db := stdlib.OpenDBFromPool(dbPool)
	defer db.Close()

	// リポジトリ初期化
	userRepo := repository.NewUserRepository(db)
	resourceRepo := repository.NewResourceRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	// サービス初期化
	authService := service.NewAuthService(oidcClient, userRepo, auditLogRepo)
	reservationService := service.NewReservationService(
		reservationRepo,
		resourceRepo,
		userRepo,
		auditLogRepo,
	)
	approvalService := service.NewApprovalService(
		reservationRepo,
		userRepo,
		auditLogRepo,
	)

	// ルーター初期化
	router := handler.NewRouter(
		authService,
		reservationService,
		approvalService,
		userRepo,
		resourceRepo,
	)

	// HTTPサーバー設定
	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router.GetRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// サーバー起動（ゴルーチン）
	go func() {
		log.Printf("Server listening on port %s", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// グレースフルシャットダウン
	gracefulShutdown(server, dbPool)
}

// loadConfig は環境変数から設定を読み込みます
func loadConfig() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/esms?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		OIDCIssuer:   getEnv("OIDC_ISSUER", ""),
		OIDCClientID: getEnv("OIDC_CLIENT_ID", ""),
		OIDCSecret:   getEnv("OIDC_CLIENT_SECRET", ""),
	}
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// initDatabase はデータベース接続プールを初期化します
func initDatabase(databaseURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// 接続プール設定
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// 接続テスト
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// initOIDCClient はOIDCクライアントを初期化します
func initOIDCClient(config *Config) (*oidc.Client, error) {
	if config.OIDCIssuer == "" {
		log.Println("Warning: OIDC not configured, using mock client")
		// 開発環境用のモッククライアント
		return nil, nil
	}

	oidcConfig := &oidc.Config{
		IssuerURL:    config.OIDCIssuer,
		ClientID:     config.OIDCClientID,
		ClientSecret: config.OIDCSecret,
		RedirectURL:  fmt.Sprintf("http://localhost:%s/api/v1/auth/callback", config.Port),
	}

	return oidc.NewClient(context.Background(), oidcConfig)
}

// healthCheck は依存サービスのヘルスチェックを実行します
func healthCheck(dbPool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// データベース接続確認
	if err := dbPool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// 将来的にRedisやその他のサービスのヘルスチェックを追加
	// if err := redisClient.Ping(ctx).Err(); err != nil {
	//     return fmt.Errorf("redis health check failed: %w", err)
	// }

	return nil
}

// gracefulShutdown はグレースフルシャットダウンを処理します
func gracefulShutdown(server *http.Server, dbPool *pgxpool.Pool) {
	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("Received signal: %v. Shutting down server...", sig)

	// シャットダウンのタイムアウト設定
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// HTTPサーバーのシャットダウン
	log.Println("Stopping HTTP server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("HTTP server stopped gracefully")
	}

	// データベース接続のクローズ
	log.Println("Closing database connections...")
	dbPool.Close()
	log.Println("Database connections closed")

	log.Println("Server exited")
}
