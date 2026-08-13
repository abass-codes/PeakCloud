"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import DriveManager from "@/components/drive/DriveManager";
import { getCurrentUser, logout, User } from "@/lib/api";

function DriveIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      className="h-5 w-5"
    >
      <path d="M4 5.5A1.5 1.5 0 0 1 5.5 4h4l2 2h7A1.5 1.5 0 0 1 20 7.5v10A1.5 1.5 0 0 1 18.5 19h-13A1.5 1.5 0 0 1 4 17.5z" />
    </svg>
  );
}

function SharedIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      className="h-5 w-5"
    >
      <circle cx="9" cy="8" r="3" />
      <path d="M3.5 19c.6-3.2 2.5-5 5.5-5s4.9 1.8 5.5 5" />
      <circle cx="17" cy="9" r="2.5" />
      <path d="M15.5 14.5c2.9-.2 4.7 1.3 5 4.5" />
    </svg>
  );
}

function TrashIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      className="h-5 w-5"
    >
      <path d="M4 7h16" />
      <path d="M9 3h6l1 4H8z" />
      <path d="M7 7l1 13h8l1-13" />
      <path d="M10 11v5M14 11v5" />
    </svg>
  );
}

export default function DashboardPage() {
  const router = useRouter();

  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [signingOut, setSigningOut] = useState(false);

  useEffect(() => {
    let cancelled = false;

    async function loadUser() {
      try {
        const currentUser = await getCurrentUser();

        if (!cancelled) {
          setUser(currentUser);
        }
      } catch {
        router.replace("/login");
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void loadUser();

    return () => {
      cancelled = true;
    };
  }, [router]);

  async function handleLogout() {
    try {
      setSigningOut(true);
      await logout();
      router.replace("/login");
      router.refresh();
    } finally {
      setSigningOut(false);
    }
  }

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[#09090b] text-zinc-100">
        <div className="flex items-center gap-3 text-sm text-zinc-500">
          <span className="h-2 w-2 animate-pulse rounded-full bg-blue-400" />
          Loading PeakCloud...
        </div>
      </main>
    );
  }

  if (!user) {
    return null;
  }

  const initial = user.display_name.trim().charAt(0).toUpperCase() || "P";

  return (
    <main className="min-h-screen bg-[#09090b] text-zinc-100">
      <div className="flex min-h-screen">
        <aside className="fixed inset-y-0 left-0 hidden w-[260px] flex-col border-r border-white/[0.06] bg-[#0b0b0e] lg:flex">
          <div className="px-6 pb-8 pt-7">
            <Link href="/" className="inline-block">
              <p className="text-2xl font-semibold tracking-[0.18em] text-blue-400">
                PEAKCLOUD
              </p>
            </Link>
          </div>

          <div className="px-3">
            <p className="mb-2 px-3 text-[10px] font-semibold uppercase tracking-[0.22em] text-zinc-500">
              Workspace
            </p>

            <nav className="space-y-1">
              <Link
                href="/dashboard"
                className="flex items-center gap-3 rounded-xl bg-white/[0.07] px-3 py-2.5 text-sm font-medium text-white"
              >
                <DriveIcon />
                My Drive
              </Link>

              <Link
                href="/shared"
                className="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-zinc-500 transition hover:bg-white/[0.05] hover:text-zinc-200"
              >
                <SharedIcon />
                Shared with me
              </Link>

              <Link
                href="/trash"
                className="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-zinc-500 transition hover:bg-white/[0.05] hover:text-zinc-200"
              >
                <TrashIcon />
                Trash
              </Link>
            </nav>
          </div>

          <div className="mt-8 px-6">
            <div className="border-t border-white/[0.06] pt-6">
              <p className="text-[10px] font-semibold uppercase tracking-[0.22em] text-zinc-500">
                PeakCloud
              </p>

              <p className="mt-3 text-sm font-medium text-zinc-500">
                Secure workspace
              </p>

              <p className="mt-1 text-xs leading-5 text-zinc-500">
                Private storage, sharing, version history and recovery.
              </p>
            </div>
          </div>

          <div className="mt-auto border-t border-white/[0.06] p-4">
            <div className="flex items-center gap-3 rounded-xl px-2 py-2">
              <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-zinc-100 text-sm font-semibold text-zinc-950">
                {initial}
              </div>

              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-medium text-zinc-200">
                  {user.display_name}
                </p>

                <p className="truncate text-xs text-zinc-500">{user.email}</p>
              </div>
            </div>

            <button
              type="button"
              disabled={signingOut}
              onClick={() => void handleLogout()}
              className="mt-2 w-full rounded-xl px-3 py-2 text-left text-sm text-zinc-500 transition hover:bg-white/[0.05] hover:text-zinc-200 disabled:opacity-50"
            >
              {signingOut ? "Signing out..." : "Sign out"}
            </button>
          </div>
        </aside>

        <div className="min-w-0 flex-1 lg:ml-[260px]">
          <header className="sticky top-0 z-20 border-b border-white/[0.06] bg-[#09090b]/90 backdrop-blur-xl">
            <div className="mx-auto flex max-w-[1440px] items-center justify-between px-6 py-4 lg:px-10">
              <div>
                <p className="text-xs font-medium text-zinc-500">
                  PeakCloud workspace
                </p>

                <h1 className="mt-0.5 text-lg font-semibold tracking-tight text-zinc-100">
                  My Drive
                </h1>
              </div>

              <div className="hidden items-center gap-2 rounded-full border border-white/[0.07] bg-white/[0.03] px-3 py-1.5 text-xs text-zinc-500 sm:flex">
                <span className="h-1.5 w-1.5 rounded-full bg-emerald-400" />
                Secure session
              </div>
            </div>
          </header>

          <div className="mx-auto max-w-[1440px] px-6 py-8 lg:px-10 lg:py-10">
            <section className="mb-10">
              <p className="text-sm text-zinc-500">
                Welcome back, {user.display_name}.
              </p>

              <div className="mt-2 flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
                <div>
                  <h2 className="text-3xl font-semibold tracking-tight text-white sm:text-4xl">
                    Your files
                  </h2>

                  <p className="mt-3 max-w-xl text-sm leading-6 text-zinc-500">
                    Store, organize, share, and recover your files from one
                    secure workspace.
                  </p>
                </div>

                <div className="flex items-center gap-2 text-xs text-zinc-500">
                  <span className="rounded-full border border-white/[0.06] px-3 py-1.5">
                    Encrypted storage
                  </span>

                  <span className="rounded-full border border-white/[0.06] px-3 py-1.5">
                    Version history
                  </span>
                </div>
              </div>
            </section>

            <section className="mt-8">
              <DriveManager />
            </section>
          </div>
        </div>
      </div>
    </main>
  );
}
