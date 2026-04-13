import { env } from "@/shared/config/env";

import { clearStoredToken, getStoredToken } from "./auth-token";
import { ApiError } from "./errors";

type RequestOptions = {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  body?: unknown;
  headers?: HeadersInit;
  token?: string | null;
  multipart?: boolean;
};

async function parseResponse<T>(response: Response): Promise<T> {
  const contentType = response.headers.get("content-type") ?? "";
  const isJson = contentType.includes("application/json");
  const payload = isJson ? await response.json() : await response.text();

  if (!response.ok) {
    const message =
      typeof payload === "object" && payload && "error" in payload
        ? String(payload.error)
        : typeof payload === "string" && payload
          ? payload
          : "Request failed";

    if (response.status === 401) {
      clearStoredToken();
    }

    throw new ApiError(message, response.status);
  }

  return payload as T;
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}) {
  const token = options.token ?? getStoredToken();
  const headers = new Headers(options.headers);

  if (!options.multipart) {
    headers.set("Content-Type", "application/json");
  }
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(`${env.apiBaseUrl}${path}`, {
    method: options.method ?? "GET",
    headers,
    body:
      options.body == null
        ? undefined
        : options.multipart
          ? (options.body as BodyInit)
          : JSON.stringify(options.body),
    cache: "no-store"
  });

  return parseResponse<T>(response);
}

export const apiClient = {
  get: <T>(path: string, options?: Omit<RequestOptions, "method" | "body">) =>
    apiRequest<T>(path, { ...options, method: "GET" }),
  post: <T>(path: string, body?: unknown, options?: Omit<RequestOptions, "method" | "body">) =>
    apiRequest<T>(path, { ...options, method: "POST", body }),
  put: <T>(path: string, body?: unknown, options?: Omit<RequestOptions, "method" | "body">) =>
    apiRequest<T>(path, { ...options, method: "PUT", body }),
  patch: <T>(path: string, body?: unknown, options?: Omit<RequestOptions, "method" | "body">) =>
    apiRequest<T>(path, { ...options, method: "PATCH", body }),
  upload: <T>(path: string, formData: FormData, options?: Omit<RequestOptions, "method" | "body" | "multipart">) =>
    apiRequest<T>(path, { ...options, method: "POST", body: formData, multipart: true })
};
