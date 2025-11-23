#!/usr/bin/env bash
# scripts/clean.sh
# Cleans containers, volumes, and local build artifacts.

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

if [ -z "$COMPOSE_CMD" ]; then
  echo "❌ docker-compose or docker compose is required. Install Docker first." >&2
  exit 1
fi

compose_binary=${COMPOSE_CMD%% *}
if ! command -v $compose_binary >/dev/null 2>&1; then
  echo "❌ $COMPOSE_CMD is required. Install Docker and docker-compose first." >&2
  exit 1
fi

echo "==> Stopping containers and removing volumes"
$COMPOSE_CMD down -v

echo "==> Removing build caches and coverage reports"
rm -rf backend/coverage frontend/coverage coverage

if command -v docker >/dev/null 2>&1; then
  docker system prune -f >/dev/null 2>&1 || true
fi

echo "✅ Cleanup completed"
