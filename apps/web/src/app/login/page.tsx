import Link from "next/link";

import LoginForm from "@/components/auth/LoginForm";

export default function LoginPage() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-zinc-950 px-6 text-white">
      <section className="w-full max-w-md">
        <Link
          href="/"
          className="mb-10 inline-block text-sm font-semibold tracking-[0.25em] text-zinc-400"
        >
          PEAKCLOUD
        </Link>

        <h1 className="text-3xl font-semibold tracking-tight">
          Welcome back
        </h1>

        <p className="mt-2 text-zinc-400">
          Sign in to access your PeakCloud workspace.
        </p>

        <div className="mt-8">
          <LoginForm />
        </div>

        <p className="mt-6 text-sm text-zinc-400">
          No account?{" "}
          <Link href="/register" className="text-white hover:underline">
            Create one
          </Link>
        </p>
      </section>
    </main>
  );
}
