import Constants from 'expo-constants';

const ENV_BASE_URL =
  (Constants.expoConfig?.extra as { apiUrl?: string } | undefined)?.apiUrl ??
  process.env.EXPO_PUBLIC_API_URL;

export const API_BASE_URL =
  ENV_BASE_URL ??
  (__DEV__ ? 'http://localhost:8080' : 'https://api.sportz44.com');

type RequestOptions = RequestInit & { skipAuth?: boolean };

let getAuthToken: (() => string | null) | null = null;
export function setAuthTokenGetter(fn: () => string | null) {
  getAuthToken = fn;
}

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(opts.headers as Record<string, string>),
  };
  if (!opts.skipAuth && getAuthToken) {
    const token = getAuthToken();
    if (token) headers.Authorization = `Bearer ${token}`;
  }
  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...opts,
    headers,
    credentials: 'include',
  });
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
