import { Header } from "@/components/blog/header";
import { TopicNav } from "@/components/blog/topic-nav";
import { ArticleCard } from "@/components/blog/article-card";
import { Sidebar } from "@/components/blog/sidebar";
import { articles, featuredArticle } from "@/lib/articles-data";

export default function BlogPage() {
  return (
    <div className="min-h-screen bg-background">
      <Header />
      <TopicNav />

      <main className="mx-auto max-w-5xl px-4 py-8">
        <div className="flex gap-12">
          {/* Main Content */}
          <div className="flex-1">
            {/* Featured Article */}
            <section className="mb-10">
              <ArticleCard article={featuredArticle} variant="featured" />
            </section>

            {/* Article Feed */}
            <section>
              <h2 className="sr-only">Latest Articles</h2>
              <div className="divide-y divide-border">
                {articles.map((article) => (
                  <ArticleCard key={article.id} article={article} />
                ))}
              </div>
            </section>

            {/* Load More */}
            <div className="mt-8 text-center">
              <button className="rounded-full border border-border px-6 py-2 text-sm font-medium text-foreground transition-colors hover:bg-secondary">
                Show more
              </button>
            </div>
          </div>

          {/* Sidebar */}
          <Sidebar />
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-border bg-background py-8">
        <div className="mx-auto flex max-w-5xl flex-col items-center justify-between gap-4 px-4 sm:flex-row">
          <span className="font-serif text-lg font-bold text-foreground">Thoughts</span>
          <nav className="flex flex-wrap justify-center gap-4 text-sm text-muted-foreground">
            <a href="/about" className="hover:text-foreground">
              About
            </a>
            <a href="/help" className="hover:text-foreground">
              Help
            </a>
            <a href="/terms" className="hover:text-foreground">
              Terms
            </a>
            <a href="/privacy" className="hover:text-foreground">
              Privacy
            </a>
          </nav>
        </div>
      </footer>
    </div>
  );
}
