# CI/CD Testing & Observability

Feature 11 introduces automated continuous integration and application-level observability for PeakCloud.

## Goals

The feature establishes automated validation for backend, frontend, and integration changes while improving visibility into HTTP request behavior.

## Continuous Integration

PeakCloud includes separate GitHub Actions workflows for:

- backend validation
- frontend validation
- integration validation

The workflows are stored in `.github/workflows`.

### Backend CI

The backend workflow validates the Go application using automated tests, static analysis, and formatting checks.

### Frontend CI

The frontend workflow validates the Next.js application using linting and a production build.

### Integration CI

The integration workflow validates application behavior that depends on PeakCloud's supporting infrastructure.

## Structured Request Logging

PeakCloud includes HTTP request logging middleware.

Each completed HTTP request records structured fields including:

- HTTP method
- request path
- response status
- request ID
- request latency
- client IP

The request logger uses the same application logger created during API startup.

## Request Correlation

Request IDs are integrated into structured HTTP request logs.

This allows requests to be correlated across incoming HTTP traffic, application logs, debugging sessions, and operational investigations.

## HTTP Metrics

PeakCloud includes HTTP metrics middleware that tracks application request activity.

Metrics are exposed through:

GET /metrics

## Testing

Feature 11 includes tests for structured request logging, HTTP metrics collection, and metrics endpoint behavior.

## Operational Value

Feature 11 provides repeatable CI validation, automated regression detection, structured HTTP logs, request correlation, application request metrics, and an operational metrics endpoint.
