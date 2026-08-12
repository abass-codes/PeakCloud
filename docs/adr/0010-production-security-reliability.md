# ADR 0010: Production Security and Reliability Boundaries

## Status

Accepted

## Context

PeakCloud requires production-oriented safeguards around HTTP requests,
configuration, abuse prevention, and server lifecycle behavior.

The application handles both small API requests and potentially large file
uploads. Applying identical limits to every request would either leave small
endpoints insufficiently protected or incorrectly restrict legitimate file
transfers.

Production deployments also require predictable server timeouts, graceful
shutdown, and separate liveness and dependency-health signals.

## Decision

PeakCloud will:

1. Validate security-sensitive configuration during production startup.
2. Install request ID and security-header middleware globally.
3. Apply a general API rate limiter.
4. Apply dedicated rate limiting to authentication and public sharing
   traffic.
5. Apply small request-body limits only where large uploads are not expected.
6. Keep file-upload sizing governed by the upload-specific configuration.
7. Configure HTTP server and shutdown timeouts through environment-backed
   reliability settings.
8. Expose `/live` independently from the dependency-aware `/health`
   endpoint.

## Consequences

### Positive

- Unsafe production configuration fails early.
- Request tracing is easier through request IDs.
- Common abuse paths receive explicit protection.
- File uploads are not accidentally restricted by a small global body limit.
- Server lifecycle behavior is configurable and bounded.
- Infrastructure can distinguish liveness from dependency health.

### Trade-offs

- In-memory rate limiting is local to an API process and does not provide a
  globally shared quota across multiple replicas.
- Production operators must configure security and reliability environment
  variables correctly.
- Additional middleware and configuration increase the number of operational
  controls that must be maintained.

## Future considerations

A multi-instance deployment may move rate-limit state to shared
infrastructure such as Redis and add production metrics around rejected
requests, latency, and graceful shutdown behavior.
