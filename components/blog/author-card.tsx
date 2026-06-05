import Image from "next/image";
import Link from "next/link";
import { getTranslations } from "next-intl/server";
import { FollowAuthorButton } from "@/components/blog/follow-author-button";
import type { ArticleCardData } from "@/lib/article-utils";

interface AuthorCardProps {
  article: ArticleCardData;
  bio?: string;
  followers?: number;
  isFollowing?: boolean;
}

export async function AuthorCard({ article, bio, followers, isFollowing }: AuthorCardProps) {
  const t = await getTranslations("common");

  return (
    <div className="mt-12 border-t border-border pt-8">
      <div className="flex gap-4">
        <Link href={`/author/${article.author.id}`} className="shrink-0">
          <div className="relative h-16 w-16 overflow-hidden rounded-full bg-muted">
            <Image
              src={article.author.avatar}
              alt={article.author.name}
              fill
              className="object-cover"
            />
          </div>
        </Link>
        <div className="flex-1">
          <div className="mb-1 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                {t("writtenBy")}
              </p>
              <Link href={`/author/${article.author.id}`} className="text-lg font-bold text-foreground hover:underline">
                {article.author.name}
              </Link>
            </div>
            <FollowAuthorButton authorId={article.author.id} initialFollowing={isFollowing} />
          </div>
          {followers !== undefined && (
            <p className="mb-2 text-sm text-muted-foreground">
              {followers} {t("followers")}
            </p>
          )}
          {bio && <p className="text-foreground">{bio}</p>}
        </div>
      </div>
    </div>
  );
}
