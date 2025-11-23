.PHONY: help setup up down restart logs logs-f ps test clean migrate seed

help: ## このヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## 初回セットアップ（.env作成、ビルド、マイグレーション）
	@echo "==> Setting up project..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env from .env.example"; fi
	@docker-compose build
	@echo "==> Running migrations..."
	@$(MAKE) migrate
	@echo "==> Seeding database..."
	@$(MAKE) seed
	@echo "✅ Setup complete! Run 'make up' to start services."

up: ## 開発環境起動
	@docker-compose up -d
	@echo ""
	@echo "✅ Services started!"
	@echo "   Frontend:  http://localhost:3000"
	@echo "   Backend:   http://localhost:8080"
	@echo "   MailHog:   http://localhost:8025"
	@echo ""
	@echo "Run 'make logs' to view logs"

down: ## 開発環境停止
	@docker-compose down

restart: ## 開発環境再起動
	@docker-compose restart

logs: ## ログ表示（全サービス）
	@docker-compose logs --tail=100

logs-f: ## ログ表示（フォロー）
	@docker-compose logs -f

ps: ## サービス状態確認
	@docker-compose ps

test: ## テスト実行
	@echo "==> Running backend tests..."
	@docker-compose run --rm backend go test ./... -v
	@echo "==> Running frontend tests..."
	@docker-compose run --rm frontend npm test

migrate: ## DBマイグレーション実行
	@echo "==> Running database migrations..."
	@docker-compose run --rm backend sh -c "cd /app && go run cmd/migrate/main.go up"

migrate-down: ## DBマイグレーションロールバック
	@echo "==> Rolling back database migrations..."
	@docker-compose run --rm backend sh -c "cd /app && go run cmd/migrate/main.go down"

seed: ## シードデータ投入
	@echo "==> Seeding database..."
	@docker-compose exec -T postgres psql -U esms_user -d esms < database/seed/users.sql
	@docker-compose exec -T postgres psql -U esms_user -d esms < database/seed/resources.sql

clean: ## クリーンアップ（コンテナ、ボリューム、ビルドキャッシュ削除）
	@echo "==> Cleaning up..."
	@docker-compose down -v
	@docker system prune -f
	@echo "✅ Cleanup complete!"

shell-backend: ## バックエンドコンテナにシェル接続
	@docker-compose exec backend sh

shell-frontend: ## フロントエンドコンテナにシェル接続
	@docker-compose exec frontend sh

shell-db: ## PostgreSQLに接続
	@docker-compose exec postgres psql -U esms_user -d esms

redis-cli: ## Redisに接続
	@docker-compose exec redis redis-cli
