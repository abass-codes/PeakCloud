"use client";

import {
  ChangeEvent,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";

import {
  FileVersion,
  downloadFileVersion,
  getFileVersions,
  restoreFileVersion,
  uploadFileVersion,
} from "@/lib/api";

type Props = {
  fileId: string;
  fileName: string;
  onClose: () => void;
  onChanged: () => void | Promise<void>;
};

function formatBytes(bytes: number) {
  if (bytes < 1024) {
    return `${bytes} B`;
  }

  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }

  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}

export default function VersionHistoryModal({
  fileId,
  fileName,
  onClose,
  onChanged,
}: Props) {
  const [versions, setVersions] = useState<FileVersion[]>([]);
  const [loading, setLoading] = useState(true);
  const [working, setWorking] = useState(false);
  const [error, setError] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  const loadVersions = useCallback(async () => {
    try {
      setError("");

      const result = await getFileVersions(fileId);

      setVersions(result);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Unable to load version history",
      );
    } finally {
      setLoading(false);
    }
  }, [fileId]);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      void loadVersions();
    }, 0);

    return () => {
      window.clearTimeout(timer);
    };
  }, [loadVersions]);

  async function handleUpload(
    event: ChangeEvent<HTMLInputElement>,
  ) {
    const file = event.target.files?.[0];

    if (!file) {
      return;
    }

    try {
      setWorking(true);
      setError("");

      await uploadFileVersion(fileId, file);
      await loadVersions();
      await onChanged();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Unable to upload new version",
      );
    } finally {
      setWorking(false);
      event.target.value = "";
    }
  }

  async function handleRestore(version: FileVersion) {
    const confirmed = window.confirm(
      `Restore version ${version.version_number}? A new version will be created from this historical version.`,
    );

    if (!confirmed) {
      return;
    }

    try {
      setWorking(true);
      setError("");

      await restoreFileVersion(
        fileId,
        version.version_number,
      );

      await loadVersions();
      await onChanged();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Unable to restore version",
      );
    } finally {
      setWorking(false);
    }
  }

  const latestVersion =
    versions.length > 0
      ? Math.max(
          ...versions.map(
            (version) => version.version_number,
          ),
        )
      : 0;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/75 p-4"
      role="dialog"
      aria-modal="true"
      aria-labelledby="version-history-title"
    >
      <div className="w-full max-w-3xl overflow-hidden rounded-2xl border border-zinc-800 bg-zinc-950 shadow-2xl">
        <div className="flex items-start justify-between border-b border-zinc-800 p-5">
          <div>
            <h2
              id="version-history-title"
              className="text-lg font-semibold text-white"
            >
              Version history
            </h2>

            <p className="mt-1 max-w-xl truncate text-sm text-zinc-400">
              {fileName}
            </p>
          </div>

          <button
            type="button"
            onClick={onClose}
            className="rounded-lg px-3 py-1 text-zinc-400 hover:bg-zinc-900 hover:text-white"
          >
            Close
          </button>
        </div>

        <div className="flex items-center justify-between border-b border-zinc-800 p-5">
          <div>
            <p className="text-sm font-medium text-white">
              Historical versions
            </p>

            <p className="mt-1 text-xs text-zinc-500">
              Restoring an older version creates a new
              version and preserves history.
            </p>
          </div>

          <div>
            <input
              ref={inputRef}
              type="file"
              className="hidden"
              disabled={working}
              onChange={(event) =>
                void handleUpload(event)
              }
            />

            <button
              type="button"
              disabled={working}
              onClick={() => inputRef.current?.click()}
              className="rounded-lg bg-white px-4 py-2 text-sm font-medium text-black disabled:cursor-not-allowed disabled:opacity-50"
            >
              {working
                ? "Working..."
                : "Upload new version"}
            </button>
          </div>
        </div>

        {error && (
          <div className="mx-5 mt-5 rounded-lg border border-red-900 bg-red-950/30 p-3 text-sm text-red-300">
            {error}
          </div>
        )}

        <div className="max-h-[60vh] overflow-y-auto p-5">
          {loading ? (
            <div className="py-10 text-center text-sm text-zinc-500">
              Loading version history...
            </div>
          ) : versions.length === 0 ? (
            <div className="py-10 text-center text-sm text-zinc-500">
              No versions found.
            </div>
          ) : (
            <div className="space-y-3">
              {versions.map((version) => {
                const current =
                  version.version_number ===
                  latestVersion;

                return (
                  <div
                    key={version.id}
                    className="flex flex-wrap items-center justify-between gap-4 rounded-xl border border-zinc-800 bg-zinc-900/40 p-4"
                  >
                    <div>
                      <div className="flex items-center gap-2">
                        <span className="font-medium text-white">
                          Version{" "}
                          {version.version_number}
                        </span>

                        {current && (
                          <span className="rounded-full border border-zinc-700 px-2 py-0.5 text-xs text-zinc-300">
                            Current
                          </span>
                        )}
                      </div>

                      <div className="mt-2 space-y-1 text-xs text-zinc-500">
                        <p>
                          {formatDate(
                            version.created_at,
                          )}
                        </p>

                        <p>
                          {formatBytes(
                            version.size_bytes,
                          )}{" "}
                          · {version.content_type}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-3">
                      <button
                        type="button"
                        disabled={working}
                        onClick={() =>
                          void downloadFileVersion(
                            fileId,
                            version.version_number,
                            `${fileName}.v${version.version_number}`,
                          )
                        }
                        className="text-sm text-zinc-300 hover:text-white disabled:opacity-50"
                      >
                        Download
                      </button>

                      {!current && (
                        <button
                          type="button"
                          disabled={working}
                          onClick={() =>
                            void handleRestore(
                              version,
                            )
                          }
                          className="text-sm text-zinc-300 hover:text-white disabled:opacity-50"
                        >
                          Restore
                        </button>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
