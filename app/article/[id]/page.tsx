import { notFound } from "next/navigation";
import { articles, featuredArticle, articleContent } from "@/lib/articles-data";
import { ArticleHeader } from "@/components/blog/article-header";
import { ArticleContent } from "@/components/blog/article-content";
import { AuthorCard } from "@/components/blog/author-card";
import { RelatedArticles } from "@/components/blog/related-articles";
import { Header } from "@/components/blog/header";

interface ArticlePageProps {
  params: Promise<{ id: string }>;
}

export async function generateMetadata({ params }: ArticlePageProps) {
  const { id } = await params;
  const allArticles = [featuredArticle, ...articles];
  const article = allArticles.find((a) => a.id === id);

  if (!article) {
    return {
      title: "Article Not Found",
    };
  }

  return {
    title: `${article.title} — Thoughts`,
    description: article.excerpt,
  };
}

export default async function ArticlePage({ params }: ArticlePageProps) {
  const { id } = await params;
  const allArticles = [featuredArticle, ...articles];
  const article = allArticles.find((a) => a.id === id);

  if (!article) {
    notFound();
  }

  const content = articleContent[id] || "";
  const relatedArticles = articles
    .filter((a) => a.id !== id && a.tags.some((tag) => article.tags.includes(tag)))
    .slice(0, 3);

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="pb-20">
        <ArticleHeader article={article} />
        <ArticleContent content={content} />
        <div className="mx-auto max-w-2xl px-4">
          <AuthorCard article={article} />
          {relatedArticles.length > 0 && (
            <RelatedArticles articles={relatedArticles} />
          )}
        </div>
      </main>
    </div>
  );
}
