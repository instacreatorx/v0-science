'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'sonner';
import { Header } from '@/components/blog/header';
import { ProtectedRoute } from '@/components/auth/protected-route';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { authApi } from '@/lib/api';
import { useAuth } from '@/lib/auth/auth-context';
import { setLocalePreference } from '@/lib/locale';
import { locales, localeNames, type Locale } from '@/i18n/config';
import { Loader2 } from 'lucide-react';

function SettingsContent() {
  const t = useTranslations('settings');
  const { user, refreshUser } = useAuth();
  const [selectedLocale, setSelectedLocale] = useState<Locale>('en');
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (user?.locale && (user.locale === 'en' || user.locale === 'fa')) {
      setSelectedLocale(user.locale as Locale);
    }
  }, [user?.locale]);

  const handleSave = async () => {
    setIsSaving(true);
    try {
      await authApi.updateUser({ locale: selectedLocale });
      setLocalePreference(selectedLocale);
      await refreshUser();
      toast.success(t('saved'));
      window.location.reload();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to save');
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-lg px-4 py-12">
        <h1 className="mb-8 font-serif text-3xl font-bold text-foreground">{t('title')}</h1>

        <div className="space-y-6 rounded-xl border border-border bg-card p-6">
          <div className="space-y-2">
            <Label>{t('language')}</Label>
            <p className="text-sm text-muted-foreground">{t('languageDescription')}</p>
            <Select
              value={selectedLocale}
              onValueChange={(value) => setSelectedLocale(value as Locale)}
            >
              <SelectTrigger className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {locales.map((locale) => (
                  <SelectItem key={locale} value={locale}>
                    {localeNames[locale]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <Button onClick={handleSave} disabled={isSaving} className="w-full rounded-full">
            {isSaving ? <Loader2 className="h-4 w-4 animate-spin" /> : t('save')}
          </Button>
        </div>
      </main>
    </div>
  );
}

export default function SettingsPage() {
  return (
    <ProtectedRoute>
      <SettingsContent />
    </ProtectedRoute>
  );
}
