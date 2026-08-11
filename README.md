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

## File and Folder Management

PeakCloud provides hierarchical filesystem-style organization on top of private object storage.

Current management capabilities include:

- root drive
- nested folders
- breadcrumb navigation
- folder creation
- file and folder rename
- file and folder move APIs
- file copying
- recursive folder deletion
- folder cycle prevention
- folder-aware uploads
- sorting and filtering
- multi-item selection
- bulk deletion
- bulk file download as ZIP
- account-isolated file and folder access
- stable generated object-storage keys

## File Preview

PeakCloud supports authenticated in-browser previews for common file formats.

Current preview capabilities include:

- image preview
- PDF preview
- plain-text preview
- source-code preview
- browser-supported video playback
- browser-supported audio playback
- unsupported-format download fallback
- owner-scoped preview authorization
- private inline content streaming
- text-preview size limits

## File Sharing

PeakCloud supports secure file and folder sharing through:

- private user-to-user shares
- viewer and editor permissions
- per-share download controls
- Shared With Me
- permission updates and revocation
- public share links
- optional link passwords
- expiration timestamps
- public-link revocation
- public file preview/content delivery
- controlled public downloads

Public share tokens use cryptographically secure randomness and are persisted only as hashes. Password-protected public links use bcrypt password hashing. Object-storage keys remain private.

See:

- `docs/api/file-sharing.md`
- `docs/security/file-sharing.md`
- `docs/adr/0006-file-sharing.md`

## Centralized Authorization

PeakCloud uses centralized authorization for authenticated file and folder operations.

The authorization system supports:

- owner, editor, and viewer permissions
- direct resource grants
- inherited folder permissions
- download-policy enforcement
- owner-only destructive operations
- source and destination checks for resource movement
- private-resource enumeration protection
- immediate enforcement of share revocation

Effective permissions are resolved using:

`owner > editor > viewer > no access`

See:

- `docs/security/authorization.md`
- `docs/adr/0007-centralized-authorization.md`

## File Versioning and Recovery

PeakCloud preserves immutable historical versions of file content while maintaining a stable logical file identity.

Current versioning capabilities include:

- automatic version 1 backfill for existing files
- immutable historical file versions
- unique object-storage keys for every version
- authenticated version history
- historical version metadata
- historical content streaming
- historical version downloads
- new-version uploads
- restore without destructive overwrite
- centralized read, download, and edit authorization
- inherited folder-sharing authorization
- viewer/editor permission enforcement
- per-share download enforcement

Restoring an older version creates a new version rather than rewriting history.

For example:

    v1
    v2
    v3
    v4 <- restored copy of v1

See:

- `docs/features/file-versioning.md`
- `docs/api/file-versioning.md`
- `docs/security/file-versioning.md`
- `docs/adr/0008-file-versioning.md`
