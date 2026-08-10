# ADR 0003: Store Binary Files in S3-Compatible Object Storage

## Status

Accepted

## Context

PeakCloud needs to persist user-uploaded binary objects while supporting future horizontal scaling and cloud deployment.

Storing large binary objects directly in PostgreSQL would tightly couple application metadata and object data and would make independent storage scaling more difficult.

## Decision

PeakCloud will separate file metadata from binary object storage.

PostgreSQL stores file metadata and ownership information.

Binary file contents are stored in S3-compatible object storage.

The local development environment uses MinIO.

Object keys are generated independently of user-supplied filenames.

## Consequences

### Positive

- Binary storage can scale independently from PostgreSQL.
- PeakCloud can migrate between S3-compatible providers with limited application changes.
- User-supplied filenames do not determine physical object paths.
- Metadata queries remain relational and efficient.
- Local development closely resembles production object-storage architecture.

### Negative

- File operations span two persistence systems.
- Upload and deletion operations require failure-handling strategies.
- Object storage becomes another production dependency.

## Failure Handling

For uploads, PeakCloud writes the object first and then persists metadata.

If metadata persistence fails, PeakCloud attempts to remove the uploaded object.

Future work may introduce asynchronous reconciliation for rare partial failures.

## Security

Every file record is associated with an authenticated owner.

File metadata and object keys are resolved through owner-scoped database queries.

PeakCloud does not expose internal object keys through the public API.
