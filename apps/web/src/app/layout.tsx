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
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
