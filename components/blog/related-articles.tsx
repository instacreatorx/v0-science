import Image from "next/image";
import Link from "next/link";
import { Article } from "@/components/blog/article-card";

interface RelatedArticlesProps {
  articles: Article[];
}

export function RelatedArticles({ articles }: RelatedArticlesProps) {
  return (
    <section className="mt-12 border-t border-border pt-8">
      <h2 className="mb-6 font-serif text-xl font-bold text-foreground">
        More from Thoughts
      </h2>
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {articles.map((article) => (
          <article key={article.id} className="group">
            <Link href={`/article/${article.id}`}>
              {article.image && (
                <div className="relative mb-3 aspect-[4/3] w-full overflow-hidden rounded bg-muted">
                  <Image
                    src={article.image}
                    alt={article.title}
                    fill
                    className="object-cover transition-transform duration-300 group-hover:scale-105"
                  />
                </div>
              )}
              <div className="mb-2 flex items-center gap-2">
                <div className="relative h-5 w-5 overflow-hidden rounded-full bg-muted">
                  <Image
                    src={article.author.avatar}
                    alt={article.author.name}
                    fill
                    className="object-cover"
                  />
                </div>
                <span className="text-xs font-medium text-foreground">
                  {article.author.name}
                </span>
              </div>
              <h3 className="mb-1 font-serif text-base font-bold leading-snug text-foreground group-hover:underline">
                {article.title}
              </h3>
              <p className="mb-2 line-clamp-2 text-sm text-muted-foreground">
                {article.excerpt}
              </p>
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <span>{article.date}</span>
                <span>·</span>
                <span>{article.readTime}</span>
                {article.isMemberOnly && (
                  <>
                    <span>·</span>
                    <span className="rounded bg-secondary px-1 py-0.5 text-[9px] font-medium uppercase text-secondary-foreground">
                      Member only
                    </span>
                  </>
                )}
              </div>
            </Link>
          </article>
        ))}
      </div>
    </section>
  );
}
