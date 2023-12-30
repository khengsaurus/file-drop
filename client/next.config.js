/** @type {import('next').NextConfig} */

const nextConfig = {
  basePath: "/file-drop",
  env: {
    DEV: process.env.DEV,
    SERVER_URL: process.env.SERVER_URL,
    SERVER_URL_DEV: process.env.SERVER_URL_DEV,
    SERVICE: process.env.SERVICE,
  },
  async rewrites() {
    const serverUrl =
      process.env.DEV === "1"
        ? process.env.SERVER_URL_DEV
        : process.env.SERVER_URL;
    return [
      {
        source: "/file/:path*",
        destination: `${serverUrl}/file/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
