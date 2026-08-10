"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { register } from "@/lib/api";

export default function RegisterForm() {
  const router = useRouter();

  const [displayName, setDisplayName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setLoading(true);

    try {
      await register(
        email,
        displayName,
        password,
      );

      router.push("/dashboard");
      router.refresh();
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : "Unable to create account",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <div>
        <label
          htmlFor="display-name"
          className="mb-2 block text-sm font-medium text-zinc-300"
        >
          Display name
        </label>

        <input
          id="display-name"
          required
          minLength={2}
          maxLength={100}
          autoComplete="name"
          value={displayName}
          onChange={(event) => setDisplayName(event.target.value)}
          className="w-full rounded-xl border border-zinc-800 bg-zinc-900 px-4 py-3 outline-none transition focus:border-zinc-600"
        />
      </div>

      <div>
        <label
          htmlFor="email"
          className="mb-2 block text-sm font-medium text-zinc-300"
        >
          Email
        </label>

        <input
          id="email"
          type="email"
          required
          autoComplete="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          className="w-full rounded-xl border border-zinc-800 bg-zinc-900 px-4 py-3 outline-none transition focus:border-zinc-600"
        />
      </div>

      <div>
        <label
          htmlFor="password"
          className="mb-2 block text-sm font-medium text-zinc-300"
        >
          Password
        </label>

        <input
          id="password"
          type="password"
          required
          minLength={12}
          maxLength={128}
          autoComplete="new-password"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          className="w-full rounded-xl border border-zinc-800 bg-zinc-900 px-4 py-3 outline-none transition focus:border-zinc-600"
        />
      </div>

      <p className="text-xs leading-5 text-zinc-500">
        Use at least 12 characters.
      </p>

      {error ? (
        <p className="text-sm text-red-400">{error}</p>
      ) : null}

      <button
        type="submit"
        disabled={loading}
        className="w-full rounded-xl bg-white px-4 py-3 font-medium text-black transition hover:bg-zinc-200 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {loading ? "Creating account..." : "Create account"}
      </button>
    </form>
  );
}
