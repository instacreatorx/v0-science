'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { Phone, ArrowRight, Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { authApi } from '@/lib/api';
import { useAuth } from '@/lib/auth-context';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { InputOTP, InputOTPGroup, InputOTPSlot } from '@/components/ui/input-otp';

interface PhoneAuthFormProps {
  mode: 'sign-in' | 'sign-up';
}

export function PhoneAuthForm({ mode }: PhoneAuthFormProps) {
  const t = useTranslations('auth');
  const tCommon = useTranslations('common');
  const router = useRouter();
  const { login } = useAuth();

  const [step, setStep] = useState<'phone' | 'otp'>('phone');
  const [phone, setPhone] = useState('');
  const [name, setName] = useState('');
  const [otp, setOtp] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [resendTimer, setResendTimer] = useState(0);

  useEffect(() => {
    if (resendTimer <= 0) return;
    const interval = setInterval(() => setResendTimer((prev) => prev - 1), 1000);
    return () => clearInterval(interval);
  }, [resendTimer]);

  const validatePhone = (value: string) => {
    const cleaned = value.replace(/\s/g, '');
    return cleaned.length >= 10 && /^\+?[\d\s-]+$/.test(value);
  };

  const handleSendOtp = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validatePhone(phone)) {
      toast.error(t('invalidPhone'));
      return;
    }

    if (mode === 'sign-up' && !name.trim()) {
      toast.error(t('nameRequired'));
      return;
    }

    setIsLoading(true);
    try {
      const response = await authApi.sendOtp({
        phone: phone.replace(/\s/g, ''),
        name: mode === 'sign-up' ? name.trim() : undefined,
      });
      setStep('otp');
      setResendTimer(response.expires_in || 60);
      toast.success(t('otpSent', { phone }));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('authError'));
    } finally {
      setIsLoading(false);
    }
  };

  const handleVerifyOtp = useCallback(
    async (code: string) => {
      if (code.length !== 6) return;

      setIsLoading(true);
      try {
        const response = await authApi.verifyOtp({
          phone: phone.replace(/\s/g, ''),
          code,
          name: mode === 'sign-up' ? name.trim() : undefined,
        });
        login(response.token, response.user);
        router.push('/');
        router.refresh();
      } catch (error) {
        toast.error(error instanceof Error ? error.message : t('authError'));
        setOtp('');
      } finally {
        setIsLoading(false);
      }
    },
    [phone, name, mode, login, router, t]
  );

  useEffect(() => {
    if (otp.length === 6) {
      handleVerifyOtp(otp);
    }
  }, [otp, handleVerifyOtp]);

  const handleResend = async () => {
    if (resendTimer > 0) return;
    setIsLoading(true);
    try {
      const response = await authApi.sendOtp({
        phone: phone.replace(/\s/g, ''),
        name: mode === 'sign-up' ? name.trim() : undefined,
      });
      setResendTimer(response.expires_in || 60);
      toast.success(t('otpSent', { phone }));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('authError'));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md space-y-8">
      <div className="text-center">
        <Link href="/" className="mb-8 inline-block font-serif text-3xl font-bold text-foreground">
          {tCommon('brand')}
        </Link>
        <h1 className="font-serif text-2xl font-bold text-foreground">
          {mode === 'sign-in' ? t('signInTitle') : t('signUpTitle')}
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          {mode === 'sign-in' ? t('signInSubtitle') : t('signUpSubtitle')}
        </p>
      </div>

      {step === 'phone' ? (
        <form onSubmit={handleSendOtp} className="space-y-5">
          {mode === 'sign-up' && (
            <div className="space-y-2">
              <Label htmlFor="name">{t('nameLabel')}</Label>
              <Input
                id="name"
                type="text"
                placeholder={t('namePlaceholder')}
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="h-12 rounded-xl"
                autoComplete="name"
              />
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="phone">{t('phoneLabel')}</Label>
            <div className="relative">
              <Phone className="absolute start-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                id="phone"
                type="tel"
                placeholder={t('phonePlaceholder')}
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                className="h-12 rounded-xl ps-10"
                autoComplete="tel"
                dir="ltr"
              />
            </div>
          </div>

          <Button
            type="submit"
            disabled={isLoading}
            className="h-12 w-full rounded-xl bg-primary text-base font-medium text-primary-foreground hover:bg-primary/90"
          >
            {isLoading ? (
              <Loader2 className="h-5 w-5 animate-spin" />
            ) : (
              <>
                {t('sendCode')}
                <ArrowRight className="ms-2 h-4 w-4 rtl:rotate-180" />
              </>
            )}
          </Button>
        </form>
      ) : (
        <div className="space-y-6">
          <div className="space-y-2 text-center">
            <Label>{t('otpLabel')}</Label>
            <div className="flex justify-center" dir="ltr">
              <InputOTP maxLength={6} value={otp} onChange={setOtp} disabled={isLoading}>
                <InputOTPGroup>
                  <InputOTPSlot index={0} className="h-12 w-12 rounded-lg text-lg" />
                  <InputOTPSlot index={1} className="h-12 w-12 rounded-lg text-lg" />
                  <InputOTPSlot index={2} className="h-12 w-12 rounded-lg text-lg" />
                  <InputOTPSlot index={3} className="h-12 w-12 rounded-lg text-lg" />
                  <InputOTPSlot index={4} className="h-12 w-12 rounded-lg text-lg" />
                  <InputOTPSlot index={5} className="h-12 w-12 rounded-lg text-lg" />
                </InputOTPGroup>
              </InputOTP>
            </div>
            {process.env.NODE_ENV === 'development' && (
              <p className="text-xs text-muted-foreground">{t('devOtpHint')}</p>
            )}
          </div>

          {isLoading && (
            <div className="flex justify-center">
              <Loader2 className="h-6 w-6 animate-spin text-primary" />
            </div>
          )}

          <div className="text-center">
            <button
              type="button"
              onClick={handleResend}
              disabled={resendTimer > 0 || isLoading}
              className="text-sm text-primary hover:underline disabled:cursor-not-allowed disabled:text-muted-foreground disabled:no-underline"
            >
              {resendTimer > 0 ? t('resendIn', { seconds: resendTimer }) : t('resendCode')}
            </button>
          </div>

          <button
            type="button"
            onClick={() => {
              setStep('phone');
              setOtp('');
            }}
            className="w-full text-center text-sm text-muted-foreground hover:text-foreground"
          >
            ← {t('phoneLabel')}
          </button>
        </div>
      )}

      <p className="text-center text-sm text-muted-foreground">
        {mode === 'sign-in' ? t('noAccount') : t('hasAccount')}{' '}
        <Link
          href={mode === 'sign-in' ? '/sign-up' : '/sign-in'}
          className="font-medium text-primary hover:underline"
        >
          {mode === 'sign-in' ? t('signUpLink') : t('signInLink')}
        </Link>
      </p>
    </div>
  );
}
