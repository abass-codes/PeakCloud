# Trash and Recovery API

All trash endpoints require authentication.

## List Trash

`GET /api/v1/trash`

Returns the authenticated user's trashed files and folders.

## Move Resource to Trash

`POST /api/v1/trash/:type/:id`

Supported resource types:

- `file`
- `folder`

The resource must belong to the authenticated user.

## Restore Resource

`POST /api/v1/trash/:type/:id/restore`

Restores a resource currently in trash.

## Permanently Delete Resource

`DELETE /api/v1/trash/:type/:id`

Permanently deletes a resource that is already in trash.

Active resources cannot be permanently deleted through this
endpoint.

## Errors

Trash endpoints may return:

- `400 Bad Request`
- `401 Unauthorized`
- `403 Forbidden`
- `404 Not Found`
- `409 Conflict`
- `500 Internal Server Error`
