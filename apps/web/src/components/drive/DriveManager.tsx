"use client";

import FilePreviewModal from "@/components/preview/FilePreviewModal";

import {
  ChangeEvent,
  DragEvent,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";

import {
  BulkItem,
  DriveContents,
  StoredFile,
  bulkDelete,
  bulkDownload,
  copyFile,
  createFolder,
  deleteFile,
  deleteFolder,
  downloadFile,
  getDrive,
  renameFile,
  renameFolder,
  uploadFile,
} from "@/lib/api";

type SortMode = "name" | "modified" | "size";

function formatBytes(bytes: number) {
  if (bytes < 1024) {
    return `${bytes} B`;
  }

  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }

  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export default function DriveManager() {
  const [drive, setDrive] = useState<DriveContents>({
    breadcrumbs: [],
    folders: [],
    files: [],
  });

  const [currentFolderId, setCurrentFolderId] = useState<string | undefined>();

  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [sortMode, setSortMode] = useState<SortMode>("name");
  const [filter, setFilter] = useState("");
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [dragging, setDragging] = useState(false);
  const [error, setError] = useState("");
  const [previewFileId, setPreviewFileId] = useState<string | null>(null);

  const refresh = useCallback(async (folderId?: string) => {
    try {
      setError("");
      const result = await getDrive(folderId);
      setDrive(result);
      setSelected(new Set());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to load drive");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const result = await getDrive();

        if (!cancelled) {
          setDrive(result);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Unable to load drive");
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

  const folders = useMemo(() => {
    const query = filter.trim().toLowerCase();

    return [...drive.folders]
      .filter((folder) => folder.name.toLowerCase().includes(query))
      .sort((a, b) => {
        if (sortMode === "modified") {
          return (
            new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()
          );
        }

        return a.name.localeCompare(b.name);
      });
  }, [drive.folders, filter, sortMode]);

  const files = useMemo(() => {
    const query = filter.trim().toLowerCase();

    return [...drive.files]
      .filter((file) => file.name.toLowerCase().includes(query))
      .sort((a, b) => {
        if (sortMode === "modified") {
          return (
            new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()
          );
        }

        if (sortMode === "size") {
          return b.size_bytes - a.size_bytes;
        }

        return a.name.localeCompare(b.name);
      });
  }, [drive.files, filter, sortMode]);

  function selectKey(type: "file" | "folder", id: string) {
    return `${type}:${id}`;
  }

  function toggle(type: "file" | "folder", id: string) {
    const key = selectKey(type, id);

    setSelected((current) => {
      const next = new Set(current);

      if (next.has(key)) {
        next.delete(key);
      } else {
        next.add(key);
      }

      return next;
    });
  }

  function selectedItems(): BulkItem[] {
    return Array.from(selected).map((value) => {
      const [type, id] = value.split(":");

      return {
        type: type as "file" | "folder",
        id,
      };
    });
  }

  async function openFolder(id?: string) {
    setLoading(true);
    setCurrentFolderId(id);
    await refresh(id);
  }

  async function handleCreateFolder() {
    const name = window.prompt("Folder name");

    if (!name) {
      return;
    }

    try {
      await createFolder(name, currentFolderId);
      await refresh(currentFolderId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to create folder");
    }
  }

  async function handleUpload(file: File) {
    try {
      setUploading(true);
      setError("");

      await uploadFile(file, currentFolderId);
      await refresh(currentFolderId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to upload file");
    } finally {
      setUploading(false);
    }
  }

  async function handleInput(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];

    if (file) {
      await handleUpload(file);
    }

    event.target.value = "";
  }

  async function handleDrop(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    setDragging(false);

    const file = event.dataTransfer.files?.[0];

    if (file) {
      await handleUpload(file);
    }
  }

  async function handleRenameFolder(id: string, current: string) {
    const name = window.prompt("Rename folder", current);

    if (!name || name === current) {
      return;
    }

    try {
      await renameFolder(id, name);
      await refresh(currentFolderId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to rename folder");
    }
  }

  async function handleRenameFile(file: StoredFile) {
    const name = window.prompt("Rename file", file.name);

    if (!name || name === file.name) {
      return;
    }

    try {
      await renameFile(file.id, name);
      await refresh(currentFolderId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to rename file");
    }
  }

  async function handleDeleteFolder(id: string) {
    if (!window.confirm("Delete this folder and everything inside it?")) {
      return;
    }

    try {
      await deleteFolder(id);
      await refresh(currentFolderId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to delete folder");
    }
  }

  async function handleDeleteFile(id: string) {
    if (!window.confirm("Delete this file permanently?")) {
      return;
    }

    try {
      await deleteFile(id);
      await refresh(currentFolderId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to delete file");
    }
  }

  async function handleCopy(file: StoredFile) {
    try {
      await copyFile(file.id, currentFolderId);
      await refresh(currentFolderId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to copy file");
    }
  }

  async function handleBulkDelete() {
    const items = selectedItems();

    if (
      items.length === 0 ||
      !window.confirm(`Delete ${items.length} selected item(s)?`)
    ) {
      return;
    }

    try {
      await bulkDelete(items);
      await refresh(currentFolderId);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Bulk delete failed");
    }
  }

  async function handleBulkDownload() {
    const fileIds = selectedItems()
      .filter((item) => item.type === "file")
      .map((item) => item.id);

    if (fileIds.length === 0) {
      setError("Select at least one file to download.");
      return;
    }

    try {
      await bulkDownload(fileIds);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Bulk download failed");
    }
  }

  return (
    <section className="mt-8">
      <div className="flex flex-wrap items-center gap-3">
        <button
          type="button"
          onClick={() => void handleCreateFolder()}
          className="rounded-lg bg-white px-4 py-2 text-sm font-medium text-black"
        >
          New folder
        </button>

        <label className="cursor-pointer rounded-lg border border-zinc-700 px-4 py-2 text-sm">
          {uploading ? "Uploading..." : "Upload file"}
          <input
            type="file"
            disabled={uploading}
            onChange={(event) => void handleInput(event)}
            className="hidden"
          />
        </label>

        <input
          value={filter}
          onChange={(event) => setFilter(event.target.value)}
          placeholder="Filter this folder..."
          className="min-w-56 rounded-lg border border-zinc-800 bg-zinc-950 px-3 py-2 text-sm outline-none"
        />

        <select
          value={sortMode}
          onChange={(event) => setSortMode(event.target.value as SortMode)}
          className="rounded-lg border border-zinc-800 bg-zinc-950 px-3 py-2 text-sm"
        >
          <option value="name">Name</option>
          <option value="modified">Last modified</option>
          <option value="size">Size</option>
        </select>
      </div>

      <div className="mt-6 flex flex-wrap items-center gap-2 text-sm text-zinc-400">
        <button
          type="button"
          onClick={() => void openFolder()}
          className="hover:text-white"
        >
          My Drive
        </button>

        {drive.breadcrumbs.map((folder) => (
          <div key={folder.id} className="flex items-center gap-2">
            <span>/</span>
            <button
              type="button"
              onClick={() => void openFolder(folder.id)}
              className="hover:text-white"
            >
              {folder.name}
            </button>
          </div>
        ))}
      </div>

      {selected.size > 0 && (
        <div className="mt-5 flex items-center gap-3 rounded-xl border border-zinc-800 bg-zinc-900/60 p-3">
          <span className="text-sm">{selected.size} selected</span>

          <button
            type="button"
            onClick={() => void handleBulkDownload()}
            className="text-sm text-zinc-300 hover:text-white"
          >
            Download files
          </button>

          <button
            type="button"
            onClick={() => void handleBulkDelete()}
            className="text-sm text-red-400 hover:text-red-300"
          >
            Delete
          </button>
        </div>
      )}

      {error && (
        <div className="mt-5 rounded-lg border border-red-900 bg-red-950/30 p-3 text-sm text-red-300">
          {error}
        </div>
      )}

      <div
        onDragEnter={(event) => {
          event.preventDefault();
          setDragging(true);
        }}
        onDragOver={(event) => event.preventDefault()}
        onDragLeave={() => setDragging(false)}
        onDrop={(event) => void handleDrop(event)}
        className={`mt-6 overflow-hidden rounded-xl border ${
          dragging ? "border-zinc-400 bg-zinc-900" : "border-zinc-800"
        }`}
      >
        <div className="grid grid-cols-[40px_1fr_140px_140px_260px] border-b border-zinc-800 bg-zinc-900/40 px-4 py-3 text-xs uppercase tracking-wide text-zinc-500">
          <span />
          <span>Name</span>
          <span>Type</span>
          <span>Size</span>
          <span>Actions</span>
        </div>

        {loading ? (
          <div className="p-8 text-sm text-zinc-500">Loading drive...</div>
        ) : folders.length === 0 && files.length === 0 ? (
          <div className="p-12 text-center text-sm text-zinc-500">
            Drop a file here or create a folder.
          </div>
        ) : (
          <>
            {folders.map((folder) => {
              const key = selectKey("folder", folder.id);

              return (
                <div
                  key={folder.id}
                  className="grid grid-cols-[40px_1fr_140px_140px_260px] items-center border-b border-zinc-900 px-4 py-3 text-sm"
                >
                  <input
                    type="checkbox"
                    checked={selected.has(key)}
                    onChange={() => toggle("folder", folder.id)}
                  />

                  <button
                    type="button"
                    onClick={() => void openFolder(folder.id)}
                    className="text-left font-medium hover:underline"
                  >
                    📁 {folder.name}
                  </button>

                  <span className="text-zinc-500">Folder</span>
                  <span className="text-zinc-600">—</span>

                  <div className="flex gap-4">
                    <button
                      type="button"
                      onClick={() =>
                        void handleRenameFolder(folder.id, folder.name)
                      }
                      className="text-zinc-400 hover:text-white"
                    >
                      Rename
                    </button>

                    <button
                      type="button"
                      onClick={() => void handleDeleteFolder(folder.id)}
                      className="text-red-400 hover:text-red-300"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              );
            })}

            {files.map((file) => {
              const key = selectKey("file", file.id);

              return (
                <div
                  key={file.id}
                  className="grid grid-cols-[40px_1fr_140px_140px_260px] items-center border-b border-zinc-900 px-4 py-3 text-sm last:border-b-0"
                >
                  <input
                    type="checkbox"
                    checked={selected.has(key)}
                    onChange={() => toggle("file", file.id)}
                  />

                  <button
                    type="button"
                    onClick={() => setPreviewFileId(file.id)}
                    className="truncate pr-4 text-left hover:underline"
                  >
                    📄 {file.name}
                  </button>

                  <span className="truncate text-zinc-500">
                    {file.content_type}
                  </span>

                  <span className="text-zinc-500">
                    {formatBytes(file.size_bytes)}
                  </span>

                  <div className="flex gap-4">
                    <button
                      type="button"
                      onClick={() => void downloadFile(file.id, file.name)}
                      className="text-zinc-400 hover:text-white"
                    >
                      Download
                    </button>

                    <button
                      type="button"
                      onClick={() => void handleRenameFile(file)}
                      className="text-zinc-400 hover:text-white"
                    >
                      Rename
                    </button>

                    <button
                      type="button"
                      onClick={() => void handleCopy(file)}
                      className="text-zinc-400 hover:text-white"
                    >
                      Copy
                    </button>

                    <button
                      type="button"
                      onClick={() => void handleDeleteFile(file.id)}
                      className="text-red-400 hover:text-red-300"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              );
            })}
          </>
        )}
      </div>

      <FilePreviewModal
        fileId={previewFileId}
        onClose={() => setPreviewFileId(null)}
        onDownload={(fileId) => {
          const file = drive.files.find((item) => item.id === fileId);

          if (file) {
            void downloadFile(file.id, file.name);
          }
        }}
      />
    </section>
  );
}
