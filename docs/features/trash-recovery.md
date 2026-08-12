# Trash and Recovery

PeakCloud uses soft deletion for files and folders.

Deleting a resource from the normal drive does not immediately
destroy its metadata. Instead, PeakCloud records a `deleted_at`
timestamp and removes the resource from active drive queries.

## Capabilities

- Move files to trash
- Move folders to trash
- List trashed resources
- Restore files
- Restore folders
- Permanently delete trashed resources
- Bulk trash actions in the web application
- Dedicated trash interface
- Active-resource isolation
- Owner-scoped trash operations

## Active Resource Isolation

Normal file and folder operations exclude resources whose
`deleted_at` value is not null.

This prevents trashed resources from appearing in normal drive
navigation or being modified through standard active-resource
operations.

## Recovery

Restoring a resource clears its `deleted_at` timestamp and returns
it to active storage.

## Permanent Deletion

Permanent deletion is only allowed for resources that are already
in trash.

This prevents the permanent-delete endpoint from bypassing the
normal trash lifecycle.
