const TOKEN_STORAGE_KEY = 'auth_token';
const TOKEN_COOKIE = 'auth_token';
const TOKEN_MAX_AGE_SECONDS = 60 * 60 * 24 * 7; // 7 days — matches JWT expiry

let memoryToken: string | null = null;

function readCookieToken(): string | null {
  if (typeof document === 'undefined') return null;
  const match = document.cookie.match(new RegExp(`(?:^|; )${TOKEN_COOKIE}=([^;]*)`));
  if (!match?.[1]) return null;
  try {
    return decodeURIComponent(match[1]);
  } catch {
    return match[1];
  }
}

function writeCookieToken(token: string) {
  if (typeof document === 'undefined') return;
  const secure = window.location.protocol === 'https:' ? ';Secure' : '';
  document.cookie = `${TOKEN_COOKIE}=${encodeURIComponent(token)};path=/;max-age=${TOKEN_MAX_AGE_SECONDS};SameSite=Lax${secure}`;
}

function clearCookieToken() {
  if (typeof document === 'undefined') return;
  document.cookie = `${TOKEN_COOKIE}=;path=/;max-age=0;SameSite=Lax`;
}

/** Persist JWT in memory, localStorage, and cookie (middleware reads cookie). */
export function setAccessToken(token: string): void {
  if (!token?.trim()) return;
  memoryToken = token;
  if (typeof window === 'undefined') return;
  localStorage.setItem(TOKEN_STORAGE_KEY, token);
  writeCookieToken(token);
}

/** Read JWT — memory → localStorage → cookie, keeping stores in sync. */
export function getAccessToken(): string | null {
  if (memoryToken) return memoryToken;
  if (typeof window === 'undefined') return null;

  const fromStorage = localStorage.getItem(TOKEN_STORAGE_KEY);
  if (fromStorage) {
    memoryToken = fromStorage;
    writeCookieToken(fromStorage);
    return fromStorage;
  }

  const fromCookie = readCookieToken();
  if (fromCookie) {
    memoryToken = fromCookie;
    localStorage.setItem(TOKEN_STORAGE_KEY, fromCookie);
    return fromCookie;
  }

  return null;
}

export function clearAccessToken(): void {
  memoryToken = null;
  if (typeof window === 'undefined') return;
  localStorage.removeItem(TOKEN_STORAGE_KEY);
  clearCookieToken();
}

export function hasAccessToken(): boolean {
  return !!getAccessToken();
}

export const AUTH_SESSION_EXPIRED_EVENT = 'auth:session-expired';

export function emitSessionExpired(): void {
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent(AUTH_SESSION_EXPIRED_EVENT));
  }
}
