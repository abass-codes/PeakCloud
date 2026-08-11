# File Sharing API

PeakCloud supports direct user-to-user resource sharing and public share links.

## Direct shares

### Create

POST /api/v1/shares

A share contains:

- recipient email
- file or folder resource
- viewer or editor permission
- download permission

### List owned shares

GET /api/v1/shares

### Shared With Me

GET /api/v1/shares/shared-with-me

### Update

PATCH /api/v1/shares/:id

### Revoke

DELETE /api/v1/shares/:id

## Public links

### Create

POST /api/v1/share-links

Supports:

- viewer/editor metadata
- download controls
- optional password
- optional expiration

### List

GET /api/v1/share-links

### Revoke

DELETE /api/v1/share-links/:id

### Resolve

POST /api/v1/public/shares/:token

### Public file content

GET /api/v1/public/shares/:token/content

### Public file download

GET /api/v1/public/shares/:token/download

Password-protected content and downloads require the share password.

PeakCloud stores only SHA-256 hashes of public share tokens.
