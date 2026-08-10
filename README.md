# PeakCloud

PeakCloud is a production-grade cloud storage and file synchronization platform.

## Status

PeakCloud is under active development.

Current milestone:

**Feature 2 — User Authentication & Account System**

## Tech Stack

### Frontend

- Next.js
- React
- TypeScript
- Tailwind CSS

### Backend

- Go
- Gin

### Data & Infrastructure

- PostgreSQL
- Redis
- S3-compatible object storage with MinIO
- Docker Compose

### Engineering

- GitHub Actions
- Automated testing
- Structured logging
- Health checks
- Graceful shutdown

## Local Development

Start infrastructure:

    make infra-up

Run the API:

    make api

Run the web application:

    make web

Run all checks:

    make check

## API

Health endpoint:

    GET /health

## Development Workflow

PeakCloud is developed incrementally using GitHub Issues, feature branches, pull requests, automated CI, tests, Architecture Decision Records, semantic versioning, and production-focused documentation.

## Authentication

PeakCloud currently supports:

- Account registration
- Argon2id password hashing
- Login
- Server-side sessions
- HttpOnly session cookies
- Protected API routes
- Session persistence
- Logout
- Authentication rate limiting

Authentication API documentation is available in:

`docs/api/authentication.md`

Security details are documented in:

`docs/security/authentication.md`

## File Storage

PeakCloud supports authenticated private file storage backed by S3-compatible object storage.

Current file capabilities include:

- authenticated file upload
- PostgreSQL metadata persistence
- MinIO/S3-compatible binary storage
- private per-user ownership
- file listing
- metadata retrieval
- secure downloads
- file deletion
- drag-and-drop browser uploads
- upload-size enforcement
- generated object keys
- rollback when metadata persistence fails
