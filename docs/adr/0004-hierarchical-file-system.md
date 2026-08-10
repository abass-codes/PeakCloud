# ADR 0004: Hierarchical File and Folder Model

## Status

Accepted

## Context

PeakCloud needs filesystem-style organization while storing binary objects independently from relational metadata.

Object storage itself should not be responsible for representing the user-visible folder hierarchy.

## Decision

PeakCloud represents folders relationally in PostgreSQL.

Each folder contains:

- an owner
- an optional parent folder
- a name
- timestamps

Files contain an optional `folder_id`.

A null parent or folder ID represents the user's drive root.

Binary object keys remain independent of folder names and file names.

## Hierarchy

Nested folders use a self-referencing `parent_id`.

PeakCloud validates folder moves before persistence.

A folder cannot:

- become its own parent
- move into one of its descendants
- move into a folder owned by another user

## Object Storage

Moving or renaming a file does not move its binary object.

The user-visible hierarchy exists in PostgreSQL metadata while immutable generated object keys remain stable.

This avoids unnecessary object-storage copies for metadata-only operations.

## Recursive Deletion

Folder deletion requires coordination between PostgreSQL and object storage.

PeakCloud identifies files in the folder subtree, removes their binary objects, and then removes the folder hierarchy and metadata.

## Consequences

### Positive

- Efficient rename operations
- Efficient move operations
- Stable object keys
- Relational hierarchy queries
- User ownership enforcement
- Object storage remains independent from presentation paths

### Negative

- Recursive operations require additional logic
- Cross-system deletion requires failure handling
- Deep hierarchy operations may require recursive database queries
