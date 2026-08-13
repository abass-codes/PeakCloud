"use client";

import { useCallback, useEffect, useState } from "react";

import {
  DriveContents,
  getDrive,
} from "@/lib/api";

type MoveResourceType = "file" | "folder";

type Props = {
  resourceId: string;
  resourceName: string;
  resourceType: MoveResourceType;
  onClose: () => void;
  onMove: (destinationFolderId?: string) => Promise<void>;
};

export default function MoveModal({
  resourceId,
  resourceName,
  resourceType,
  onClose,
  onMove,
}: Props) {
  const [drive, setDrive] = useState<DriveContents>({
    breadcrumbs: [],
    folders: [],
    files: [],
  });
  const [currentFolderId, setCurrentFolderId] = useState<string | undefined>();
  const [loading, setLoading] = useState(true);
  const [moving, setMoving] = useState(false);
  const [error, setError] = useState("");

  const load = useCallback(async (folderId?: string) => {
    try {
      setLoading(true);
      setError("");

      const result = await getDrive(folderId);
      setDrive(result);
      setCurrentFolderId(folderId);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Unable to load folders",
      );
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    let cancelled = false;

    async function loadInitialFolders() {
      try {
        const result = await getDrive();

        if (!cancelled) {
          setDrive(result);
          setCurrentFolderId(undefined);
          setError("");
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof Error
              ? err.message
              : "Unable to load folders",
          );
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void loadInitialFolders();

    return () => {
      cancelled = true;
    };
  }, []);

  async function handleMove() {
    try {
      setMoving(true);
      setError("");

      await onMove(currentFolderId);
      onClose();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Unable to move item",
      );
    } finally {
      setMoving(false);
    }
  }

  const folders = drive.folders.filter(
    (folder) =>
      !(resourceType === "folder" && folder.id === resourceId),
  );

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
      onMouseDown={(event) => {
        if (event.target === event.currentTarget) {
          onClose();
        }
      }}
    >
      <div className="w-full max-w-xl rounded-2xl border border-zinc-800 bg-zinc-950 p-6 shadow-2xl">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h2 className="text-xl font-semibold text-white">
              Move {resourceName}
            </h2>
            <p className="mt-1 text-sm text-zinc-500">
              Choose a destination folder.
            </p>
          </div>

          <button
            type="button"
            onClick={onClose}
            className="text-sm text-zinc-400 hover:text-white"
          >
            Close
          </button>
        </div>

        {error && (
          <div className="mt-5 rounded-lg border border-red-900 bg-red-950/30 p-3 text-sm text-red-300">
            {error}
          </div>
        )}

        <div className="mt-6 flex flex-wrap items-center gap-2 border-b border-zinc-800 pb-4 text-sm text-zinc-400">
          <button
            type="button"
            onClick={() => void load()}
            className="hover:text-white"
          >
            My Drive
          </button>

          {drive.breadcrumbs.map((folder) => (
            <div key={folder.id} className="flex items-center gap-2">
              <span>/</span>

              <button
                type="button"
                onClick={() => void load(folder.id)}
                className="hover:text-white"
              >
                {folder.name}
              </button>
            </div>
          ))}
        </div>

        <div className="mt-4 max-h-72 overflow-y-auto rounded-xl border border-zinc-800">
          {loading ? (
            <div className="p-6 text-sm text-zinc-500">
              Loading folders...
            </div>
          ) : folders.length === 0 ? (
            <div className="p-6 text-sm text-zinc-500">
              No folders inside this location.
            </div>
          ) : (
            folders.map((folder) => (
              <button
                key={folder.id}
                type="button"
                onClick={() => void load(folder.id)}
                className="flex w-full items-center gap-3 border-b border-zinc-900 px-4 py-3 text-left text-sm last:border-b-0 hover:bg-zinc-900"
              >
                <span>📁</span>
                <span className="font-medium text-zinc-200">
                  {folder.name}
                </span>
              </button>
            ))
          )}
        </div>

        <div className="mt-6 flex items-center justify-between gap-4">
          <p className="text-xs text-zinc-500">
            Move to the folder currently shown above.
          </p>

          <div className="flex gap-3">
            <button
              type="button"
              onClick={onClose}
              className="rounded-lg border border-zinc-700 px-4 py-2 text-sm text-zinc-300 hover:text-white"
            >
              Cancel
            </button>

            <button
              type="button"
              disabled={moving}
              onClick={() => void handleMove()}
              className="rounded-lg bg-white px-4 py-2 text-sm font-medium text-black disabled:opacity-50"
            >
              {moving ? "Moving..." : "Move here"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
