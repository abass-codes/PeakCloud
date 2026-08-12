# ADR 0009: Soft Delete Trash and Recovery

## Status

Accepted

## Context

Immediate hard deletion provides poor recovery guarantees and makes
accidental deletion difficult to reverse.

PeakCloud needs a recoverable deletion workflow similar to modern
cloud-storage systems.

## Decision

Files and folders use nullable `deleted_at` timestamps.

A null value means the resource is active.

A non-null value means the resource is in trash.

Normal repository operations exclude trashed resources.

Trash operations explicitly operate on deleted resources.

Permanent deletion is only permitted after soft deletion.

## Consequences

Advantages:

- accidental deletions can be recovered
- active and deleted resources have explicit lifecycle states
- destructive operations require an additional step
- trash can be exposed through a dedicated API and UI

Tradeoffs:

- repository queries must consistently enforce active-resource
  filtering
- permanent deletion requires careful metadata and object-storage
  cleanup
- folder relationships must account for deleted ancestors
