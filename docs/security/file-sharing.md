# File Sharing Security

## Private by default

Files and folders remain private unless their owner explicitly creates a share.

## Ownership verification

A user can create a share only for a resource they own.

Unauthorized ownership checks do not reveal private resource metadata.

## Direct shares

Direct shares reference registered PeakCloud users and contain:

- viewer/editor permission
- download permission

Owners can update or revoke direct shares.

## Public tokens

Public links use 256 bits of cryptographically secure randomness.

Only SHA-256 hashes of public tokens are persisted.

The raw token is returned only when a link is created.

## Passwords

Optional public-link passwords are hashed using bcrypt.

Plaintext passwords are never stored.

## Expiration

Expired public links are rejected by the backend.

## Revocation

Revoked public links cannot be resolved or used to retrieve content.

## Downloads

Public download endpoints enforce allow_download before serving a file.

## Object storage

Public links never expose PeakCloud object-storage keys.

All public content continues to pass through the application backend.

## Authorization boundary

Feature 6 creates and manages sharing grants.

Feature 7 provides centralized authorization enforcement across authenticated file and folder operations.
