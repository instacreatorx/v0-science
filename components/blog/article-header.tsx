import Image from "next/image";
import Link from "next/link";
import { getTranslations } from "next-intl/server";
import { VerifiedBadge } from "@/components/blog/verified-badge";
import { ArticleActions } from "@/components/blog/article-actions";
import { FollowAuthorButton } from "@/components/blog/follow-author-button";
import type { ArticleCardData } from "@/lib/article-utils";

interface ArticleHeaderProps {
  article: ArticleCardData;
  isFollowing?: boolean;
}

export async function ArticleHeader({ article, isFollowing }: ArticleHeaderProps) {
  const t = await getTranslations("common");

  return (
    <header className="mx-auto max-w-3xl px-4 pt-10">
      <h1 className="mb-4 font-serif text-3xl font-bold leading-tight text-foreground md:text-4xl lg:text-5xl">
        {article.title}
      </h1>
      <p className="mb-6 text-lg text-muted-foreground md:text-xl">{article.excerpt}</p>

      <div className="mb-6 flex items-center gap-4">
        <Link href={`/author/${article.author.id}`} className="shrink-0">
          <div className="relative h-12 w-12 overflow-hidden rounded-full bg-muted">
            <Image
              src={article.author.avatar}
              alt={article.author.name}
              fill
              className="object-cover"
            />
          </div>
        </Link>
        <div className="flex-1">
          <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
            <Link href={`/author/${article.author.id}`} className="font-medium text-foreground hover:underline">
              {article.author.name}
            </Link>
            {article.team && (
              <>
                <span className="text-muted-foreground">{t("in")}</span>
                <Link href={`/team/${article.team.slug}`} className="flex items-center gap-1 font-medium text-foreground hover:underline">
                  {article.team.name}
                  <VerifiedBadge verified={!!article.team.is_verified} />
                </Link>
              </>
            )}
          </div>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <span>{article.readTime}</span>
            <span>·</span>
            <span>{article.date}</span>
            {article.isMemberOnly && (
              <>
                <span>·</span>
                <span className="rounded bg-secondary px-1.5 py-0.5 text-[10px] font-medium uppercase text-secondary-foreground">
                  {t("memberOnly")}
                </span>
              </>
            )}
          </div>
        </div>
        <FollowAuthorButton authorId={article.author.id} initialFollowing={isFollowing} />
      </div>

      <ArticleActions article={article} />

      {article.image && (
        <div className="relative mb-8 aspect-[16/9] w-full overflow-hidden rounded-sm bg-muted">
          <Image src={article.image} alt={article.title} fill className="object-cover" priority />
        </div>
      )}
    </header>
  );
}
