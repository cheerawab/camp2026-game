import { createAppError } from "./error"

type RequestOptions = Omit<RequestInit, "body"> & {
  searchParams?: Record<string, string | number | boolean | null | undefined>
}

type JsonRequestOptions = RequestOptions & {
  json?: unknown
}

function resolveRequestUrl(
  path: string,
  searchParams?: RequestOptions["searchParams"],
) {
  const url =
    path.startsWith("http://") || path.startsWith("https://")
      ? new URL(path)
      : typeof window === "undefined"
        ? new URL(
            path,
            process.env.APP_ORIGIN ??
              process.env.VITE_APP_ORIGIN ??
              "http://localhost:3000",
          )
        : new URL(path, window.location.origin)

  Object.entries(searchParams ?? {}).forEach(([key, value]) => {
    if (value != null) url.searchParams.set(key, String(value))
  })

  return url
}

async function parseJsonSafely(response: Response) {
  const text = await response.text()
  if (!text) return null

  try {
    return JSON.parse(text) as unknown
  } catch {
    return text
  }
}

export const apiClient = {
  async get(path: string, options: RequestOptions = {}) {
    return request(path, { ...options, method: "GET" })
  },

  async post(path: string, options: JsonRequestOptions = {}) {
    return request(path, { ...options, method: "POST" })
  },

  async put(path: string, options: JsonRequestOptions = {}) {
    return request(path, { ...options, method: "PUT" })
  },

  async delete(path: string, options: JsonRequestOptions = {}) {
    return request(path, { ...options, method: "DELETE" })
  },
}

async function request(path: string, options: JsonRequestOptions = {}) {
  const { json, searchParams, headers, ...init } = options
  const response = await fetch(resolveRequestUrl(path, searchParams), {
    ...init,
    credentials: init.credentials ?? "include",
    headers: {
      Accept: "application/json",
      ...(json === undefined ? {} : { "Content-Type": "application/json" }),
      ...headers,
    },
    body: json === undefined ? undefined : JSON.stringify(json),
  })
  const body = await parseJsonSafely(response)

  if (!response.ok) {
    throw createAppError({
      status: response.status,
      body,
      fallbackMessage: response.statusText,
    })
  }

  return body
}
