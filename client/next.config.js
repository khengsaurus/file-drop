/** @type {import('next').NextConfig} */

const isDev = process.env.DEV === "1";

const nextConfig = {
  basePath: isDev ? "" : "/file-drop",
  env: {
    API_BASE_PATH: isDev ? "/api" : "/file-drop/api",
    SERVICE: process.env.SERVICE
  },
  async redirects() {
    return [
      { source: "/file", destination: "/", permanent: true },
      { source: "/url", destination: "/", permanent: true },
      { source: "/not-found", destination: `/404`, permanent: true }
    ];
  },
  async rewrites() {
    const serverUrl = isDev
      ? process.env.SERVER_URL_DEV
      : process.env.SERVER_URL;
    return [
      {
        source: "/api/:path*",
        destination: `${serverUrl}/api/:path*`
      },
      {
        source: "/file/:path*",
        destination: `${serverUrl}/stream/:path*`
      },
      {
        source: "/url/:path*",
        destination: `${serverUrl}/url/:path*`
      }
    ];
  }
};

module.exports = nextConfig;
