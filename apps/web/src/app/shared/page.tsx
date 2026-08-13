"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import FilePreviewModal from "@/components/preview/FilePreviewModal";

import {
  DriveContents,
  ResourceShare,
  downloadFile,
  getDrive,
  getSharedWithMe,
} from "@/lib/api";

export default function SharedPage() {
  const [shares, setShares] = useState<ResourceShare[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [previewFile, setPreviewFile] = useState<{
    id: string;
    name: string;
    allowDownload: boolean;
  } | null>(null);

  const [openFolder, setOpenFolder] = useState<{
    share: ResourceShare;
    drive: DriveContents;
  } | null>(null);

  const [folderLoading, setFolderLoading] = useState(false);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const result = await getSharedWithMe();

        if (!cancelled) {
          setShares(result);
          setError("");
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof Error
              ? err.message
              : "Unable to load shared resources",
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
  }, []);

  async function handleOpenFolder(share: ResourceShare) {
    try {
      setFolderLoading(true);
      setError("");

      const drive = await getDrive(share.resource_id);

      setOpenFolder({
        share,
        drive,
      });
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Unable to open shared folder",
      );
    } finally {
      setFolderLoading(false);
    }
  }

  async function handleDownload(share: ResourceShare) {
    if (!share.allow_download) {
      return;
    }

    try {
      setError("");
      await downloadFile(share.resource_id, share.resource_name);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Unable to download shared file",
      );
    }
  }

  return (
    <main className="min-h-screen bg-zinc-950 p-8 text-white">
      <div className="mx-auto max-w-6xl">
        <nav className="mb-6 flex items-center gap-4 text-sm">
          <Link href="/dashboard" className="text-zinc-400 hover:text-white">
            My Drive
          </Link>

          <span className="text-zinc-700">/</span>

          <span className="text-zinc-200">Shared with me</span>
        </nav>

        <h1 className="text-3xl font-semibold">Shared with me</h1>

        <p className="mt-2 text-zinc-500">
          Files and folders other PeakCloud users have shared with you.
        </p>

        {loading && <p className="mt-8 text-zinc-400">Loading...</p>}

        {error && (
          <div className="mt-8 rounded-lg border border-red-900 bg-red-950/30 p-4 text-red-300">
            {error}
          </div>
        )}

        {!loading && !error && shares.length === 0 && (
          <div className="mt-8 rounded-xl border border-zinc-800 p-8 text-center text-zinc-500">
            Nothing has been shared with you yet.
          </div>
        )}

        <div className="mt-8 space-y-3">
          {shares.map((share) => (
            <div
              key={share.id}
              className="grid grid-cols-[1fr_110px_150px_220px] items-center gap-4 rounded-xl border border-zinc-800 bg-zinc-900/40 p-4"
            >
              <div>
                <p className="font-medium">
                  {share.resource_type === "folder" ? "📁" : "📄"}{" "}
                  {share.resource_name}
                </p>

                <p className="mt-1 text-xs text-zinc-500">
                  {share.resource_type}
                </p>
              </div>

              <span className="text-sm text-zinc-400">{share.permission}</span>

              <span className="text-sm text-zinc-400">
                {share.allow_download
                  ? "Download allowed"
                  : "Download disabled"}
              </span>

              <div className="flex justify-end gap-4 text-sm">
                {share.resource_type === "file" && (
                  <button
                    type="button"
                    onClick={() =>
                      setPreviewFile({
                        id: share.resource_id,
                        name: share.resource_name,
                        allowDownload: share.allow_download,
                      })
                    }
                    className="text-blue-400 hover:text-blue-300"
                  >
                    Open
                  </button>
                )}

                {share.resource_type === "folder" && (
                  <button
                    type="button"
                    disabled={folderLoading}
                    onClick={() => void handleOpenFolder(share)}
                    className="text-blue-400 hover:text-blue-300 disabled:opacity-50"
                  >
                    Open
                  </button>
                )}

                {share.resource_type === "file" && share.allow_download && (
                  <button
                    type="button"
                    onClick={() => void handleDownload(share)}
                    className="text-zinc-300 hover:text-white"
                  >
                    Download
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>

        {openFolder && (
          <div className="mt-10 rounded-xl border border-zinc-800 bg-zinc-900/40">
            <div className="flex items-center justify-between border-b border-zinc-800 p-5">
              <div>
                <p className="text-xs uppercase tracking-wide text-zinc-500">
                  Shared folder
                </p>

                <h2 className="mt-1 text-xl font-semibold">
                  📁 {openFolder.share.resource_name}
                </h2>
              </div>

              <button
                type="button"
                onClick={() => setOpenFolder(null)}
                className="text-sm text-zinc-400 hover:text-white"
              >
                Close
              </button>
            </div>

            <div className="divide-y divide-zinc-800">
              {openFolder.drive.folders.map((folder) => (
                <div
                  key={folder.id}
                  className="flex items-center justify-between p-4"
                >
                  <span>📁 {folder.name}</span>

                  <button
                    type="button"
                    onClick={() =>
                      void handleOpenFolder({
                        ...openFolder.share,
                        resource_id: folder.id,
                        resource_name: folder.name,
                        resource_type: "folder",
                      })
                    }
                    className="text-sm text-blue-400 hover:text-blue-300"
                  >
                    Open
                  </button>
                </div>
              ))}

              {openFolder.drive.files.map((file) => (
                <div
                  key={file.id}
                  className="flex items-center justify-between p-4"
                >
                  <span>📄 {file.name}</span>

                  <div className="flex gap-4 text-sm">
                    <button
                      type="button"
                      onClick={() =>
                        setPreviewFile({
                          id: file.id,
                          name: file.name,
                          allowDownload: openFolder.share.allow_download,
                        })
                      }
                      className="text-blue-400 hover:text-blue-300"
                    >
                      Open
                    </button>

                    {openFolder.share.allow_download && (
                      <button
                        type="button"
                        onClick={() =>
                          void downloadFile(file.id, file.name).catch(
                            (err: unknown) => {
                              setError(
                                err instanceof Error
                                  ? err.message
                                  : "Unable to download shared file",
                              );
                            },
                          )
                        }
                        className="text-zinc-300 hover:text-white"
                      >
                        Download
                      </button>
                    )}
                  </div>
                </div>
              ))}

              {openFolder.drive.folders.length === 0 &&
                openFolder.drive.files.length === 0 && (
                  <div className="p-8 text-center text-zinc-500">
                    This folder is empty.
                  </div>
                )}
            </div>
          </div>
        )}
      </div>

      {previewFile && (
        <FilePreviewModal
          fileId={previewFile.id}
          allowDownload={previewFile.allowDownload}
          onClose={() => setPreviewFile(null)}
          onDownload={(fileId) => {
            if (!previewFile.allowDownload) {
              return;
            }

            void downloadFile(fileId, previewFile.name).catch(
              (err: unknown) => {
                setError(
                  err instanceof Error
                    ? err.message
                    : "Unable to download shared file",
                );
              },
            );
          }}
        />
      )}
    </main>
  );
}
