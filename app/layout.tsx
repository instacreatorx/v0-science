import type { Metadata } from 'next';
import { Source_Serif_4, Inter, Vazirmatn } from 'next/font/google';
import { Analytics } from '@vercel/analytics/next';
import { getLocale, getMessages, getTranslations } from 'next-intl/server';
import { isRtlLocale } from '@/i18n/config';
import { Providers } from '@/components/providers';
import './globals.css';

const sourceSerif = Source_Serif_4({
  subsets: ['latin'],
  variable: '--font-serif',
  display: 'swap',
});

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-sans',
  display: 'swap',
});

const vazirmatn = Vazirmatn({
  subsets: ['arabic'],
  variable: '--font-fa',
  display: 'swap',
});

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations('metadata');

  return {
    title: t('title'),
    description: t('description'),
    generator: 'v0.app',
    icons: {
      icon: [
        {
          url: '/icon-light-32x32.png',
          media: '(prefers-color-scheme: light)',
        },
        {
          url: '/icon-dark-32x32.png',
          media: '(prefers-color-scheme: dark)',
        },
        {
          url: '/icon.svg',
          type: 'image/svg+xml',
        },
      ],
      apple: '/apple-icon.png',
    },
  };
}

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const locale = await getLocale();
  const messages = await getMessages();
  const rtl = isRtlLocale(locale);

  return (
    <html
      lang={locale}
      dir={rtl ? 'rtl' : 'ltr'}
      className={`${sourceSerif.variable} ${inter.variable} ${vazirmatn.variable} bg-background`}
    >
      <body className={`font-sans antialiased ${rtl ? 'font-[family-name:var(--font-fa)]' : ''}`}>
        <Providers locale={locale} messages={messages}>
          {children}
        </Providers>
        {process.env.NODE_ENV === 'production' && <Analytics />}
      </body>
    </html>
  );
}
