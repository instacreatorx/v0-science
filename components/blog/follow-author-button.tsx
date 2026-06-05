'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { userApi } from '@/lib/api';
import { useAuth } from '@/lib/auth/auth-context';

interface FollowAuthorButtonProps {
  authorId: number;
  initialFollowing?: boolean;
}

export function FollowAuthorButton({ authorId, initialFollowing }: FollowAuthorButtonProps) {
  const { isAuthenticated } = useAuth();
  const router = useRouter();
  const [following, setFollowing] = useState(initialFollowing ?? false);

  const handleFollow = async () => {
    if (!isAuthenticated) {
      router.push('/sign-in');
      return;
    }
    try {
      const result = await userApi.followUser(authorId);
      setFollowing(result.following);
      toast.success(result.message);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to follow');
    }
  };

  return (
    <Button
      variant="outline"
      size="sm"
      className="hidden border-primary text-primary hover:bg-primary hover:text-primary-foreground sm:inline-flex"
      onClick={handleFollow}
    >
      {following ? 'Following' : 'Follow'}
    </Button>
  );
}
