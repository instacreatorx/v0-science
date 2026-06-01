import Image from "next/image";
import Link from "next/link";
import { getTranslations } from "next-intl/server";
import { Article } from "@/components/blog/article-card";
import { Button } from "@/components/ui/button";
import { Bookmark, Share, MoreHorizontal, MessageCircle, ThumbsUp } from "lucide-react";

interface ArticleHeaderProps {
  article: Article;
}

export async function ArticleHeader({ article }: ArticleHeaderProps) {
  const t = await getTranslations("common");

  return (
    <header className="mx-auto max-w-3xl px-4 pt-10">
      <h1 className="mb-4 font-serif text-3xl font-bold leading-tight text-foreground md:text-4xl lg:text-5xl">
        {article.title}
      </h1>
      <p className="mb-6 text-lg text-muted-foreground md:text-xl">{article.excerpt}</p>

      <div className="mb-6 flex items-center gap-4">
        <Link href="#" className="shrink-0">
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
            <Link href="#" className="font-medium text-foreground hover:underline">
              {article.author.name}
            </Link>
            {article.publication && (
              <>
                <span className="text-muted-foreground">{t("in")}</span>
                <Link href="#" className="font-medium text-foreground hover:underline">
                  {article.publication}
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
        <Button
          variant="outline"
          size="sm"
          className="hidden border-primary text-primary hover:bg-primary hover:text-primary-foreground sm:inline-flex"
        >
          {t("follow")}
        </Button>
      </div>

      <div className="mb-8 flex items-center justify-between border-y border-border py-3">
        <div className="flex items-center gap-4">
          <button className="flex items-center gap-1.5 text-muted-foreground transition-colors hover:text-foreground">
            <ThumbsUp className="h-5 w-5" />
            <span className="text-sm">1.2K</span>
          </button>
          <button className="flex items-center gap-1.5 text-muted-foreground transition-colors hover:text-foreground">
            <MessageCircle className="h-5 w-5" />
            <span className="text-sm">48</span>
          </button>
        </div>
        <div className="flex items-center gap-1">
          <Button variant="ghost" size="icon" className="h-9 w-9 text-muted-foreground hover:text-foreground">
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

      {article.image && (
        <div className="relative mb-8 aspect-[16/9] w-full overflow-hidden rounded-sm bg-muted">
          <Image src={article.image} alt={article.title} fill className="object-cover" priority />
        </div>
      )}
    </header>
  );
}
