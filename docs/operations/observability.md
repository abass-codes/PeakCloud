# PeakCloud Observability

This document describes PeakCloud's application observability capabilities.

## Structured Logging

PeakCloud uses structured HTTP request logging.

The API creates the application logger during startup and passes that logger into the HTTP router.

Each completed request produces an `http_request` event containing:

- method
- path
- status
- request_id
- latency_ms
- client_ip

## Request IDs

PeakCloud uses the `X-Request-ID` header for request correlation.

When a request already contains an ID, PeakCloud preserves it. Otherwise, PeakCloud generates one.

The request ID is stored in the Gin context, returned in the response header, and included in structured request logs.

## HTTP Metrics

PeakCloud exposes HTTP request metrics through:

GET /metrics

The metrics middleware records request activity as traffic passes through the application.

## Operational Endpoints

PeakCloud provides:

GET /live
GET /health
GET /metrics

The liveness endpoint confirms the API process is responding.

The health endpoint checks configured application dependencies.

The metrics endpoint exposes HTTP request metrics.

## CI as an Operational Control

Continuous integration validates changes automatically before they enter the production codebase.

PeakCloud separates CI into backend, frontend, and integration workflows so failures can be identified at the appropriate layer.

## Incident Investigation

For an HTTP failure, an operator can capture the X-Request-ID, search structured logs for the corresponding request_id, inspect the request method and path, inspect the response status and latency, and compare the event with health and metrics information.

## Security Considerations

Observability data should not intentionally record secrets, authentication credentials, session cookies, or sensitive request bodies.

Request logging is limited to operational metadata required for debugging and monitoring.
