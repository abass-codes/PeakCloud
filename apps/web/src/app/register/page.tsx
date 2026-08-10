import Link from "next/link";

import RegisterForm from "@/components/auth/RegisterForm";

export default function RegisterPage() {
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
          Create your account
        </h1>

        <p className="mt-2 text-zinc-400">
          Start your PeakCloud workspace.
        </p>

        <div className="mt-8">
          <RegisterForm />
        </div>

        <p className="mt-6 text-sm text-zinc-400">
          Already have an account?{" "}
          <Link href="/login" className="text-white hover:underline">
            Sign in
          </Link>
        </p>
      </section>
    </main>
  );
}
