# Production Deployment & Final Polish

Feature 12 completes PeakCloud's production deployment foundation and prepares the application for repeatable production-oriented deployment.

## Production Architecture

PeakCloud runs as five containerized services:

- PostgreSQL for relational application data
- Redis for caching and session infrastructure
- MinIO for S3-compatible object storage
- Go API for backend services
- Next.js for the web application

Production orchestration is defined in `docker-compose.production.yml`.

## Production Containers

The root `Dockerfile` builds the Go API using a multi-stage build. The builder uses the Go version required by `go.mod`, while the final runtime image contains only the compiled application and required runtime dependencies.

`apps/web/Dockerfile` builds the Next.js application using multiple stages for dependency installation, compilation, and runtime.

Next.js uses `output: "standalone"` so the final image contains only the files required to serve the application.

Both application containers use non-root runtime users.

## Production Configuration

Production configuration is provided through environment variables.

`.env.production.example` documents the required configuration while `.env.production` contains deployment-specific values and is excluded from Git.

Real secrets must never be committed to the repository.

## Security

PeakCloud preserves the production security validation introduced in Feature 10.

When `APP_ENV=production`, insecure configuration causes application startup to fail.

Real production deployments require secure session cookies, HTTPS application URLs, and TLS-enabled object storage.

Local verification can use development-mode configuration with localhost HTTP endpoints without weakening the requirements for a real production deployment.

## Health and Observability

The API exposes:

- `GET /live` for process liveness
- `GET /health` for dependency-aware health
- `GET /metrics` for HTTP metrics

These endpoints support container health checks, operational monitoring, and deployment verification.

## Deployment Verification

PeakCloud includes:

- `scripts/production-smoke-test.sh`
- `scripts/verify-production.sh`

The smoke test verifies API liveness, dependency health, metrics availability, and web application availability.

The complete production-style stack was successfully validated locally with PostgreSQL, Redis, MinIO, the PeakCloud API, and the PeakCloud web application.

## Result

Feature 12 gives PeakCloud a repeatable production-oriented deployment foundation with containerized builds, service orchestration, persistent infrastructure, environment isolation, health checks, smoke testing, and documented operational procedures.
