# File Versioning API

PeakCloud exposes authenticated endpoints for immutable file version history.

## Create a Version

`POST /api/v1/files/:id/versions`

Uploads new content for an existing logical file.

The caller must have edit permission.

The request uses multipart form data with the file in the `file` field.

A successful request creates a new immutable object, stores version metadata, and updates the logical file to point to the new object.

## List Versions

`GET /api/v1/files/:id/versions`

Returns version history for a logical file.

The caller must have read access.

## Get Version Metadata

`GET /api/v1/files/:id/versions/:version`

Returns metadata for a specific version.

## Read Historical Content

`GET /api/v1/files/:id/versions/:version/content`

Returns stored content for a historical version.

## Download Historical Content

`GET /api/v1/files/:id/versions/:version/download`

Downloads a historical version.

Download authorization is enforced separately from basic read access.

## Restore a Version

`POST /api/v1/files/:id/versions/:version/restore`

Restores historical content without mutating existing history.

If versions 1 through 3 exist, restoring version 1 creates version 4 containing a copy of version 1's content.

The caller must have edit permission.

## Errors

Version endpoints may return:

- `400 Bad Request`
- `401 Unauthorized`
- `403 Forbidden`
- `404 Not Found`
- `413 Request Entity Too Large`
- `500 Internal Server Error`
