# PeakCloud Deployment Runbook

This runbook describes how to configure, build, start, verify, inspect, and stop the PeakCloud deployment.

## Prerequisites

The deployment host requires Docker and Docker Compose.

## Configure the Environment

Create the deployment environment file:

    cp .env.production.example .env.production

Replace all `CHANGE_ME` values with deployment-specific credentials.

Important variables include:

- `DATABASE_URL`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`
- `MINIO_ROOT_USER`
- `MINIO_ROOT_PASSWORD`
- `S3_ACCESS_KEY`
- `S3_SECRET_KEY`
- `WEB_URL`
- `NEXT_PUBLIC_API_URL`

Confirm that the real environment file is ignored:

    git check-ignore -v .env.production

Never commit `.env.production`.

## Validate Configuration

Run:

    docker compose --env-file .env.production -f docker-compose.production.yml config

The command must complete successfully before deployment.

## Build

Run:

    docker compose --env-file .env.production -f docker-compose.production.yml build

This builds the production API and web images.

## Start

Run:

    docker compose --env-file .env.production -f docker-compose.production.yml up -d

## Check Services

Run:

    docker compose --env-file .env.production -f docker-compose.production.yml ps

The expected services are:

- postgres
- redis
- minio
- api
- web

PostgreSQL, Redis, MinIO, and the API should become healthy.

## Verify Liveness

Run:

    curl -fsS http://localhost:8080/live

Expected local response:

    {"status":"ok"}

## Verify Dependency Health

Run:

    curl -fsS http://localhost:8080/health

The API should report healthy PostgreSQL and Redis dependencies.

## Verify Metrics

Run:

    curl -fsS http://localhost:8080/metrics

The endpoint should return PeakCloud HTTP metrics.

## Verify Web Application

Run:

    curl -fsS http://localhost:3000/

The web application should return successfully.

## Run Smoke Test

Run:

    bash scripts/production-smoke-test.sh

A successful run ends with:

    Production smoke test passed

## Run Production Verification

Run:

    bash scripts/verify-production.sh

## Inspect Logs

All services:

    docker compose --env-file .env.production -f docker-compose.production.yml logs

API:

    docker compose --env-file .env.production -f docker-compose.production.yml logs --tail=100 api

Web:

    docker compose --env-file .env.production -f docker-compose.production.yml logs --tail=100 web

## Stop

Run:

    docker compose --env-file .env.production -f docker-compose.production.yml down

Normal shutdown preserves persistent volumes.

Do not remove production volumes unless permanent deletion of stored data is intended.

## Production Security

A real production deployment should use `APP_ENV=production`, HTTPS URLs, secure cookies, TLS-enabled object storage, and deployment-specific credentials.

Local container verification may use `APP_ENV=development`, localhost HTTP URLs, `SESSION_SECURE=false`, and `S3_USE_SSL=false`.

That configuration is for local testing only.
