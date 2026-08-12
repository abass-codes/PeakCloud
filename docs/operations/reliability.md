# Production Reliability

## HTTP timeouts

PeakCloud configures HTTP server timeouts through environment variables:

- `HTTP_READ_HEADER_TIMEOUT_SECONDS`
- `HTTP_READ_TIMEOUT_SECONDS`
- `HTTP_WRITE_TIMEOUT_SECONDS`
- `HTTP_IDLE_TIMEOUT_SECONDS`
- `SHUTDOWN_TIMEOUT_SECONDS`

The default values are documented in `.env.example`.

Timeout values must be positive integers.

## Liveness

`GET /live`

The liveness endpoint reports whether the API process is running and able to
serve HTTP requests.

A successful response returns HTTP 200.

Liveness intentionally does not depend on PostgreSQL, Redis, or object
storage availability.

## Health

`GET /health`

The health endpoint performs the application's dependency-aware health
checks.

This endpoint is appropriate for determining whether the service and its
required dependencies are ready for normal operation.

## Graceful shutdown

When PeakCloud receives a termination signal, the API starts graceful
shutdown using the configured shutdown timeout.

This provides active requests an opportunity to complete while still placing
an upper bound on shutdown duration.

## Production operation

Production deployments should:

1. Use secure production configuration.
2. Configure infrastructure probes against `/live` and `/health` according
   to their intended semantics.
3. Set timeout values appropriate for the deployment environment.
4. Keep upload limits separate from small JSON request-body limits.
5. Monitor rate-limit responses and request IDs during incident analysis.
