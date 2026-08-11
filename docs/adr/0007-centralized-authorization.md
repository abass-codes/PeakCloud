# ADR 0007: Centralized Authorization

## Status

Accepted

## Context

PeakCloud originally protected files and folders primarily through ownership checks.

Feature 6 introduced direct sharing grants with viewer and editor permissions.

Continuing to implement authorization independently inside each file and folder operation would duplicate security logic and increase the risk of inconsistent permission enforcement.

## Decision

PeakCloud uses a centralized authorization package for authenticated resource access.

The authorization layer resolves:

- resource ownership
- direct grants
- inherited folder grants
- effective permissions
- download policy

File and folder services request authorization decisions from this layer before performing protected operations.

Folder grants are inherited by descendant folders and files.

Effective permissions follow this priority:

    owner > editor > viewer > no access

Owners retain exclusive delete authority.

Move and copy operations must also validate the destination authorization boundary.

## Consequences

Authorization behavior is consistent across file and folder operations.

Sharing grants can be reused without duplicating authorization logic.

Folder sharing applies naturally to descendant resources.

Revoking a sharing grant removes the access derived from that grant.

Future permission features can extend a single authorization layer rather than modifying every resource handler independently.

## Drive operations

Drive endpoints reuse the centralized authorization service rather than
implementing separate ownership rules.

Bulk operations use a preflight model. Authorization is evaluated for the
complete batch before mutation begins.

This provides two important guarantees:

1. file, folder, and drive operations use the same permission semantics;
2. a batch containing an unauthorized resource is rejected before authorized
   resources earlier in the request are modified.

Bulk deletion remains owner-only. Bulk download respects inherited
`allow_download` policy. Bulk move requires edit access to both the selected
resources and the destination folder.
