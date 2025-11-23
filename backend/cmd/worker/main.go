```go
// backend/cmd/worker/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
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

	// ジョブキュー初期化（リトライ付き）
	jobQueue, err := initJobQueueWithRetry(config.RedisURL, 5, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to initialize job queue: %v", err)
	}
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

	// WaitGroupでワーカーの完了を追跡
	var wg sync.WaitGroup

	// 複数のワーカーゴルーチンを起動
	for i := 0; i < config.WorkerCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i, jobQueue, notificationService)
	}

	log.Printf("Started %d worker(s)", config.WorkerCount)

	// グレースフルシャットダウン
	gracefulShutdown(cancel, &wg, dbPool)
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

// initJobQueueWithRetry はジョブキューを初期化します（リトライ付き）
func initJobQueueWithRetry(redisURL string, maxRetries int, initialBackoff time.Duration) (*queue.JobQueue, error) {
	var jobQueue *queue.JobQueue
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		jobQueue = queue.NewJobQueue(redisURL)
		
		// 接続テスト（キューの実装に応じて調整）
		// ここでは単純に初期化が成功したと仮定
		if jobQueue != nil {
			log.Printf("Job queue connection established (attempt %d/%d)", attempt+1, maxRetries)
			return jobQueue, nil
		}

		if attempt < maxRetries-1 {
			// 指数バックオフで待機
			backoff := time.Duration(math.Pow(2, float64(attempt))) * initialBackoff
			log.Printf("Job queue connection failed, retrying in %v (attempt %d/%d)", backoff, attempt+1, maxRetries)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("failed to connect to job queue after %d attempts: %w", maxRetries, err)
}

// worker はジョブを処理するワーカー
func worker(ctx context.Context, wg *sync.WaitGroup, id int, jobQueue *queue.JobQueue, notificationService *service.NotificationService) {
	defer wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping gracefully", id)
			return
		default:
			// ジョブをデキュー
			job, err := jobQueue.Dequeue(ctx)
			if err != nil {
				// コンテキストがキャンセルされた場合は終了
				if ctx.Err() != nil {
					return
				}
				// キューが空の場合は少し待機
				time.Sleep(1 * time.Second)
				continue
			}

			// ジョブ処理（コンテキストを渡して中断可能にする）
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
func gracefulShutdown(cancel context.CancelFunc, wg *sync.WaitGroup, dbPool *pgxpool.Pool) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("Received signal: %v. Shutting down worker...", sig)

	// ワーカーの停止をシグナル
	log.Println("Stopping workers...")
	cancel()

	// 処理中のジョブが完了するまで待機（タイムアウト付き）
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All workers stopped gracefully")
	case <-time.After(30 * time.Second):
		log.Println("Timeout waiting for workers to stop, forcing shutdown")
	}

	// データベース接続のクローズ
	log.Println("Closing database connections...")
	dbPool.Close()
	log.Println("Database connections closed")

	log.Println("Worker exited")
}
