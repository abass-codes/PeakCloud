# Centralized Authorization

PeakCloud uses a centralized authorization layer for authenticated file and folder operations.

## Permission model

PeakCloud supports three effective permission levels:

- owner
- editor
- viewer

Owners have full control over their resources.

Editors can read and modify shared resources but cannot delete resources owned by another user.

Viewers can read shared resources but cannot modify them.

Downloads are additionally controlled by the sharing grant's `allow_download` policy.

## Centralized resolution

Authorization decisions are resolved through the `internal/authorization` package rather than being implemented independently by individual file and folder handlers.

The authorization layer determines:

1. resource ownership
2. direct sharing grants
3. inherited folder grants
4. effective permission
5. download policy

## Folder inheritance

Permissions granted to a folder apply to resources beneath that folder.

For example:

    Documents/
    └── School/
        └── report.txt

A viewer or editor grant on `Documents` can authorize access to `School` and `report.txt`.

This avoids requiring a separate sharing row for every descendant resource.

## Effective permissions

When multiple applicable grants exist, PeakCloud selects the strongest effective permission:

    owner > editor > viewer > no access

Applicable download permissions are also resolved as part of the effective access decision.

## Private resource protection

Unauthorized reads return resource-not-found responses rather than exposing whether another user's private resource exists.

This prevents resource enumeration through authorization errors.

## File authorization

File operations enforce centralized authorization for:

- metadata reads
- content and preview access
- downloads
- rename
- move
- copy
- delete

Delete remains owner-only.

## Folder authorization

Folder operations enforce centralized authorization for:

- reads
- listing
- creation within authorized folders
- rename
- move
- breadcrumbs
- delete

Delete remains owner-only.

## Move boundaries

Move operations verify authorization for both the source resource and destination folder.

This prevents a shared user from moving a resource outside the authorization boundary.

## Revocation

Authorization is evaluated against current sharing grants.

When a share is revoked, inherited access from that grant disappears.

## Drive and bulk-operation authorization

Drive operations use the same centralized authorization policy as individual
file and folder operations.

Bulk operations perform an authorization preflight before mutation begins.

### Bulk move

Every selected resource must grant edit access to the requesting user.
When a destination folder is supplied, the destination must also grant edit
access.

The entire authorization preflight completes before the first resource is
moved.

### Bulk delete

Deletion remains owner-only. Every selected resource must be owned by the
requesting user before the first deletion begins.

Editors cannot delete resources owned by another user.

### Bulk download

Every selected file must grant download access before the ZIP response is
written. Recipient access therefore respects the `allow_download` policy
inherited from direct or ancestor-folder grants.

### Failure behavior

PeakCloud rejects a bulk operation when any selected resource fails its
authorization requirement. This prevents an unauthorized item later in a
request from causing a partially authorized batch to execute.
