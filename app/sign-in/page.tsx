import { PhoneAuthForm } from '@/components/auth/phone-auth-form';

export default function SignInPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-background via-background to-secondary/30 px-4">
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -start-40 -top-40 h-80 w-80 rounded-full bg-primary/5 blur-3xl" />
        <div className="absolute -bottom-40 -end-40 h-80 w-80 rounded-full bg-primary/10 blur-3xl" />
      </div>
      <div className="relative z-10 w-full max-w-md rounded-2xl border border-border/50 bg-card/80 p-8 shadow-xl backdrop-blur-sm">
        <PhoneAuthForm mode="sign-in" />
      </div>
    </div>
  );
}
