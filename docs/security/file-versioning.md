# File Versioning Security

File version history uses PeakCloud's centralized authorization system.

## Authorization

Read authorization applies to:

- listing versions
- retrieving version metadata
- reading historical content

Download authorization applies to historical downloads.

Edit authorization applies to:

- creating versions
- restoring versions

Inherited folder sharing therefore continues to apply to version history.

## Immutable Historical Objects

Historical objects use unique object keys.

New uploads never overwrite historical objects.

Restoring a version creates another object rather than modifying the source version.

## Database Integrity

The `file_versions` table enforces:

- positive version numbers
- non-negative sizes
- non-empty object keys
- non-empty content types
- unique `(file_id, version_number)` pairs
- unique object keys

## Existing Files

Migration `000005_create_file_versions` backfills existing files as version 1.

## Restore Safety

Restore is append-only.

If versions 1, 2, and 3 exist, restoring version 1 creates version 4.

Versions 1 through 3 remain unchanged.
