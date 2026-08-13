import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "PeakCloud",
  description:
    "A production-grade cloud storage and file synchronization platform.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const currentYear = new Date().getFullYear();

  return (
    <html lang="en">
      <body>
        {children}

        <footer className="border-t border-zinc-900 bg-zinc-950 px-6 py-6 text-center text-xs text-zinc-500">
          © {currentYear} Yakubu Mohammed Abass. All rights reserved.
        </footer>
      </body>
    </html>
  );
}
