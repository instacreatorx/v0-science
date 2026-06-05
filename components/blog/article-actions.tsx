'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Bookmark, MessageCircle, Share, MoreHorizontal, ThumbsUp } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { articlesApi } from '@/lib/api';
import { useAuth } from '@/lib/auth/auth-context';
import type { ArticleCardData } from '@/lib/article-utils';

interface ArticleActionsProps {
  article: ArticleCardData;
}

export function ArticleActions({ article }: ArticleActionsProps) {
  const { isAuthenticated } = useAuth();
  const router = useRouter();
  const [liked, setLiked] = useState(article.likedByMe ?? false);
  const [likesCount, setLikesCount] = useState(article.likesCount);
  const [bookmarked, setBookmarked] = useState(article.bookmarkedByMe ?? false);

  const requireAuth = () => {
    toast.error('Please sign in to continue');
    router.push('/sign-in');
    return false;
  };

  const handleLike = async () => {
    if (!isAuthenticated && !requireAuth()) return;
    try {
      const result = await articlesApi.like(article.numericId);
      setLiked(result.liked);
      setLikesCount((prev) => (result.liked ? prev + 1 : Math.max(prev - 1, 0)));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to like');
    }
  };

  const handleBookmark = async () => {
    if (!isAuthenticated && !requireAuth()) return;
    try {
      const result = await articlesApi.bookmark(article.numericId);
      setBookmarked(result.bookmarked);
      toast.success(result.message);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to bookmark');
    }
  };

  return (
    <div className="mb-8 flex items-center justify-between border-y border-border py-3">
      <div className="flex items-center gap-4">
        <button
          type="button"
          onClick={handleLike}
          className={`flex items-center gap-1.5 transition-colors ${liked ? 'text-primary' : 'text-muted-foreground hover:text-foreground'}`}
        >
          <ThumbsUp className="h-5 w-5" />
          <span className="text-sm">{likesCount}</span>
        </button>
        <div className="flex items-center gap-1.5 text-muted-foreground">
          <MessageCircle className="h-5 w-5" />
          <span className="text-sm">{article.commentsCount}</span>
        </div>
      </div>
      <div className="flex items-center gap-1">
        <Button
          variant="ghost"
          size="icon"
          className={`h-9 w-9 ${bookmarked ? 'text-primary' : 'text-muted-foreground hover:text-foreground'}`}
          onClick={handleBookmark}
        >
          <Bookmark className="h-5 w-5" />
        </Button>
        <Button variant="ghost" size="icon" className="h-9 w-9 text-muted-foreground hover:text-foreground">
          <Share className="h-5 w-5" />
        </Button>
        <Button variant="ghost" size="icon" className="h-9 w-9 text-muted-foreground hover:text-foreground">
          <MoreHorizontal className="h-5 w-5" />
        </Button>
      </div>
    </div>
  );
}
