"use client";

import {
  useCallback,
  useEffect,
  useState,
} from "react";

import {
  getTrash,
  permanentlyDeleteTrashItem,
  restoreTrashItem,
  type TrashItem,
} from "@/lib/api";

function formatDeletedAt(value: string) {
  return new Intl.DateTimeFormat("en", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}

export default function TrashManager() {
  const [items, setItems] = useState<TrashItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [workingId, setWorkingId] =
    useState<string | null>(null);
  const [error, setError] = useState("");

  const loadTrash = useCallback(async () => {
    try {
      setError("");

      const result = await getTrash();

      setItems(result.items);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Unable to load trash",
      );
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      void loadTrash();
    }, 0);

    return () => {
      window.clearTimeout(timer);
    };
  }, [loadTrash]);

  async function handleRestore(item: TrashItem) {
    try {
      setWorkingId(item.id);
      setError("");

      await restoreTrashItem(
        item.resource_type,
        item.id,
      );

      await loadTrash();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Unable to restore item",
      );
    } finally {
      setWorkingId(null);
    }
  }

  async function handlePermanentDelete(
    item: TrashItem,
  ) {
    const confirmed = window.confirm(
      `Permanently delete "${item.name}"? This cannot be undone.`,
    );

    if (!confirmed) {
      return;
    }

    try {
      setWorkingId(item.id);
      setError("");

      await permanentlyDeleteTrashItem(
        item.resource_type,
        item.id,
      );

      await loadTrash();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Unable to permanently delete item",
      );
    } finally {
      setWorkingId(null);
    }
  }

  return (
    <section className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-white">
          Trash
        </h1>

        <p className="mt-2 text-sm text-zinc-400">
          Restore deleted files and folders or remove
          them permanently.
        </p>
      </div>

      {error && (
        <div className="rounded-lg border border-red-900/50 bg-red-950/30 px-4 py-3 text-sm text-red-300">
          {error}
        </div>
      )}

      <div className="overflow-hidden rounded-xl border border-zinc-800 bg-zinc-950">
        <div className="grid grid-cols-[minmax(0,1fr)_120px_190px_180px] gap-4 border-b border-zinc-800 px-5 py-3 text-xs font-medium uppercase tracking-wide text-zinc-500">
          <span>Name</span>
          <span>Type</span>
          <span>Deleted</span>
          <span>Actions</span>
        </div>

        {loading ? (
          <div className="px-5 py-10 text-center text-sm text-zinc-500">
            Loading trash...
          </div>
        ) : items.length === 0 ? (
          <div className="px-5 py-10 text-center">
            <p className="text-sm font-medium text-zinc-300">
              Trash is empty
            </p>

            <p className="mt-1 text-sm text-zinc-500">
              Deleted files and folders will appear here.
            </p>
          </div>
        ) : (
          <div className="divide-y divide-zinc-800">
            {items.map((item) => {
              const working =
                workingId === item.id;

              return (
                <div
                  key={`${item.resource_type}-${item.id}`}
                  className="grid grid-cols-[minmax(0,1fr)_120px_190px_180px] items-center gap-4 px-5 py-4"
                >
                  <div className="min-w-0">
                    <p className="truncate text-sm font-medium text-zinc-200">
                      {item.name}
                    </p>
                  </div>

                  <p className="text-sm capitalize text-zinc-400">
                    {item.resource_type}
                  </p>

                  <p className="text-sm text-zinc-500">
                    {formatDeletedAt(
                      item.deleted_at,
                    )}
                  </p>

                  <div className="flex items-center gap-4">
                    <button
                      type="button"
                      disabled={working}
                      onClick={() =>
                        void handleRestore(item)
                      }
                      className="text-sm text-zinc-300 hover:text-white disabled:opacity-50"
                    >
                      Restore
                    </button>

                    <button
                      type="button"
                      disabled={working}
                      onClick={() =>
                        void handlePermanentDelete(
                          item,
                        )
                      }
                      className="text-sm text-red-400 hover:text-red-300 disabled:opacity-50"
                    >
                      Delete forever
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </section>
  );
}
