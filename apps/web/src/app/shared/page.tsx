"use client";

import { useEffect, useState } from "react";

import { ResourceShare, getSharedWithMe } from "@/lib/api";

export default function SharedPage() {
  const [shares, setShares] = useState<ResourceShare[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const result = await getSharedWithMe();

        if (!cancelled) {
          setShares(result);
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

  return (
    <main className="min-h-screen bg-zinc-950 p-8 text-white">
      <div className="mx-auto max-w-6xl">
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
              className="grid grid-cols-[1fr_120px_160px] items-center rounded-xl border border-zinc-800 bg-zinc-900/40 p-4"
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
            </div>
          ))}
        </div>
      </div>
    </main>
  );
}
