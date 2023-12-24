export function get(url: string, params?: Record<string, any>) {
  const urlWithParams = params ? `${url}?${new URLSearchParams(params)}` : url;
  return fetch(urlWithParams, {
    method: "GET",
    headers: {
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
  isDev = false,
  abortSignal?: AbortController
) {
  return fetch(url, {
    method: "PUT",
    headers: isDev
      ? {
          "Content-Type": "application/json",
          "x-amz-acl": "public-read-write",
        }
      : { "Content-Type": "multipart/form-data" },
    body: file,
    signal: abortSignal?.signal,
  });
}

export * from "./consts";
