#!/usr/bin/env bash
# scripts/setup.sh
# Initial project bootstrap: env creation, dependency containers, migrations, and seed data.

set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)
resolve_compose() {
  if [ -n "${COMPOSE_CMD:-}" ]; then
    echo "$COMPOSE_CMD"
    return
  fi

  if command -v docker-compose >/dev/null 2>&1; then
    echo "docker-compose"
    return
  fi

  if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    echo "docker compose"
    return
  fi

  echo ""
}

COMPOSE_CMD=$(resolve_compose)

cd "$PROJECT_ROOT"

echo "==> ESMS Setup Script"

if [ -z "$COMPOSE_CMD" ]; then
  echo "❌ docker-compose or docker compose is required. Install Docker first." >&2
  exit 1
fi

compose_binary=${COMPOSE_CMD%% *}
if ! command -v $compose_binary >/dev/null 2>&1; then
  echo "❌ $COMPOSE_CMD is required. Install Docker and docker-compose first." >&2
  exit 1
fi

if [ ! -f .env ]; then
  echo "Creating .env from .env.example..."
  cp .env.example .env
  echo "✅ .env created. Please update values as needed."
else
  echo "ℹ️  .env already exists, skipping creation."
fi

echo "==> Building Docker images"
$COMPOSE_CMD build

echo "==> Starting core services (postgres, redis, mailhog)"
$COMPOSE_CMD up -d postgres redis mailhog

wait_for_postgres() {
  for _ in $(seq 1 20); do
    if $COMPOSE_CMD exec -T postgres pg_isready -U esms_user -d esms >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done
  return 1
}

if wait_for_postgres; then
  echo "✅ postgres is ready"
else
  echo "❌ postgres did not become ready in time" >&2
  exit 1
fi

echo "==> Running database migrations"
"$PROJECT_ROOT/backend/scripts/migrate.sh" up

echo "==> Seeding database"
$COMPOSE_CMD run --rm backend sh -c "cd /app && go run scripts/seed.go"

echo "\n✅ Setup complete!"
echo "Use ./scripts/dev.sh to start the full stack."
