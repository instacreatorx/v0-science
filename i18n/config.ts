export const locales = ['en', 'fa'] as const;
export type Locale = (typeof locales)[number];

export const defaultLocale: Locale = 'en';

export const rtlLocales: Locale[] = ['fa'];

export function isRtlLocale(locale: string): boolean {
  return rtlLocales.includes(locale as Locale);
}

export function isValidLocale(locale: string): locale is Locale {
  return locales.includes(locale as Locale);
}

export const localeNames: Record<Locale, string> = {
  en: 'English',
  fa: 'فارسی',
};
