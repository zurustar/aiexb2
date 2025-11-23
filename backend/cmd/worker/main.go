// backend/cmd/worker/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/esms/internal/queue"
	"github.com/your-org/esms/internal/repository"
	"github.com/your-org/esms/internal/service"
)

// Config はワーカー設定
type Config struct {
	DatabaseURL string
	RedisURL    string
	WorkerCount int
}

func main() {
	// 設定読み込み
	config := loadConfig()

	// ロガー初期化
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting ESMS Background Worker...")

	// データベース接続
	dbPool, err := initDatabase(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbPool.Close()
	log.Println("Database connection established")

	// ジョブキュー初期化
	jobQueue := queue.NewJobQueue(config.RedisURL)
	log.Println("Job queue initialized")

	// リポジトリ初期化
	userRepo := repository.NewUserRepository(dbPool)
	auditLogRepo := repository.NewAuditLogRepository(dbPool)

	// サービス初期化
	notificationService := service.NewNotificationService(
		userRepo,
		jobQueue,
		nil, // EmailSender は実装に応じて初期化
	)

	// ワーカー起動
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 複数のワーカーゴルーチンを起動
	for i := 0; i < config.WorkerCount; i++ {
		go worker(ctx, i, jobQueue, notificationService)
	}

	log.Printf("Started %d worker(s)", config.WorkerCount)

	// グレースフルシャットダウン
	gracefulShutdown(cancel, dbPool)
}

// loadConfig は環境変数から設定を読み込みます
func loadConfig() *Config {
	workerCount := 5
	if count := os.Getenv("WORKER_COUNT"); count != "" {
		// Parse worker count
		workerCount = 5 // デフォルト
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/esms?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		WorkerCount: workerCount,
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
		return nil, err
	}

	// 接続プール設定（ワーカー用に少なめ）
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

// worker はジョブを処理するワーカー
func worker(ctx context.Context, id int, jobQueue *queue.JobQueue, notificationService *service.NotificationService) {
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping", id)
			return
		default:
			// ジョブをデキュー
			job, err := jobQueue.Dequeue(ctx)
			if err != nil {
				// キューが空の場合は少し待機
				time.Sleep(1 * time.Second)
				continue
			}

			// ジョブ処理
			if err := processJob(ctx, job, notificationService); err != nil {
				log.Printf("Worker %d: Failed to process job %s: %v", id, job.ID, err)
			} else {
				log.Printf("Worker %d: Successfully processed job %s", id, job.ID)
			}
		}
	}
}

// processJob はジョブを処理します
func processJob(ctx context.Context, job *queue.Job, notificationService *service.NotificationService) error {
	switch job.Type {
	case "send_email":
		// メール送信ジョブ
		log.Printf("Processing email job: %s", job.ID)
		// notificationService を使用してメール送信
		return nil

	case "cleanup":
		// クリーンアップジョブ
		log.Printf("Processing cleanup job: %s", job.ID)
		return nil

	default:
		log.Printf("Unknown job type: %s", job.Type)
		return nil
	}
}

// gracefulShutdown はグレースフルシャットダウンを処理します
func gracefulShutdown(cancel context.CancelFunc, dbPool *pgxpool.Pool) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")

	// ワーカーの停止
	cancel()

	// 処理中のジョブが完了するまで待機
	time.Sleep(5 * time.Second)

	// データベース接続のクローズ
	dbPool.Close()

	log.Println("Worker exited")
}
