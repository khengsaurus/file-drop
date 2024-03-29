/** @type {import('next').NextConfig} */

const nextConfig = {
  basePath: process.env.DEV === "1" ? "" : "/file-drop",
  env: {
    DEV: process.env.DEV,
    SERVER_URL: process.env.SERVER_URL,
    SERVER_URL_DEV: process.env.SERVER_URL_DEV,
    SERVICE: process.env.SERVICE,
  },
  async redirects() {
    return [
      {
        source: "/file",
        destination: "/",
        permanent: true,
      },
      {
        source: "/url",
        destination: "/",
        permanent: true,
      },
      {
        source: "/not-found",
        destination: `/404`,
        permanent: true,
      },
    ];
  },
  async rewrites() {
    const serverUrl =
      process.env.DEV === "1"
        ? process.env.SERVER_URL_DEV
        : process.env.SERVER_URL;
    return [
      {
        source: "/file/:path*",
        destination: `${serverUrl}/stream/:path*`,
      },
      {
        source: "/url/:path*",
        destination: `${serverUrl}/url/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
