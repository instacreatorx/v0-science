import {
  clearAccessToken,
  emitSessionExpired,
  getAccessToken,
} from '@/lib/auth/session';

export class ApiError extends Error {
  readonly status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError;
}

interface RequestOptions {
  method?: string;
  body?: object;
  headers?: Record<string, string>;
  /** When true, a 401 clears the stored session. Default: true for authenticated calls. */
  auth?: boolean;
}

function resolveApiBaseUrl(): string {
  // Same-origin proxy avoids CORS and keeps auth headers reliable in production.
  if (typeof window !== 'undefined') {
    return '/api/backend';
  }
  return process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
}

export async function apiRequest<T>(
  endpoint: string,
  options: RequestOptions = {}
): Promise<T> {
  const { method = 'GET', body, headers = {}, auth = true } = options;
  const token = auth ? getAccessToken() : null;

  const requestHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
    Accept: 'application/json',
    ...headers,
  };

  if (token) {
    requestHeaders.Authorization = `Bearer ${token}`;
  }

  let response: Response;

  try {
    response = await fetch(`${resolveApiBaseUrl()}${endpoint}`, {
      method,
      headers: requestHeaders,
      body: body ? JSON.stringify(body) : undefined,
      credentials: 'same-origin',
      cache: 'no-store',
    });
  } catch {
    throw new ApiError('Network error. Please check your connection.', 0);
  }

  if (!response.ok) {
    const payload = await response.json().catch(() => ({ error: 'Request failed' }));
    const message =
      typeof payload.error === 'string' ? payload.error : 'Request failed';

    if (auth && response.status === 401 && token) {
      clearAccessToken();
      emitSessionExpired();
    }

    throw new ApiError(message, response.status);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json() as Promise<T>;
}
