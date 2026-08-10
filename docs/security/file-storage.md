# File Storage Security

PeakCloud file storage is private by default.

## Ownership

Every file metadata record contains an `owner_id` referencing the authenticated PeakCloud user.

File queries are scoped by both file ID and owner ID.

This prevents one authenticated account from retrieving another account's file by guessing its UUID.

## Object Keys

PeakCloud does not use the original filename as the physical object key.

Object keys are generated server-side using:

- authenticated user namespace
- date namespace
- random UUID

Internal object keys are not returned through the public API.

## Upload Validation

PeakCloud:

- rejects blank filenames
- rejects path traversal filenames
- enforces a maximum upload size
- stores content type metadata
- does not trust filenames as storage paths

## Downloads

Downloads require an authenticated session.

PeakCloud first verifies file ownership in PostgreSQL and only then resolves the corresponding object.

## Deletion

Deletion requires authenticated ownership.

PeakCloud removes the binary object and associated metadata.

## Future Hardening

Future storage features may add:

- malware scanning
- content sniffing
- file extension policies
- per-account storage quotas
- upload checksums
- multipart upload controls
- object encryption policies
- signed download URLs
- audit logging
