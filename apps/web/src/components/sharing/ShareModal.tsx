"use client";

import { useCallback, useEffect, useState } from "react";

import {
  PublicShareLink,
  ResourceShare,
  SharePermission,
  ShareResourceType,
  createPublicShareLink,
  createResourceShare,
  deleteResourceShare,
  getPublicShareLinks,
  getResourceShares,
  revokePublicShareLink,
  updateResourceShare,
} from "@/lib/api";

type Props = {
  resourceId: string | null;
  resourceName: string;
  resourceType: ShareResourceType;
  onClose: () => void;
};

export default function ShareModal({
  resourceId,
  resourceName,
  resourceType,
  onClose,
}: Props) {
  const [email, setEmail] = useState("");
  const [permission, setPermission] = useState<SharePermission>("viewer");
  const [allowDownload, setAllowDownload] = useState(true);
  const [password, setPassword] = useState("");
  const [expiresAt, setExpiresAt] = useState("");
  const [shares, setShares] = useState<ResourceShare[]>([]);
  const [links, setLinks] = useState<PublicShareLink[]>([]);
  const [publicURL, setPublicURL] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const refresh = useCallback(async () => {
    if (!resourceId) {
      return;
    }

    const [allShares, allLinks] = await Promise.all([
      getResourceShares(),
      getPublicShareLinks(),
    ]);

    setShares(
      allShares.filter(
        (share) =>
          share.resource_id === resourceId &&
          share.resource_type === resourceType,
      ),
    );

    setLinks(
      allLinks.filter(
        (link) =>
          link.resource_id === resourceId &&
          link.resource_type === resourceType &&
          !link.revoked_at,
      ),
    );
  }, [resourceId, resourceType]);

  useEffect(() => {
    if (!resourceId) {
      return;
    }

    let cancelled = false;

    Promise.all([getResourceShares(), getPublicShareLinks()])
      .then(([allShares, allLinks]) => {
        if (cancelled) {
          return;
        }

        setShares(
          allShares.filter(
            (share) =>
              share.resource_id === resourceId &&
              share.resource_type === resourceType,
          ),
        );

        setLinks(
          allLinks.filter(
            (link) =>
              link.resource_id === resourceId &&
              link.resource_type === resourceType &&
              !link.revoked_at,
          ),
        );
      })
      .catch((err) => {
        if (!cancelled) {
          setError(
            err instanceof Error
              ? err.message
              : "Unable to load sharing settings",
          );
        }
      });

    return () => {
      cancelled = true;
    };
  }, [resourceId, resourceType]);

  if (!resourceId) {
    return null;
  }

  async function shareWithUser() {
    try {
      setBusy(true);
      setError("");

      await createResourceShare({
        recipient_email: email,
        resource_type: resourceType,
        resource_id: resourceId!,
        permission,
        allow_download: allowDownload,
      });

      setEmail("");
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to share resource");
    } finally {
      setBusy(false);
    }
  }

  async function createLink() {
    try {
      setBusy(true);
      setError("");

      const link = await createPublicShareLink({
        resource_type: resourceType,
        resource_id: resourceId!,
        permission,
        allow_download: allowDownload,
        password: password || undefined,
        expires_at: expiresAt ? new Date(expiresAt).toISOString() : undefined,
      });

      setPublicURL(`${window.location.origin}/share/${link.token}`);
      setPassword("");
      setExpiresAt("");

      await refresh();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Unable to create public link",
      );
    } finally {
      setBusy(false);
    }
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
      onMouseDown={(event) => {
        if (event.target === event.currentTarget) {
          onClose();
        }
      }}
    >
      <div className="max-h-[90vh] w-full max-w-2xl overflow-y-auto rounded-2xl border border-zinc-800 bg-zinc-950 p-6">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h2 className="text-xl font-semibold text-white">
              Share {resourceName}
            </h2>
            <p className="mt-1 text-sm text-zinc-500">
              Manage private access and public links.
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

        <div className="mt-6 border-t border-zinc-800 pt-6">
          <h3 className="font-medium">Share with a PeakCloud user</h3>

          <input
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            placeholder="user@example.com"
            className="mt-4 w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2"
          />

          <div className="mt-3 flex flex-wrap items-center gap-3">
            <select
              value={permission}
              onChange={(event) =>
                setPermission(event.target.value as SharePermission)
              }
              className="rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2"
            >
              <option value="viewer">Viewer</option>
              <option value="editor">Editor</option>
            </select>

            <label className="flex items-center gap-2 text-sm text-zinc-300">
              <input
                type="checkbox"
                checked={allowDownload}
                onChange={(event) => setAllowDownload(event.target.checked)}
              />
              Allow download
            </label>

            <button
              type="button"
              disabled={busy || !email.trim()}
              onClick={() => void shareWithUser()}
              className="rounded-lg bg-white px-4 py-2 text-sm font-medium text-black disabled:opacity-50"
            >
              Share
            </button>
          </div>
        </div>

        <div className="mt-6 border-t border-zinc-800 pt-6">
          <h3 className="font-medium">Public link</h3>

          <div className="mt-4 grid gap-3">
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="Optional password"
              className="rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2"
            />

            <input
              type="datetime-local"
              value={expiresAt}
              onChange={(event) => setExpiresAt(event.target.value)}
              className="rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2"
            />

            <button
              type="button"
              disabled={busy}
              onClick={() => void createLink()}
              className="w-fit rounded-lg border border-zinc-700 px-4 py-2 text-sm"
            >
              Create public link
            </button>
          </div>

          {publicURL && (
            <div className="mt-4 rounded-lg border border-zinc-800 bg-zinc-900 p-3">
              <p className="break-all text-sm text-zinc-300">{publicURL}</p>

              <button
                type="button"
                onClick={() => void navigator.clipboard.writeText(publicURL)}
                className="mt-2 text-sm underline"
              >
                Copy link
              </button>
            </div>
          )}
        </div>

        <div className="mt-6 border-t border-zinc-800 pt-6">
          <h3 className="font-medium">People with access</h3>

          <div className="mt-3 space-y-3">
            {shares.length === 0 && (
              <p className="text-sm text-zinc-500">
                Not shared with another user.
              </p>
            )}

            {shares.map((share) => (
              <div
                key={share.id}
                className="rounded-lg border border-zinc-800 p-3"
              >
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <p className="text-sm">{share.recipient_email}</p>
                    <p className="text-xs text-zinc-500">
                      {share.permission} ·{" "}
                      {share.allow_download
                        ? "download allowed"
                        : "download disabled"}
                    </p>
                  </div>

                  <div className="flex gap-3">
                    <button
                      type="button"
                      onClick={() =>
                        void updateResourceShare(share.id, {
                          permission:
                            share.permission === "viewer" ? "editor" : "viewer",
                          allow_download: share.allow_download,
                        }).then(refresh)
                      }
                      className="text-sm text-zinc-300"
                    >
                      Make {share.permission === "viewer" ? "editor" : "viewer"}
                    </button>

                    <button
                      type="button"
                      onClick={() =>
                        void updateResourceShare(share.id, {
                          permission: share.permission,
                          allow_download: !share.allow_download,
                        }).then(refresh)
                      }
                      className="text-sm text-zinc-300"
                    >
                      {share.allow_download
                        ? "Disable download"
                        : "Allow download"}
                    </button>

                    <button
                      type="button"
                      onClick={() =>
                        void deleteResourceShare(share.id).then(refresh)
                      }
                      className="text-sm text-red-400"
                    >
                      Remove
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="mt-6 border-t border-zinc-800 pt-6">
          <h3 className="font-medium">Active public links</h3>

          <div className="mt-3 space-y-3">
            {links.length === 0 && (
              <p className="text-sm text-zinc-500">No active public links.</p>
            )}

            {links.map((link) => (
              <div
                key={link.id}
                className="flex items-center justify-between gap-4 rounded-lg border border-zinc-800 p-3"
              >
                <div>
                  <p className="text-sm">{link.permission}</p>
                  <p className="text-xs text-zinc-500">
                    {link.password_set ? "Password protected" : "No password"}
                    {" · "}
                    {link.allow_download
                      ? "download allowed"
                      : "download disabled"}
                  </p>

                  {link.expires_at && (
                    <p className="mt-1 text-xs text-zinc-500">
                      Expires {new Date(link.expires_at).toLocaleString()}
                    </p>
                  )}
                </div>

                <button
                  type="button"
                  onClick={() =>
                    void revokePublicShareLink(link.id).then(refresh)
                  }
                  className="text-sm text-red-400"
                >
                  Revoke
                </button>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
