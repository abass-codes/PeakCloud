# ADR 0008: Immutable File Versioning

## Status

Accepted

## Context

PeakCloud originally stored one active object for each logical file.

Replacing content would otherwise destroy the previous state.

Feature 8 requires historical versions and recovery.

## Decision

PeakCloud separates logical file identity from immutable content versions.

The `files` table remains the logical file.

The `file_versions` table stores immutable content snapshots.

Each version stores:

- logical file ID
- version number
- object key
- size
- content type
- ETag
- creator
- creation time

The active `files.object_key` points to current content.

## Version Creation

Uploading replacement content creates a new object and version row.

Historical objects are not overwritten.

## Restoration

Restoration is append-only.

For example:

`v1 -> v2 -> v3 -> restore v1 -> v4`

Version 4 contains restored content from version 1.

Versions 1 through 3 remain unchanged.

## Authorization

Version operations use PeakCloud's centralized authorization service.

This preserves direct-share and inherited-folder authorization behavior.

## Migration

Existing files are backfilled as version 1 by migration `000005_create_file_versions`.

## Consequences

Benefits include:

- immutable history
- recoverable content
- stable logical file IDs
- auditable restoration
- consistent authorization

Costs include:

- increased storage consumption
- additional database metadata
- cross-system object-storage/database operations
