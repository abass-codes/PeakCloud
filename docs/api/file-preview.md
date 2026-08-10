# File Preview API

PeakCloud provides authenticated browser previews for supported file types.

## Preview Metadata

```http
GET /api/v1/files/:id/preview
```

Returns file metadata and the preview classification.

Preview kinds include:

- `image`
- `pdf`
- `text`
- `code`
- `video`
- `audio`
- `unsupported`

The endpoint is owner-scoped.

Files belonging to another account are treated as nonexistent.

## Preview Content

```http
GET /api/v1/files/:id/content
```

Streams supported file content for inline browser rendering.

The response uses:

```http
Content-Disposition: inline
X-Content-Type-Options: nosniff
Cache-Control: private, no-store
```

## Supported Content

PeakCloud supports browser previews for:

- images
- PDFs
- plain text
- common source-code files
- browser-supported video
- browser-supported audio

Unsupported files fall back to the normal download workflow.

## Text Preview Limits

Large text and source-code files are not rendered inline.

PeakCloud limits text-style previews to 1 MiB to prevent large objects from being loaded into the browser as text.

## Authorization

Preview metadata and preview content require an authenticated PeakCloud session.

Every file lookup is scoped by both file ID and authenticated owner ID.

A user requesting another user's file receives `404 Not Found` rather than an authorization response that would reveal the object's existence.

## Object Storage

Preview endpoints do not expose internal object keys.

The backend resolves the authenticated file record and streams the corresponding object from S3-compatible storage.
