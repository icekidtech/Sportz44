import Constants from 'expo-constants';

const ENV_BASE_URL =
  (Constants.expoConfig?.extra as { apiUrl?: string } | undefined)?.apiUrl ??
  process.env.EXPO_PUBLIC_API_URL;

export const API_BASE_URL =
  ENV_BASE_URL ??
  (__DEV__ ? 'http://localhost:8080' : 'https://api.sportz44.com');

type RequestOptions = RequestInit;

let onUnauthorized: (() => void) | null = null;
export function setUnauthorizedHandler(fn: () => void) {
  onUnauthorized = fn;
}

let isRefreshing = false;
let refreshPromise: Promise<void> | null = null;

async function doRefresh(): Promise<void> {
  if (isRefreshing && refreshPromise) return refreshPromise;
  isRefreshing = true;
  refreshPromise = (async () => {
    const res = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    });
    if (!res.ok) {
      onUnauthorized?.();
      throw new Error('Session expired');
    }
  })();
  try {
    await refreshPromise;
  } finally {
    isRefreshing = false;
    refreshPromise = null;
  }
}

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(opts.headers as Record<string, string>),
  };

  const doFetch = () =>
    fetch(`${API_BASE_URL}${path}`, {
      ...opts,
      headers,
      credentials: 'include',
    });

  let res = await doFetch();

  // Auto-refresh on 401 (access cookie expired) — retry once.
  // Skip for auth endpoints themselves to avoid loops.
  const isAuthPath = path.startsWith('/api/auth/');
  if (res.status === 401 && !isAuthPath) {
    try {
      await doRefresh();
      res = await doFetch();
    } catch {
      // doRefresh already called onUnauthorized
      throw new Error('Session expired — please sign in again');
    }
  }

  if (!res.ok) {
    const body = await res.text();
    let msg = body;
    try {
      const j = JSON.parse(body);
      msg = j.error ?? j.message ?? body;
    } catch {}
    throw new Error(msg || `HTTP ${res.status}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export const api = {
  get: <T>(path: string, opts?: RequestOptions) => request<T>(path, { ...opts, method: 'GET' }),
  post: <T>(path: string, body?: unknown, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'POST', body: body ? JSON.stringify(body) : undefined }),
  put: <T>(path: string, body?: unknown, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'PUT', body: body ? JSON.stringify(body) : undefined }),
  del: <T>(path: string, opts?: RequestOptions) => request<T>(path, { ...opts, method: 'DELETE' }),
};
