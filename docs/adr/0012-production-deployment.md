# ADR 0012: Containerized Production Deployment

## Status

Accepted

## Context

PeakCloud contains a Go API, Next.js web application, PostgreSQL, Redis, and S3-compatible object storage.

The project requires a repeatable deployment model rather than relying on manually configured development processes.

The deployment architecture must provide reproducible builds, explicit dependencies, environment isolation, persistent infrastructure, health checks, and repeatable verification.

## Decision

PeakCloud will use Docker for application packaging and Docker Compose for production-oriented service orchestration.

The deployment consists of:

- PostgreSQL
- Redis
- MinIO
- PeakCloud API
- PeakCloud web application

## Application Images

The API uses a multi-stage Go Docker build.

The web application uses a multi-stage Node.js build with Next.js standalone output.

Application runtime containers use non-root users.

Docker ignore files reduce build context size and prevent unnecessary development artifacts from entering images.

## Configuration

Deployment configuration is supplied through environment variables.

`.env.production.example` documents the required variables.

The real `.env.production` file is excluded from version control.

Deployment secrets must remain outside the repository.

## Orchestration

`docker-compose.production.yml` defines service networking, dependencies, health checks, restart behavior, persistent volumes, and application ports.

The API depends on healthy infrastructure services.

The web application depends on API health.

## Security

Production configuration validation remains enforced.

When the API runs with `APP_ENV=production`, insecure production configuration causes startup to fail.

Local production-image testing may use development mode to support localhost HTTP services, but this does not weaken real production requirements.

## Observability

The deployment uses the existing operational endpoints:

- `/live`
- `/health`
- `/metrics`

These endpoints provide process liveness, dependency health, and HTTP metrics.

## Verification

PeakCloud includes automated smoke testing and deployment verification scripts.

The deployment is considered ready only when configuration validates, required services are running, health checks pass, the web application is reachable, and smoke tests succeed.

## Consequences

Benefits include reproducible builds, portable deployment, explicit dependencies, smaller runtime images, non-root execution, persistent infrastructure, and repeatable operational verification.

Trade-offs include additional Docker configuration, persistent-volume management, secret management, and responsibility for production TLS and external networking.

## Result

Containerized deployment provides PeakCloud with a portable production foundation while preserving the security, reliability, testing, and observability controls introduced by earlier features.
