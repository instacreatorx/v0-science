import Image from "next/image";
import Link from "next/link";
import { getTranslations } from "next-intl/server";
import { Bookmark, MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";

export interface Article {
  id: string;
  title: string;
  excerpt: string;
  author: {
    name: string;
    avatar: string;
  };
  publication?: string;
  date: string;
  readTime: string;
  image?: string;
  tags: string[];
  isMemberOnly?: boolean;
}

interface ArticleCardProps {
  article: Article;
  variant?: "default" | "featured";
}

export async function ArticleCard({ article, variant = "default" }: ArticleCardProps) {
  const t = await getTranslations("common");

  if (variant === "featured") {
    return (
      <article className="group relative overflow-hidden rounded-lg">
        <Link href={`/article/${article.id}`} className="block">
          <div className="relative aspect-[16/9] w-full overflow-hidden rounded-lg bg-muted">
            {article.image && (
              <Image
                src={article.image}
                alt={article.title}
                fill
                className="object-cover transition-transform duration-300 group-hover:scale-105"
              />
            )}
            <div className="absolute inset-0 bg-gradient-to-t from-foreground/80 via-foreground/20 to-transparent" />
            <div className="absolute bottom-0 start-0 end-0 p-6 text-card">
              <div className="mb-3 flex items-center gap-2">
                <div className="relative h-6 w-6 overflow-hidden rounded-full bg-muted">
                  <Image
                    src={article.author.avatar}
                    alt={article.author.name}
                    fill
                    className="object-cover"
                  />
                </div>
                <span className="text-sm font-medium">{article.author.name}</span>
                {article.publication && (
                  <>
                    <span className="text-sm opacity-60">{t("in")}</span>
                    <span className="text-sm font-medium">{article.publication}</span>
                  </>
                )}
              </div>
              <h2 className="mb-2 font-serif text-2xl font-bold leading-tight md:text-3xl">
                {article.title}
              </h2>
              <p className="mb-3 line-clamp-2 text-sm opacity-80 md:text-base">{article.excerpt}</p>
              <div className="flex items-center gap-2 text-xs opacity-60">
                <span>{article.date}</span>
                <span>·</span>
                <span>{article.readTime}</span>
                {article.isMemberOnly && (
                  <>
                    <span>·</span>
                    <span className="rounded bg-card/20 px-1.5 py-0.5 text-[10px] font-medium uppercase">
                      {t("memberOnly")}
                    </span>
                  </>
                )}
              </div>
            </div>
          </div>
        </Link>
      </article>
    );
  }

  return (
    <article className="group flex gap-4 border-b border-border py-6 last:border-0 md:gap-6">
      <div className="flex flex-1 flex-col">
        <div className="mb-2 flex items-center gap-2">
          <div className="relative h-5 w-5 overflow-hidden rounded-full bg-muted">
            <Image
              src={article.author.avatar}
              alt={article.author.name}
              fill
              className="object-cover"
            />
          </div>
          <span className="text-sm font-medium text-foreground">{article.author.name}</span>
          {article.publication && (
            <>
              <span className="text-sm text-muted-foreground">{t("in")}</span>
              <span className="text-sm font-medium text-foreground">{article.publication}</span>
            </>
          )}
        </div>
        <Link href={`/article/${article.id}`} className="block">
          <h3 className="mb-1 font-serif text-lg font-bold leading-snug text-foreground group-hover:underline md:text-xl">
            {article.title}
          </h3>
          <p className="mb-3 line-clamp-2 text-sm text-muted-foreground md:text-base">
            {article.excerpt}
          </p>
        </Link>
        <div className="mt-auto flex items-center justify-between">
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <span>{article.date}</span>
            <span>·</span>
            <span>{article.readTime}</span>
            {article.isMemberOnly && (
              <>
                <span>·</span>
                <span className="rounded bg-secondary px-1.5 py-0.5 text-[10px] font-medium uppercase text-secondary-foreground">
                  {t("memberOnly")}
                </span>
              </>
            )}
          </div>
          <div className="flex items-center gap-1">
            <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground">
              <Bookmark className="h-4 w-4" />
            </Button>
            <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground">
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
      {article.image && (
        <Link href={`/article/${article.id}`} className="shrink-0">
          <div className="relative h-20 w-28 overflow-hidden rounded bg-muted md:h-28 md:w-40">
            <Image
              src={article.image}
              alt={article.title}
              fill
              className="object-cover transition-transform duration-300 group-hover:scale-105"
            />
          </div>
        </Link>
      )}
    </article>
  );
}
