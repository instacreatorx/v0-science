import type { Article, Team, User } from '@/lib/api';

export interface ArticleCardData {
  id: string;
  numericId: number;
  slug: string;
  title: string;
  excerpt: string;
  content?: string;
  author: {
    id: number;
    name: string;
    avatar: string;
  };
  team?: Team | null;
  date: string;
  readTime: string;
  image?: string;
  tags: string[];
  isMemberOnly: boolean;
  likesCount: number;
  commentsCount: number;
  likedByMe?: boolean;
  bookmarkedByMe?: boolean;
  status?: string;
}

export function formatArticleDate(iso: string): string {
  try {
    return new Intl.DateTimeFormat(undefined, {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    }).format(new Date(iso));
  } catch {
    return iso;
  }
}

export function articleHref(article: { slug?: string; id: number | string }): string {
  const slug = 'slug' in article && article.slug ? article.slug : String(article.id);
  return `/article/${slug}`;
}

export function toArticleCard(article: Article): ArticleCardData {
  const author: User | undefined = article.author;
  return {
    id: article.slug || String(article.id),
    numericId: article.id,
    slug: article.slug || String(article.id),
    title: article.title,
    excerpt: article.excerpt,
    content: article.content,
    author: {
      id: article.author_id,
      name: author?.name ?? 'Unknown',
      avatar: author?.avatar ?? '',
    },
    team: article.team ?? null,
    date: formatArticleDate(article.published_at || article.created_at),
    readTime: article.read_time,
    image: article.image || undefined,
    tags: article.tags ?? [],
    isMemberOnly: article.is_member_only,
    likesCount: article.likes_count,
    commentsCount: article.comments_count,
    likedByMe: article.liked_by_me,
    bookmarkedByMe: article.bookmarked_by_me,
    status: article.status,
  };
}
