import { defaultLocale, isValidLocale, type Locale } from '@/i18n/config';

const LOCALE_COOKIE = 'locale';
const LOCALE_STORAGE = 'user_locale';

export function getStoredLocale(): Locale {
  if (typeof window === 'undefined') return defaultLocale;

  const fromStorage = localStorage.getItem(LOCALE_STORAGE);
  if (fromStorage && isValidLocale(fromStorage)) return fromStorage;

  const match = document.cookie.match(new RegExp(`(?:^|; )${LOCALE_COOKIE}=([^;]*)`));
  const fromCookie = match?.[1];
  if (fromCookie && isValidLocale(fromCookie)) return fromCookie;

  return defaultLocale;
}

export function setLocalePreference(locale: Locale) {
  if (typeof window === 'undefined') return;

  localStorage.setItem(LOCALE_STORAGE, locale);
  document.cookie = `${LOCALE_COOKIE}=${locale};path=/;max-age=31536000;SameSite=Lax`;
}

export function syncLocaleFromProfile(locale: string | undefined) {
  if (!locale || !isValidLocale(locale)) return;

  const current = getStoredLocale();
  if (current !== locale) {
    setLocalePreference(locale);
    window.location.reload();
  }
}
