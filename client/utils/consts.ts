export const isDev = process.env.DEV === "1";

export const serverUrl = isDev
  ? process.env.SERVER_URL_DEV
  : process.env.SERVER_URL;
