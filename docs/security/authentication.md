# Authentication Security

## Passwords

PeakCloud hashes passwords with Argon2id using a unique cryptographically random salt.

Plaintext passwords are never persisted.

## Sessions

PeakCloud creates cryptographically random opaque session tokens.

Only SHA-256 hashes of session tokens are stored in PostgreSQL.

Authentication cookies are:

- HttpOnly
- SameSite=Lax
- Secure in production
- Scoped to the application path

## Authentication Errors

Login failures return a generic invalid-credentials response instead of revealing whether a particular email exists.

## Rate Limiting

Registration and login endpoints are rate limited.

The initial implementation uses an in-process limiter. A distributed implementation can later use Redis.

## Input Validation

Registration validates:

- Email format
- Display-name length
- Password length

Database constraints provide an additional integrity boundary.

## Future Security Work

Later milestones will include:

- CSRF hardening as deployment topology evolves
- Email verification
- Password reset
- Session management UI
- Account recovery
- Audit logging
- Distributed rate limiting
- Security headers
- Abuse monitoring
