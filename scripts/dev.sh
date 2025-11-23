#!/usr/bin/env bash
# scripts/dev.sh
# Bootstraps the development environment using docker-compose.

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
REBUILD=${1:-}

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

if [ ! -f .env ]; then
  echo "⚠️  .env not found. Run ./scripts/setup.sh first."
  exit 1
fi

echo "==> Starting development services"
if [ "$REBUILD" == "--rebuild" ]; then
  $COMPOSE_CMD build --no-cache
fi

$COMPOSE_CMD up -d backend worker frontend redis postgres mailhog

cat <<INFO

✅ Development stack is running.
- Backend API:   http://localhost:8080
- Frontend:      http://localhost:3000
- MailHog UI:    http://localhost:8025

Use '$COMPOSE_CMD logs -f backend' to tail backend logs.
INFO
