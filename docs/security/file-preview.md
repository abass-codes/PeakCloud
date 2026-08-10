# File Preview Security

PeakCloud file previews preserve the same private ownership model used by file downloads.

## Authentication

Preview endpoints require a valid authenticated PeakCloud session.

## Ownership

File records are resolved using both:

- file ID
- authenticated owner ID

A file owned by another account is treated as nonexistent.

## Object Storage Isolation

PeakCloud does not expose:

- internal object keys
- bucket credentials
- storage paths
- MinIO administrative endpoints

The API resolves and streams objects on behalf of authenticated users.

## MIME Handling

Preview classification uses stored MIME metadata and controlled filename-extension rules.

Inline responses include:

`X-Content-Type-Options: nosniff`

This prevents browsers from freely MIME-sniffing preview responses.

## Content Disposition

Preview content uses inline content disposition.

User-controlled filenames are sanitized before being placed into response headers.

## Caching

Private preview responses use:

`Cache-Control: private, no-store`

## Large Text Files

Text and source-code previews are size limited.

Files exceeding the configured preview threshold must be downloaded instead of rendered as text.

## Unsupported Content

Unsupported file types are not forced into an inline browser renderer.

Users are directed to the existing authenticated download path instead.
