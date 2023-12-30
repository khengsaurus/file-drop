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
      "Content-Type": "application/json"
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

export * from "./consts";
