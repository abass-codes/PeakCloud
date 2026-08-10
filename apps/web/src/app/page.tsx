import Link from "next/link";

export default function Home() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-zinc-950 px-6 text-white">
      <section className="max-w-3xl text-center">
        <p className="mb-4 text-sm font-medium uppercase tracking-[0.3em] text-zinc-400">
          PeakCloud
        </p>

        <h1 className="text-5xl font-semibold tracking-tight sm:text-6xl">
          Your files. Securely stored.
        </h1>

        <p className="mx-auto mt-6 max-w-2xl text-lg leading-8 text-zinc-400">
          A production-grade cloud storage and file synchronization platform.
        </p>

        <div className="mt-10 flex flex-wrap justify-center gap-4">
          <Link
            href="/register"
            className="rounded-xl bg-white px-6 py-3 font-medium text-black transition hover:bg-zinc-200"
          >
            Create account
          </Link>

          <Link
            href="/login"
            className="rounded-xl border border-zinc-800 px-6 py-3 font-medium text-zinc-200 transition hover:bg-zinc-900"
          >
            Sign in
          </Link>
        </div>
      </section>
    </main>
  );
}
