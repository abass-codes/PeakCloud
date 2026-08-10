# ADR 0002: Use Server-Side Sessions for Browser Authentication

## Status

Accepted

## Context

PeakCloud requires authentication for its browser-based application.

Authentication credentials must support secure login persistence, logout, session expiration, and future session revocation.

## Decision

PeakCloud will use opaque server-side session tokens.

A cryptographically random token is sent to the browser in an HttpOnly cookie.

PeakCloud stores only a SHA-256 hash of the session token in PostgreSQL.

Passwords are hashed using Argon2id.

## Consequences

### Positive

- Sessions can be revoked server-side.
- Logout immediately invalidates the persisted session.
- Authentication tokens are not exposed to application JavaScript.
- Raw session tokens are not stored in the database.
- Session behavior is straightforward to reason about.

### Negative

- Authentication requires a server-side session lookup.
- Session storage must scale with active users.
- Distributed deployments require shared session storage.

A future architecture may move session storage to Redis or another distributed session store without changing the browser-facing authentication model.
