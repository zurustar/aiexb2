// backend/cmd/worker/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/your-org/esms/internal/cache"
	"github.com/your-org/esms/internal/config"
	"github.com/your-org/esms/internal/queue"
	"github.com/your-org/esms/internal/repository"
	"github.com/your-org/esms/internal/service"
)

func main() {
	// ロガー初期化
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting ESMS Background Worker...")

	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// データベース接続
	dbPool, err := initDatabase(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbPool.Close()
	log.Println("Database connection established")

	// pgxpool.Pool を *sql.DB に変換
	db := stdlib.OpenDBFromPool(dbPool)
	defer db.Close()

	// Redis接続
	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	log.Println("Redis connection established")

	// ジョブキュー初期化
	jobQueue := queue.NewRedisJobQueue(redisClient, "default")
	log.Println("Job queue initialized")

	// リポジトリ初期化
	userRepo := repository.NewUserRepository(db)

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
	workerCount := 5 // デフォルト
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i, jobQueue, notificationService)
	}

	log.Printf("Started %d worker(s)", workerCount)

	// グレースフルシャットダウン
	gracefulShutdown(cancel, &wg, dbPool)
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
func worker(ctx context.Context, wg *sync.WaitGroup, id int, jobQueue queue.JobQueue, notificationService *service.NotificationService) {
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
				log.Printf("Worker %d: Failed to dequeue job: %v", id, err)
				time.Sleep(1 * time.Second)
				continue
			}

			// キューが空の場合
			if job == nil {
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
