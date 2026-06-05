import { notFound } from "next/navigation";
import { getTranslations } from "next-intl/server";
import { ArticleHeader } from "@/components/blog/article-header";
import { ArticleContent } from "@/components/blog/article-content";
import { AuthorCard } from "@/components/blog/author-card";
import { RelatedArticles } from "@/components/blog/related-articles";
import { Header } from "@/components/blog/header";
import { fetchArticleById, fetchArticles } from "@/lib/api-server";
import { toArticleCard } from "@/lib/article-utils";

interface ArticlePageProps {
  params: Promise<{ id: string }>;
}

export async function generateMetadata({ params }: ArticlePageProps) {
  const { id } = await params;
  const t = await getTranslations("metadata");
  const article = await fetchArticleById(id);

  if (!article) {
    return { title: t("articleNotFound") };
  }

  const tCommon = await getTranslations("common");
  return {
    title: `${article.title} — ${tCommon("brand")}`,
    description: article.excerpt,
  };
}

export default async function ArticlePage({ params }: ArticlePageProps) {
  const { id } = await params;
  const apiArticle = await fetchArticleById(id);

  if (!apiArticle) {
    notFound();
  }

  const article = toArticleCard(apiArticle);
  const listResponse = await fetchArticles(1, 20);
  const relatedArticles = (listResponse?.data ?? [])
    .filter((a) => a.id !== apiArticle.id)
    .filter((a) => a.tags?.some((tag) => apiArticle.tags?.includes(tag)))
    .slice(0, 3)
    .map(toArticleCard);

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="pb-20">
        <ArticleHeader article={article} />
        <ArticleContent content={apiArticle.content} />
        <div className="mx-auto max-w-2xl px-4">
          <AuthorCard
            article={article}
            bio={apiArticle.author?.bio}
            followers={apiArticle.author?.followers}
          />
          {relatedArticles.length > 0 && <RelatedArticles articles={relatedArticles} />}
        </div>
      </main>
    </div>
  );
}
