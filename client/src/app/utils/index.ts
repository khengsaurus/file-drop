export function get(
  url: string,
  params?: Record<string, any>,
  headers?: Record<string, any>
) {
  const urlWithParams = params ? `${url}?${new URLSearchParams(params)}` : url;
  return fetch(urlWithParams, {
    method: "GET",
    headers: headers || {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  });
}

export function post(url: string, data = {}) {
  return fetch(url, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
}

export function uploadFile(
  file: File,
  url: string,
  abortSignal?: AbortController
) {
  return fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": file.type,
      "Content-Length": String(file.size),
    },
    body: file,
    signal: abortSignal?.signal,
  });
}

export function getFileUrl(windowLocation: typeof window.location) {
  const { host, pathname } = windowLocation || {};
  if (host.endsWith(".vercel.app")) {
    return host + pathname + "/file";
  }
  return host + "/file";
}

export function getUrlUrl(windowLocation: typeof window.location) {
  const { host, pathname } = windowLocation || {};
  if (host.endsWith(".vercel.app")) {
    return host + pathname + "/url";
  }
  return host + "/url";
}

/** @see https://www.freecodecamp.org/news/check-if-a-javascript-string-is-a-url/ */
export function isValidUrl(urlString: string) {
  var urlPattern = new RegExp(
    "^(https?:\\/\\/)?" + // validate protocol
      "((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|" + // validate domain name
      "((\\d{1,3}\\.){3}\\d{1,3}))" + // validate OR ip (v4) address
      "(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*" + // validate port and path
      "(\\?[;&a-z\\d%_.~+=-]*)?" + // validate query string
      "(\\#[-a-z\\d_]*)?$",
    "i"
  ); // validate fragment locator
  return !!urlPattern.test(urlString);
}

export * from "./consts";
