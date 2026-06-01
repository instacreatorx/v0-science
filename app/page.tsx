import { getTranslations } from "next-intl/server";
import { Header } from "@/components/blog/header";
import { TopicNav } from "@/components/blog/topic-nav";
import { ArticleCard } from "@/components/blog/article-card";
import { Sidebar } from "@/components/blog/sidebar";
import { articles, featuredArticle } from "@/lib/articles-data";

export default async function BlogPage() {
  const t = await getTranslations("home");
  const tCommon = await getTranslations("common");
  const tSidebar = await getTranslations("sidebar");

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <TopicNav />

      <main className="mx-auto max-w-5xl px-4 py-8">
        <div className="flex gap-12">
          <div className="flex-1">
            <section className="mb-10">
              <ArticleCard article={featuredArticle} variant="featured" />
            </section>

            <section>
              <h2 className="sr-only">{t("latestArticles")}</h2>
              <div className="divide-y divide-border">
                {articles.map((article) => (
                  <ArticleCard key={article.id} article={article} />
                ))}
              </div>
            </section>

            <div className="mt-8 text-center">
              <button className="rounded-full border border-border px-6 py-2 text-sm font-medium text-foreground transition-colors hover:bg-secondary">
                {tCommon("showMore")}
              </button>
            </div>
          </div>

          <Sidebar />
        </div>
      </main>

      <footer className="border-t border-border bg-background py-8">
        <div className="mx-auto flex max-w-5xl flex-col items-center justify-between gap-4 px-4 sm:flex-row">
          <span className="font-serif text-lg font-bold text-foreground">{tCommon("brand")}</span>
          <nav className="flex flex-wrap justify-center gap-4 text-sm text-muted-foreground">
            <a href="/about" className="hover:text-foreground">
              {tSidebar("about")}
            </a>
            <a href="/help" className="hover:text-foreground">
              {tSidebar("help")}
            </a>
            <a href="/terms" className="hover:text-foreground">
              {tSidebar("terms")}
            </a>
            <a href="/privacy" className="hover:text-foreground">
              {tSidebar("privacy")}
            </a>
          </nav>
        </div>
      </footer>
    </div>
  );
}
