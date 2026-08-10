"use client";

import { ChangeEvent, DragEvent, useCallback, useEffect, useState } from "react";

import {
  deleteFile,
  downloadFile,
  listFiles,
  StoredFile,
  uploadFile,
} from "@/lib/api";

function formatBytes(bytes: number): string {
  if (bytes === 0) {
    return "0 B";
  }

  const units = ["B", "KB", "MB", "GB", "TB"];
  const index = Math.min(
    Math.floor(Math.log(bytes) / Math.log(1024)),
    units.length - 1,
  );

  const value = bytes / 1024 ** index;

  return `${value.toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
}

function formatDate(value: string): string {
  return new Intl.DateTimeFormat("en", {
    month: "short",
    day: "numeric",
    year: "numeric",
  }).format(new Date(value));
}

export default function FileManager() {
  const [files, setFiles] = useState<StoredFile[]>([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [dragging, setDragging] = useState(false);
  const [error, setError] = useState("");

  const refresh = useCallback(async () => {
    try {
      setError("");
      setFiles(await listFiles());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to load files");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    let cancelled = false;

    async function loadFiles() {
      try {
        const result = await listFiles();

        if (!cancelled) {
          setFiles(result);
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof Error ? err.message : "Unable to load files",
          );
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void loadFiles();

    return () => {
      cancelled = true;
    };
  }, []);

  async function handleUpload(file: File) {
    try {
      setUploading(true);
      setError("");

      await uploadFile(file);
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setUploading(false);
    }
  }

  async function handleInput(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];

    if (file) {
      await handleUpload(file);
      event.target.value = "";
    }
  }

  async function handleDrop(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    setDragging(false);

    const file = event.dataTransfer.files?.[0];

    if (file) {
      await handleUpload(file);
    }
  }

  async function handleDelete(file: StoredFile) {
    const confirmed = window.confirm(
      `Delete "${file.name}"? This action cannot be undone.`,
    );

    if (!confirmed) {
      return;
    }

    try {
      setError("");
      await deleteFile(file.id);
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Delete failed");
    }
  }

  async function handleDownload(file: StoredFile) {
    try {
      setError("");
      await downloadFile(file);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Download failed");
    }
  }

  return (
    <section className="mt-10">
      <div
        onDragEnter={(event) => {
          event.preventDefault();
          setDragging(true);
        }}
        onDragOver={(event) => {
          event.preventDefault();
          setDragging(true);
        }}
        onDragLeave={() => setDragging(false)}
        onDrop={(event) => void handleDrop(event)}
        className={`rounded-2xl border border-dashed p-8 text-center transition ${
          dragging
            ? "border-blue-400 bg-blue-500/10"
            : "border-zinc-700 bg-zinc-900/40"
        }`}
      >
        <p className="text-lg font-medium text-zinc-100">
          Drop a file here
        </p>

        <p className="mt-2 text-sm text-zinc-400">
          or choose a file from your computer
        </p>

        <label className="mt-5 inline-flex cursor-pointer rounded-lg bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-500">
          {uploading ? "Uploading..." : "Choose file"}

          <input
            type="file"
            className="hidden"
            disabled={uploading}
            onChange={(event) => void handleInput(event)}
          />
        </label>

        <p className="mt-3 text-xs text-zinc-500">
          Maximum file size: 100 MB
        </p>
      </div>

      {error && (
        <div className="mt-4 rounded-lg border border-red-900/50 bg-red-950/30 px-4 py-3 text-sm text-red-300">
          {error}
        </div>
      )}

      <div className="mt-8 overflow-hidden rounded-xl border border-zinc-800">
        <div className="grid grid-cols-[minmax(0,1fr)_110px_150px_180px] gap-4 border-b border-zinc-800 bg-zinc-900/80 px-5 py-3 text-xs font-medium uppercase tracking-wide text-zinc-500">
          <span>Name</span>
          <span>Size</span>
          <span>Uploaded</span>
          <span className="text-right">Actions</span>
        </div>

        {loading ? (
          <div className="px-5 py-12 text-center text-sm text-zinc-500">
            Loading files...
          </div>
        ) : files.length === 0 ? (
          <div className="px-5 py-12 text-center">
            <p className="text-sm font-medium text-zinc-300">
              No files yet
            </p>
            <p className="mt-1 text-sm text-zinc-500">
              Upload your first file to PeakCloud.
            </p>
          </div>
        ) : (
          files.map((file) => (
            <div
              key={file.id}
              className="grid grid-cols-[minmax(0,1fr)_110px_150px_180px] items-center gap-4 border-b border-zinc-800 px-5 py-4 last:border-b-0"
            >
              <div className="min-w-0">
                <p className="truncate text-sm font-medium text-zinc-200">
                  {file.name}
                </p>
                <p className="mt-1 truncate text-xs text-zinc-500">
                  {file.content_type}
                </p>
              </div>

              <span className="text-sm text-zinc-400">
                {formatBytes(file.size_bytes)}
              </span>

              <span className="text-sm text-zinc-400">
                {formatDate(file.created_at)}
              </span>

              <div className="flex justify-end gap-2">
                <button
                  type="button"
                  onClick={() => void handleDownload(file)}
                  className="rounded-md border border-zinc-700 px-3 py-1.5 text-xs text-zinc-300 transition hover:bg-zinc-800"
                >
                  Download
                </button>

                <button
                  type="button"
                  onClick={() => void handleDelete(file)}
                  className="rounded-md border border-red-900/70 px-3 py-1.5 text-xs text-red-300 transition hover:bg-red-950/50"
                >
                  Delete
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </section>
  );
}
