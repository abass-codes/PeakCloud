# ADR 0005: Authenticated File Preview

## Status

Accepted

## Context

PeakCloud users need to inspect common files without downloading them first.

Files are private by default and binary contents are stored in S3-compatible object storage.

Directly exposing object-storage paths would bypass PeakCloud's authorization model and leak internal storage details.

## Decision

PeakCloud provides authenticated preview endpoints through the Go API.

Preview support is derived from existing file metadata rather than stored as additional database state.

The backend classifies files using MIME type and, where appropriate, filename extension.

Supported preview categories are:

- images
- PDFs
- text
- source code
- video
- audio

Unsupported content uses the existing download workflow.

## Authorization

Preview requests resolve files through owner-scoped database queries.

Internal object keys remain private.

A user cannot preview another user's file even when the file UUID is known.

## Content Delivery

Previewable binary content is streamed through an authenticated API endpoint using inline content disposition.

Private preview responses are not intended for shared caching.

## Text Files

Text and source-code previews have a size limit.

This prevents large objects from being unnecessarily loaded into application or browser memory for text rendering.

## Frontend

The web application renders previews in a modal.

Viewer selection is based on backend preview classification.

Images, PDFs, audio, and video use native browser rendering capabilities.

Text and source-code files use a text viewer.

## Consequences

### Positive

- Private files remain behind PeakCloud authorization.
- Internal object-storage keys are not exposed.
- Preview behavior is centralized.
- Unsupported formats retain a predictable download fallback.
- No schema migration is required.
- Existing object storage remains unchanged.

### Negative

- Preview traffic passes through the API.
- Browser codec support varies for audio and video.
- Very large text files cannot be previewed inline.
- More advanced document formats may require dedicated processing later.
