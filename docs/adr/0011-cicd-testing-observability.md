# ADR 0011: CI/CD Testing and Observability

## Status

Accepted

## Context

As PeakCloud grew, manual validation alone became insufficient.

The project needed automated checks capable of detecting regressions across the backend, frontend, and infrastructure-dependent application paths.

The API also needed better runtime visibility through structured request logging and HTTP metrics.

## Decision

PeakCloud will use GitHub Actions for continuous integration.

CI is separated into three workflows:

1. backend CI
2. frontend CI
3. integration CI

Backend CI validates the Go application.

Frontend CI validates the Next.js application.

Integration CI validates behavior that depends on PeakCloud's supporting services.

PeakCloud will also provide HTTP observability through:

- request ID middleware
- structured request logging middleware
- HTTP metrics middleware

The existing application logger is injected into the router rather than creating a second logger.

HTTP request logs include method, path, status, request ID, latency, and client IP.

HTTP metrics are exposed through:

GET /metrics

## Consequences

### Positive

- regressions can be detected automatically
- backend and frontend validation are repeatable
- integration behavior receives dedicated CI coverage
- HTTP requests can be correlated through request IDs
- request logs are machine-readable
- application metrics are available to operators
- the router reuses the existing application logger

### Negative

- CI workflows require maintenance as the project changes
- integration validation can take longer than isolated tests
- metrics are currently process-local
- production deployments may eventually require external metrics and log aggregation

## Alternatives Considered

A single CI workflow was rejected because separate workflows provide clearer failure boundaries.

A separate HTTP logger was rejected because PeakCloud already has an application logger that can be reused.

Logs without metrics were rejected because aggregate HTTP behavior is easier to inspect with both structured logs and request metrics.

## Result

Feature 11 establishes the CI and observability foundation required before production deployment and final project polish.
