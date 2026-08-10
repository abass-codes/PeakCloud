# ADR 0001: Use a Monorepo

## Status

Accepted

## Context

PeakCloud contains a web application, Go API, background workers, infrastructure configuration, documentation, and shared development tooling.

Maintaining these components in separate repositories during initial development would introduce unnecessary coordination and versioning overhead.

## Decision

PeakCloud will use a monorepo containing the frontend, backend, infrastructure configuration, documentation, tests, and development tooling.

## Consequences

### Positive

- Atomic frontend and backend changes
- Simplified local development
- Centralized CI/CD
- Unified documentation
- Easier issue and release management

### Negative

- Repository size will increase as the platform grows
- CI workflows must avoid unnecessary work as additional services are introduced

The architecture may be reconsidered if independent service lifecycles eventually justify repository separation.
