#!/usr/bin/env bash
# backend/scripts/migrate.sh
# Database migration helper using golang-migrate inside the backend container.

set -euo pipefail

ACTION=${1:-up}
STEPS=${2:-1}

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
DATABASE_URL=${DATABASE_URL:-"postgres://esms_user:esms_pass@postgres:5432/esms?sslmode=disable"}
MIGRATIONS_PATH=${MIGRATIONS_PATH:-/app/migrations}

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

usage() {
  cat <<USAGE
Usage: $(basename "$0") [up|down|force <version>|version|redo|steps <n>]

Examples:
  $(basename "$0") up           # apply all pending migrations
  $(basename "$0") down         # rollback the last batch
  $(basename "$0") redo         # rollback one step then re-apply
  $(basename "$0") steps 3      # apply 3 migrations
  $(basename "$0") force 1      # force schema version
USAGE
}

if [[ "$ACTION" == "-h" || "$ACTION" == "--help" ]]; then
  usage
  exit 0
fi

if [ -z "$COMPOSE_CMD" ]; then
  echo "❌ docker-compose or docker compose is required to run migrations." >&2
  exit 1
fi

compose_binary=${COMPOSE_CMD%% *}
if ! command -v $compose_binary >/dev/null 2>&1; then
  echo "❌ $COMPOSE_CMD is required to run migrations." >&2
  exit 1
fi

run_migrate() {
  local migrate_args=$1
  $COMPOSE_CMD run --rm backend migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" $migrate_args
}

case "$ACTION" in
  up)
    run_migrate "up"
    ;;
  down)
    run_migrate "down"
    ;;
  redo)
    run_migrate "down 1"
    run_migrate "up"
    ;;
  steps)
    run_migrate "up $STEPS"
    ;;
  force)
    run_migrate "force $STEPS"
    ;;
  version)
    run_migrate "version"
    ;;
  *)
    usage
    exit 1
    ;;
fi
