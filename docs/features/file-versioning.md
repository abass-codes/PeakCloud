# File Versioning and Recovery

PeakCloud preserves immutable historical versions of file content while keeping the `files` table as the logical identity of each file.

## Architecture

A logical file can have many immutable versions:

```text
files
  |
  +-- file_versions
        |
        +-- version 1
        +-- version 2
        +-- version 3
        +-- ...
```

The current `files.object_key` points to the active object.

Each `file_versions` row stores an immutable snapshot containing:

- version number
- object key
- size
- content type
- ETag
- creator
- creation timestamp

## Existing Files

Migration `000005_create_file_versions` backfills every existing file as version 1.

This ensures files created before Feature 8 immediately participate in version history.

## Creating Versions

Authorized editors can create a new version of an existing logical file.

New content is stored under a unique object key.

Example:

```text
<owner-id>/versions/<file-id>/v2-<uuid>
```

The historical object is never overwritten.

After the version is persisted, the logical file points to the newest object.

## Version History

PeakCloud supports:

- listing versions
- retrieving version metadata
- reading historical content
- downloading historical content
- restoring historical versions

Version history is ordered by version number.

## Restore Semantics

Restoring a historical version does not overwrite history.

For example:

```text
v1
v2
v3
```

Restoring `v1` creates:

```text
v1
v2
v3
v4 <- restored copy of v1
```

The active logical file then points to the new v4 object.

This preserves immutable history and provides an audit trail of recovery operations.

## Authorization

Version operations use PeakCloud's centralized authorization service.

Read access is required for:

- listing versions
- retrieving version metadata
- reading historical content

Download permission is required for historical downloads.

Edit permission is required for:

- creating versions
- restoring versions

Inherited folder sharing therefore applies consistently to version history.

## API

```text
POST /api/v1/files/:id/versions
GET  /api/v1/files/:id/versions
GET  /api/v1/files/:id/versions/:version
GET  /api/v1/files/:id/versions/:version/content
GET  /api/v1/files/:id/versions/:version/download
POST /api/v1/files/:id/versions/:version/restore
```

## Storage Guarantees

Historical version objects use unique object keys.

Restoring a version creates a new object instead of changing the original historical object.

Database uniqueness constraints prevent duplicate version numbers for the same logical file.
