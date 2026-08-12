# Production Security

## Production configuration validation

PeakCloud validates security-sensitive configuration before serving
production traffic.

Production startup requires secure settings for:

- Session cookies.
- Object storage transport.
- The configured web application URL.

Invalid production configuration causes startup to fail instead of allowing
the service to run with known insecure settings.

## Security headers

The HTTP router installs security-header middleware globally so responses
receive the application's production security headers consistently.

## Request IDs

Requests receive an identifier for correlation and diagnostics.

When an incoming request already contains an accepted request identifier,
PeakCloud preserves it. Otherwise, the middleware generates one.

## Rate limiting

Rate limiting is applied at multiple boundaries:

- General API traffic.
- Authentication traffic.
- Public sharing traffic.

Authentication endpoints receive a dedicated limiter to reduce exposure to
automated login and registration abuse.

Public sharing endpoints receive their own limiter because they can be
reached without an authenticated session.

## Request body limits

PeakCloud limits small request bodies where a large payload is not expected.

The 2 MB middleware limit is intentionally not installed globally on
`/api/v1`. A global limit would conflict with legitimate file uploads.

File uploads continue to use the application's configurable upload-size
limit.

## Deleted-resource isolation

Authorization and repository operations preserve the deleted-resource
isolation introduced by Trash & Recovery. Soft-deleted resources are not
treated as active resources during normal access checks.

## Operational principle

Security controls should fail safely while preserving legitimate product
behavior. Abuse protection therefore targets appropriate request classes
without imposing a small global request limit on file-transfer operations.
