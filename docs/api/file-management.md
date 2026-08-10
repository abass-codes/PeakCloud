# File and Folder Management API

PeakCloud provides authenticated filesystem-style organization for stored objects.

## Drive Contents

`GET /api/v1/drive`

Optional query:

`folder_id=<uuid>`

Returns the current folder's child folders, files, and breadcrumb path.

## Create Folder

`POST /api/v1/folders`

Request:

```json
{
  "name": "Documents",
  "parent_id": null
}
```

Nested folders use the parent folder UUID.

## Rename Folder

`PATCH /api/v1/folders/:id/name`

## Move Folder

`PATCH /api/v1/folders/:id/location`

PeakCloud rejects moves that would place a folder inside itself or one of its descendants.

## Delete Folder

`DELETE /api/v1/folders/:id`

Deletion is recursive.

Binary objects belonging to files in the deleted folder tree are removed from object storage before relational metadata is removed.

## Folder Breadcrumbs

`GET /api/v1/folders/:id/breadcrumbs`

Returns the hierarchy from the top-level folder through the requested folder.

## Folder-Aware Upload

`POST /api/v1/files`

Multipart fields:

- `file`
- `folder_id` — optional

When `folder_id` is omitted, the file is stored at the drive root.

## Rename File

`PATCH /api/v1/files/:id/name`

## Move File

`PATCH /api/v1/files/:id/location`

## Copy File

`POST /api/v1/files/:id/copy`

The binary object is copied to a newly generated object key.

## Bulk Move

`POST /api/v1/drive/bulk/move`

## Bulk Delete

`POST /api/v1/drive/bulk/delete`

## Bulk Download

`POST /api/v1/drive/bulk/download`

Selected files are streamed as a ZIP archive.

## Authorization

Every operation is scoped to the authenticated user's ID.

Resources belonging to another account are treated as nonexistent and return `404 Not Found` where applicable.
