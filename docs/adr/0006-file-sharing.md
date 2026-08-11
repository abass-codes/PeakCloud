# ADR 0006: File Sharing Model

## Status

Accepted

## Context

PeakCloud needs private user-to-user sharing and revocable public links while keeping object storage private.

## Decision

PeakCloud uses direct resource shares and public share links.

Direct shares reference an owner, recipient, resource, permission, and download policy.

Public links use cryptographically random bearer tokens.

Only token hashes are stored.

Public links may additionally contain:

- password protection
- expiration
- download restrictions
- revocation state

Public file delivery remains application mediated rather than exposing object-storage URLs.

## Authorization boundary

Feature 6 persists sharing grants and protects public-link access.

Feature 7 introduces centralized authorization that consumes these grants for authenticated file and folder operations.

## Consequences

Share lifecycle and object storage remain independent.

Tokens can be revoked without moving stored objects.

Object-storage keys remain private.

The permission model can be reused by the centralized authorization layer.
