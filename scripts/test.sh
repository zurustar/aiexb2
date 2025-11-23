#!/usr/bin/env bash
# scripts/test.sh
# Runs backend and frontend test suites with coverage reports.

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
TEST_DB_URL=${TEST_DATABASE_URL:-"postgres://esms_user:esms_pass@postgres:5432/esms?sslmode=disable"}
BACKEND_COVERAGE_DIR="$PROJECT_ROOT/backend/coverage"
FRONTEND_COVERAGE_JSON="$PROJECT_ROOT/frontend/coverage/coverage-summary.json"
SUMMARY_REPORT="$PROJECT_ROOT/coverage/coverage-summary.txt"

cd "$PROJECT_ROOT"
mkdir -p "$PROJECT_ROOT/coverage" "$BACKEND_COVERAGE_DIR"

if [ -z "$COMPOSE_CMD" ]; then
  echo "❌ docker-compose or docker compose is required. Install Docker first." >&2
  exit 1
fi

compose_binary=${COMPOSE_CMD%% *}
if ! command -v $compose_binary >/dev/null 2>&1; then
  echo "❌ $COMPOSE_CMD is required. Install Docker and docker-compose first." >&2
  exit 1
fi

echo "==> Starting database dependencies for tests"
$COMPOSE_CMD up -d postgres redis

# Ensure migrations are applied before running integration tests
"$PROJECT_ROOT/backend/scripts/migrate.sh" up

echo "==> Running backend tests with coverage"
$COMPOSE_CMD run --rm -e TEST_DATABASE_URL="$TEST_DB_URL" backend sh -c "cd /app && mkdir -p coverage && go test ./... -covermode=atomic -coverprofile=coverage/backend.out -count=1"
$COMPOSE_CMD run --rm backend sh -c "cd /app && go tool cover -func=coverage/backend.out > coverage/backend-coverage.txt"

BACKEND_RATE=$(awk '/^total:/ {gsub("%", "", $3); print $3}' "$BACKEND_COVERAGE_DIR/backend-coverage.txt")

if [ -z "$BACKEND_RATE" ]; then
  echo "⚠️  Could not determine backend coverage percentage"
fi

echo "==> Running frontend tests with coverage"
$COMPOSE_CMD run --rm frontend sh -c "cd /app && npm test -- --coverage --runInBand"

FRONTEND_RATE=$(python3 - <<'PY'
import json
import pathlib
summary_path = pathlib.Path("frontend/coverage/coverage-summary.json")
if not summary_path.exists():
    print("")
else:
    data = json.loads(summary_path.read_text())
    total = data.get("total", {})
    statements = total.get("statements", {})
    print(statements.get("pct", ""))
PY
)

cat > "$SUMMARY_REPORT" <<REPORT
# Test Coverage Summary
Backend statements coverage: ${BACKEND_RATE:-"N/A"}%
Frontend statements coverage: ${FRONTEND_RATE:-"N/A"}%
Target: Backend 80% / Frontend 70%
REPORT

echo "==> Coverage summary"
cat "$SUMMARY_REPORT"

check_threshold() {
  local rate=$1
  local target=$2
  python3 - <<'PY'
import os, sys
rate = os.environ.get('RATE')
target = float(os.environ.get('TARGET', '0'))
try:
    if float(rate) + 1e-9 < target:
        sys.exit(1)
except (TypeError, ValueError):
    sys.exit(1)
PY
}

backend_ok=true
frontend_ok=true
if ! RATE="$BACKEND_RATE" TARGET="80" check_threshold; then
  echo "⚠️  Backend coverage is below 80%" >&2
  backend_ok=false
fi
if ! RATE="$FRONTEND_RATE" TARGET="70" check_threshold; then
  echo "⚠️  Frontend coverage is below 70%" >&2
  frontend_ok=false
fi

if [ "$backend_ok" = false ] || [ "$frontend_ok" = false ]; then
  echo "⚠️  Coverage targets not met. See summary above." >&2
fi
