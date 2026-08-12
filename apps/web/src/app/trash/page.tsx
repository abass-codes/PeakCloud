import Link from "next/link";

import TrashManager from "@/components/trash/TrashManager";

export default function TrashPage() {
  return (
    <main className="min-h-screen bg-black px-6 py-8">
      <div className="mx-auto max-w-6xl space-y-6">
        <nav className="flex items-center gap-4 text-sm">
          <Link
            href="/dashboard"
            className="text-zinc-400 hover:text-white"
          >
            My Drive
          </Link>

          <span className="text-zinc-700">/</span>

          <span className="text-zinc-200">
            Trash
          </span>
        </nav>

        <TrashManager />
      </div>
    </main>
  );
}
