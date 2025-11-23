#!/bin/bash
# scripts/setup.sh
# 初回セットアップスクリプト
#
# 実行内容:
# - .envファイルの作成
# - Dockerイメージのビルド
# - データベースマイグレーション
# - シードデータ投入

set -e

echo "==> ESMS Setup Script"
echo ""

# .env作成
if [ ! -f .env ]; then
    echo "Creating .env from .env.example..."
    cp .env.example .env
    echo "✅ .env created. Please edit it with your configuration."
else
    echo "ℹ️  .env already exists, skipping..."
fi

# Dockerビルド
echo ""
echo "Building Docker images..."
docker-compose build

# マイグレーション
echo ""
echo "Running database migrations..."
make migrate

# シードデータ
echo ""
echo "Seeding database..."
make seed

echo ""
echo "✅ Setup complete!"
echo "Run 'make up' to start the development environment."
