#!/usr/bin/env sh

set -eu

COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.production.yml}"
ENV_FILE="${ENV_FILE:-.env.production}"

if [ ! -f "$ENV_FILE" ]; then
  echo "FAIL: $ENV_FILE does not exist"
  echo "Create it from .env.production.example and provide real secrets."
  exit 1
fi

echo "Validating production Compose configuration..."

docker compose \
  --env-file "$ENV_FILE" \
  -f "$COMPOSE_FILE" \
  config \
  --quiet

echo "PASS: production Compose configuration"

echo
echo "Container status:"

docker compose \
  --env-file "$ENV_FILE" \
  -f "$COMPOSE_FILE" \
  ps

echo
echo "Running production smoke tests..."

API_URL="${API_URL:-http://localhost:${PEAKCLOUD_API_PORT:-8080}}" \
WEB_URL="${WEB_URL:-http://localhost:${PEAKCLOUD_WEB_PORT:-3000}}" \
./scripts/production-smoke-test.sh

echo
echo "PASS: production verification"
