# File and Folder Management Security

## Ownership

Every folder and file belongs to a PeakCloud user.

Backend operations scope resource lookup using the authenticated user ID.

Client-side visibility is never treated as authorization.

## Resource Enumeration

Requests for another user's files or folders do not reveal ownership information.

Where appropriate, inaccessible resources return `404 Not Found`.

## Folder Relationships

PeakCloud validates hierarchy changes before updating folder relationships.

Folders cannot be moved:

- into themselves
- into their descendants
- into another user's folders

This prevents cyclic folder graphs and cross-account hierarchy manipulation.

## Names

File and folder names are validated before persistence.

Path traversal syntax is rejected.

User-visible names do not control physical object-storage paths.

## Object Keys

Binary objects continue using generated server-side object keys.

Renaming or moving a file changes metadata rather than allowing users to manipulate storage keys.

## Recursive Deletion

Before folder metadata is recursively removed, PeakCloud identifies associated binary objects and removes them from object storage.

This reduces orphaned objects after folder deletion.

## Bulk Operations

Bulk operations perform the same owner-scoped service operations used by individual resource endpoints.

Selection in the frontend does not bypass backend authorization.
