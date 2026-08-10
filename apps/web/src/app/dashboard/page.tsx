"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { getMe, logout, type User } from "@/lib/api";

export default function DashboardPage() {
  const router = useRouter();

  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let active = true;

    async function loadUser() {
      try {
        const response = await getMe();

        if (active) {
          setUser(response.user);
        }
      } catch {
        router.replace("/login");
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    }

    void loadUser();

    return () => {
      active = false;
    };
  }, [router]);

  async function handleLogout() {
    try {
      await logout();
    } finally {
      router.replace("/login");
      router.refresh();
    }
  }

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-zinc-950 text-zinc-400">
        Loading PeakCloud...
      </main>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <main className="min-h-screen bg-zinc-950 text-white">
      <header className="border-b border-zinc-900">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-5">
          <span className="font-semibold tracking-[0.2em]">
            PEAKCLOUD
          </span>

          <button
            type="button"
            onClick={handleLogout}
            className="rounded-lg border border-zinc-800 px-4 py-2 text-sm text-zinc-300 transition hover:bg-zinc-900"
          >
            Sign out
          </button>
        </div>
      </header>

      <section className="mx-auto max-w-6xl px-6 py-16">
        <p className="text-sm text-zinc-500">Workspace</p>

        <h1 className="mt-2 text-4xl font-semibold tracking-tight">
          Welcome, {user.display_name}
        </h1>

        <p className="mt-3 text-zinc-400">{user.email}</p>

        <div className="mt-12 rounded-2xl border border-zinc-900 bg-zinc-900/30 p-8">
          <h2 className="text-xl font-medium">Your files</h2>

          <p className="mt-3 text-zinc-500">
            File storage arrives in the next PeakCloud feature.
          </p>
        </div>
      </section>
    </main>
  );
}
