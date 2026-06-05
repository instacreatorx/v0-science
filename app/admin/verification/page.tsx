'use client';

import { useCallback, useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Header } from '@/components/blog/header';
import { ProtectedRoute } from '@/components/auth/protected-route';
import { VerifiedBadge } from '@/components/blog/verified-badge';
import { Button } from '@/components/ui/button';
import { adminApi, type TeamVerificationRequest } from '@/lib/api';
import { useAuth } from '@/lib/auth/auth-context';
import { useRouter } from 'next/navigation';

export default function AdminVerificationPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const [requests, setRequests] = useState<TeamVerificationRequest[]>([]);

  const load = useCallback(async () => {
    try {
      const response = await adminApi.listVerificationRequests('pending');
      setRequests(response.data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to load requests');
    }
  }, []);

  useEffect(() => {
    if (!isLoading && user?.role !== 'super_admin') {
      router.replace('/');
      return;
    }
    if (user?.role === 'super_admin') {
      void load();
    }
  }, [user, isLoading, load, router]);

  const handleApprove = async (id: number) => {
    try {
      await adminApi.approveVerification(id);
      toast.success('Team verified');
      void load();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to approve');
    }
  };

  const handleReject = async (id: number) => {
    const reason = window.prompt('Rejection reason');
    if (!reason?.trim()) return;
    try {
      await adminApi.rejectVerification(id, reason.trim());
      toast.success('Request rejected');
      void load();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to reject');
    }
  };

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-background">
        <Header />
        <main className="mx-auto max-w-3xl px-4 py-8">
          <h1 className="mb-8 font-serif text-3xl font-bold">Verification requests</h1>
          {requests.length === 0 ? (
            <p className="text-muted-foreground">No pending requests.</p>
          ) : (
            <ul className="space-y-4">
              {requests.map((req) => (
                <li key={req.id} className="rounded-lg border border-border p-4">
                  <div className="mb-2 flex items-center gap-2 font-semibold">
                    {req.team?.name}
                    <VerifiedBadge verified={false} />
                  </div>
                  <p className="mb-4 text-sm text-foreground">{req.proof_text}</p>
                  <div className="flex gap-2">
                    <Button size="sm" onClick={() => void handleApprove(req.id)}>
                      Approve (blue tick)
                    </Button>
                    <Button size="sm" variant="outline" onClick={() => void handleReject(req.id)}>
                      Reject
                    </Button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </main>
      </div>
    </ProtectedRoute>
  );
}
