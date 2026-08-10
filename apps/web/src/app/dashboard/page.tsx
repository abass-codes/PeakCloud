"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import FileManager from "@/components/files/FileManager";
import { getCurrentUser, logout, User } from "@/lib/api";

export default function DashboardPage() {
  const router = useRouter();

  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [signingOut, setSigningOut] = useState(false);

  useEffect(() => {
    async function loadUser() {
      try {
        const currentUser = await getCurrentUser();
        setUser(currentUser);
      } catch {
        router.replace("/login");
      } finally {
        setLoading(false);
      }
    }

    void loadUser();
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
      <main className="min-h-screen bg-zinc-950 text-zinc-100">
        <div className="mx-auto max-w-6xl px-6 py-12">
          <p className="text-sm text-zinc-500">Loading workspace...</p>
        </div>
      </main>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <main className="min-h-screen bg-zinc-950 text-zinc-100">
      <div className="mx-auto max-w-6xl px-6 py-8">
        <header className="flex items-center justify-between border-b border-zinc-800 pb-6">
          <div>
            <p className="text-xs font-semibold tracking-[0.25em] text-blue-400">
              PEAKCLOUD
            </p>

            <h1 className="mt-2 text-2xl font-semibold">
              Your files
            </h1>

            <p className="mt-1 text-sm text-zinc-500">
              {user.display_name} · {user.email}
            </p>
          </div>

          <button
            type="button"
            disabled={signingOut}
            onClick={() => void handleLogout()}
            className="rounded-lg border border-zinc-700 px-4 py-2 text-sm text-zinc-300 transition hover:bg-zinc-900 disabled:opacity-50"
          >
            {signingOut ? "Signing out..." : "Sign out"}
          </button>
        </header>

        <FileManager />
      </div>
    </main>
  );
}
