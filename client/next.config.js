/** @type {import('next').NextConfig} */
const nextConfig = {
  env: {
    DEV: process.env.DEV,
    FILE_URL: process.env.FILE_URL,
    FILE_URL_DEV: process.env.FILE_URL_DEV,
    SERVER_URL: process.env.SERVER_URL,
    SERVER_URL_DEV: process.env.SERVER_URL_DEV,
  },
};

module.exports = nextConfig;
