"use client";

import { FormEvent, useEffect, useState } from "react";
import { useParams } from "next/navigation";

import {
  PublicShareLink,
  publicShareContentUrl,
  publicShareDownloadUrl,
  resolvePublicShare,
} from "@/lib/api";

export default function PublicSharePage() {
  const params = useParams<{ token: string }>();
  const token = params.token;

  const [link, setLink] = useState<PublicShareLink | null>(null);
  const [password, setPassword] = useState("");
  const [acceptedPassword, setAcceptedPassword] = useState("");
  const [passwordRequired, setPasswordRequired] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  async function load(value = "") {
    try {
      setLoading(true);
      setError("");

      const result = await resolvePublicShare(token, value);

      setLink(result);
      setAcceptedPassword(value);
      setPasswordRequired(false);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Unable to open shared resource";

      if (message.toLowerCase().includes("password")) {
        setPasswordRequired(true);
      } else {
        setError(message);
      }
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    let cancelled = false;

    resolvePublicShare(token)
      .then((result) => {
        if (cancelled) {
          return;
        }

        setLink(result);
        setAcceptedPassword("");
        setPasswordRequired(false);
        setError("");
      })
      .catch((err) => {
        if (cancelled) {
          return;
        }

        const message =
          err instanceof Error ? err.message : "Unable to open shared resource";

        if (message.toLowerCase().includes("password")) {
          setPasswordRequired(true);
        } else {
          setError(message);
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [token]);

  function submit(event: FormEvent) {
    event.preventDefault();
    void load(password);
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-zinc-950 p-6 text-white">
      <div className="w-full max-w-3xl rounded-2xl border border-zinc-800 bg-zinc-900/40 p-8">
        <p className="text-sm font-medium text-zinc-500">PeakCloud</p>
        <h1 className="mt-2 text-2xl font-semibold">Shared resource</h1>

        {loading && (
          <p className="mt-6 text-zinc-400">Opening secure share...</p>
        )}

        {!loading && passwordRequired && !link && (
          <form onSubmit={submit} className="mt-6">
            <label className="text-sm text-zinc-300">
              This share is password protected.
            </label>

            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="mt-3 w-full rounded-lg border border-zinc-700 bg-zinc-950 px-3 py-2"
              placeholder="Password"
            />

            <button
              type="submit"
              className="mt-3 rounded-lg bg-white px-4 py-2 text-sm font-medium text-black"
            >
              Open
            </button>
          </form>
        )}

        {!loading && error && (
          <div className="mt-6 rounded-lg border border-red-900 bg-red-950/30 p-4 text-red-300">
            {error}
          </div>
        )}

        {!loading && link && (
          <div className="mt-6">
            <div className="rounded-xl border border-zinc-800 p-5">
              <p className="text-lg font-medium">
                {link.resource_type === "folder" ? "📁" : "📄"}{" "}
                {link.resource_name}
              </p>

              <p className="mt-2 text-sm text-zinc-500">
                Permission: {link.permission}
              </p>

              <p className="mt-1 text-sm text-zinc-500">
                {link.allow_download
                  ? "Downloads allowed"
                  : "Downloads disabled"}
              </p>
            </div>

            {link.resource_type === "file" && (
              <div className="mt-5">
                <a
                  href={publicShareContentUrl(token, acceptedPassword)}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex rounded-lg border border-zinc-700 px-4 py-2 text-sm"
                >
                  Open file
                </a>

                {link.allow_download && (
                  <a
                    href={publicShareDownloadUrl(token, acceptedPassword)}
                    className="ml-3 inline-flex rounded-lg bg-white px-4 py-2 text-sm font-medium text-black"
                  >
                    Download
                  </a>
                )}
              </div>
            )}

            {link.resource_type === "folder" && (
              <p className="mt-5 text-sm text-zinc-500">
                Folder access has been shared. Full shared-folder navigation
                will use PeakCloud&apos;s centralized authorization layer.
              </p>
            )}
          </div>
        )}
      </div>
    </main>
  );
}
