"use client";

import { useEffect, useMemo, useState } from "react";
import {
  FilePreviewResponse,
  getFilePreview,
  getFilePreviewBlob,
} from "@/lib/api";

type Props = {
  fileId: string | null;
  onClose: () => void;
  onDownload: (fileId: string) => void;
  allowDownload?: boolean;
};

export default function FilePreviewModal({
  fileId,
  onClose,
  onDownload,
  allowDownload = true,
}: Props) {
  const [preview, setPreview] = useState<FilePreviewResponse | null>(null);
  const [blob, setBlob] = useState<Blob | null>(null);
  const [text, setText] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const objectUrl = useMemo(() => {
    if (!blob) {
      return null;
    }

    return URL.createObjectURL(blob);
  }, [blob]);

  useEffect(() => {
    return () => {
      if (objectUrl) {
        URL.revokeObjectURL(objectUrl);
      }
    };
  }, [objectUrl]);

  useEffect(() => {
    if (!fileId) {
      return;
    }

    const activeFileId = fileId;
    let cancelled = false;

    async function load() {
      setLoading(true);
      setError(null);
      setPreview(null);
      setBlob(null);
      setText(null);

      try {
        const metadata = await getFilePreview(activeFileId);

        if (cancelled) {
          return;
        }

        setPreview(metadata);

        if (!metadata.preview.previewable) {
          return;
        }

        const content = await getFilePreviewBlob(activeFileId);

        if (cancelled) {
          return;
        }

        if (
          metadata.preview.kind === "text" ||
          metadata.preview.kind === "code"
        ) {
          setText(await content.text());
        } else {
          setBlob(content);
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof Error ? err.message : "Unable to preview file",
          );
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void load();

    return () => {
      cancelled = true;
    };
  }, [fileId]);

  if (!fileId) {
    return null;
  }

  const filename = preview?.file.name ?? "File preview";

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
      onMouseDown={(event) => {
        if (event.target === event.currentTarget) {
          onClose();
        }
      }}
    >
      <div className="flex h-[90vh] w-full max-w-6xl flex-col overflow-hidden rounded-2xl border border-zinc-800 bg-zinc-950 shadow-2xl">
        <header className="flex items-center justify-between gap-4 border-b border-zinc-800 px-5 py-4">
          <div className="min-w-0">
            <p className="truncate font-medium text-zinc-100">{filename}</p>

            {preview && (
              <p className="mt-1 text-xs uppercase tracking-wider text-zinc-500">
                {preview.preview.kind}
              </p>
            )}
          </div>

          <div className="flex shrink-0 items-center gap-2">
            {allowDownload && (
              <button
                type="button"
                onClick={() => onDownload(fileId)}
                className="rounded-lg border border-zinc-700 px-3 py-2 text-sm text-zinc-200 transition hover:bg-zinc-900"
              >
                Download
              </button>
            )}

            <button
              type="button"
              onClick={onClose}
              className="rounded-lg border border-zinc-700 px-3 py-2 text-sm text-zinc-200 transition hover:bg-zinc-900"
            >
              Close
            </button>
          </div>
        </header>

        <div className="flex min-h-0 flex-1 items-center justify-center overflow-auto bg-zinc-900/40 p-5">
          {loading && (
            <p className="text-sm text-zinc-400">Loading preview...</p>
          )}

          {!loading && error && (
            <div className="text-center">
              <p className="text-sm text-red-400">{error}</p>
            </div>
          )}

          {!loading && !error && preview && !preview.preview.previewable && (
            <div className="max-w-md text-center">
              <h2 className="text-xl font-semibold text-zinc-100">
                Preview unavailable
              </h2>

              <p className="mt-2 text-sm text-zinc-400">
                PeakCloud cannot preview this file type in the browser. You can
                still download the original file.
              </p>

              {allowDownload && (
                <button
                  type="button"
                  onClick={() => onDownload(fileId)}
                  className="mt-5 rounded-lg bg-zinc-100 px-4 py-2 text-sm font-medium text-zinc-950"
                >
                  Download file
                </button>
              )}
            </div>
          )}

          {!loading &&
            !error &&
            preview?.preview.kind === "image" &&
            objectUrl && (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                src={objectUrl}
                alt={filename}
                className="max-h-full max-w-full object-contain"
              />
            )}

          {!loading &&
            !error &&
            preview?.preview.kind === "pdf" &&
            objectUrl && (
              <iframe
                src={objectUrl}
                title={filename}
                className="h-full min-h-[70vh] w-full rounded-lg bg-white"
              />
            )}

          {!loading &&
            !error &&
            preview?.preview.kind === "video" &&
            objectUrl && (
              <video src={objectUrl} controls className="max-h-full max-w-full">
                Your browser does not support video playback.
              </video>
            )}

          {!loading &&
            !error &&
            preview?.preview.kind === "audio" &&
            objectUrl && (
              <div className="w-full max-w-2xl rounded-xl border border-zinc-800 bg-zinc-950 p-8">
                <p className="mb-6 text-center text-zinc-300">{filename}</p>

                <audio src={objectUrl} controls className="w-full">
                  Your browser does not support audio playback.
                </audio>
              </div>
            )}

          {!loading &&
            !error &&
            (preview?.preview.kind === "text" ||
              preview?.preview.kind === "code") &&
            text !== null && (
              <pre className="min-h-full w-full overflow-auto whitespace-pre-wrap break-words rounded-xl border border-zinc-800 bg-zinc-950 p-5 text-sm leading-6 text-zinc-200">
                <code>{text}</code>
              </pre>
            )}
        </div>
      </div>
    </div>
  );
}
