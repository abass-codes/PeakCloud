import type { NextConfig } from "next";

const API_URL =
  process.env.API_INTERNAL_URL ??
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8080";

const nextConfig: NextConfig = {
  output: "standalone",

  experimental: {
    proxyClientMaxBodySize: "100mb",
  },

  async rewrites() {
    return [
      {
        source: "/api/backend/:path*",
        destination: `${API_URL}/:path*`,
      },
    ];
  },
};

export default nextConfig;
