# Production Security & Reliability

Feature 10 hardens PeakCloud for production-style operation by adding
security middleware, abuse protection, production configuration validation,
liveness checks, and configurable HTTP server reliability controls.

## Security foundation

PeakCloud now provides:

- Request IDs for request correlation.
- Production security headers.
- Production startup configuration validation.
- Validation for secure session cookies.
- Validation for secure object storage.
- Validation for an HTTPS web URL in production.

The application fails startup validation when production security
requirements are not satisfied.

## Abuse protection

PeakCloud applies rate limiting to API traffic with separate protection for
authentication and public sharing routes.

Authentication routes receive stricter rate limiting because login and
registration endpoints are more sensitive to automated abuse.

Public sharing routes are independently rate limited because they are
accessible without an authenticated session.

Small JSON-oriented authentication and public requests are protected by a
request body limit.

The general API group intentionally does not receive the same 2 MB body
limit because file uploads use a separate configurable maximum upload size.

## Reliability

HTTP server behavior is controlled through environment-backed reliability
configuration.

Configured controls include:

- Read-header timeout.
- Read timeout.
- Write timeout.
- Idle timeout.
- Graceful shutdown timeout.

PeakCloud also exposes:

- `/live` for process liveness.
- `/health` for dependency-aware health checks.

This separation allows infrastructure to distinguish a running process from
an application whose dependencies may be unavailable.

## Verification

Feature 10 includes tests covering:

- Production security validation.
- Request ID behavior.
- Security headers.
- Rate limiting.
- Request body limits.
- Reliability configuration.

The complete repository is also validated with Go tests, Go vet, frontend
linting, the frontend production build, formatting checks, and the project
`make check` workflow.
