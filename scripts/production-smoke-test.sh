#!/usr/bin/env sh

set -eu

API_URL="${API_URL:-http://localhost:8080}"
WEB_URL="${WEB_URL:-http://localhost:3000}"

echo "PeakCloud production smoke test"
echo "API: $API_URL"
echo "Web: $WEB_URL"

echo
echo "Checking API liveness..."
curl --fail --silent --show-error \
  "$API_URL/live" >/dev/null
echo "PASS: API liveness"

echo
echo "Checking API health..."
curl --fail --silent --show-error \
  "$API_URL/health" >/dev/null
echo "PASS: API dependency health"

echo
echo "Checking metrics..."
curl --fail --silent --show-error \
  "$API_URL/metrics" >/dev/null
echo "PASS: API metrics"

echo
echo "Checking web application..."
curl --fail --silent --show-error \
  "$WEB_URL/" >/dev/null
echo "PASS: web application"

echo
echo "Production smoke test passed"
