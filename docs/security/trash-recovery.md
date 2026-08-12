# Trash and Recovery Security

Trash operations are authenticated and owner scoped.

## Authorization

A user cannot trash, restore, inspect, or permanently delete another
user's private resources through the trash API.

Repository operations include owner identifiers when mutating trash
state.

## Active Resource Isolation

Resources with a non-null `deleted_at` timestamp are excluded from
normal active-resource queries.

This prevents soft-deleted resources from continuing to behave as
ordinary drive objects.

## Permanent Deletion

Permanent deletion requires the resource to already have a
`deleted_at` timestamp.

This creates a two-stage destructive workflow:

1. Move to trash.
2. Permanently delete from trash.

## Recovery

Restore operations only target resources currently marked as
deleted.
