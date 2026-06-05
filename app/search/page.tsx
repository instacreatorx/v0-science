import { Header } from '@/components/blog/header';
import { ArticleCard } from '@/components/blog/article-card';
import { fetchSearchArticles } from '@/lib/api-server';
import { toArticleCard } from '@/lib/article-utils';

interface SearchPageProps {
  searchParams: Promise<{ q?: string }>;
}

export default async function SearchPage({ searchParams }: SearchPageProps) {
  const { q } = await searchParams;
  const query = q?.trim() ?? '';
  const response = query ? await fetchSearchArticles(query) : null;
  const articles = (response?.data ?? []).map(toArticleCard);

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-3xl px-4 py-8">
        <h1 className="mb-6 font-serif text-2xl font-bold">
          {query ? `Results for "${query}"` : 'Search'}
        </h1>
        {!query ? (
          <p className="text-muted-foreground">Enter a search term from the header.</p>
        ) : articles.length === 0 ? (
          <p className="text-muted-foreground">No articles found.</p>
        ) : (
          <div className="divide-y divide-border">
            {articles.map((article) => (
              <ArticleCard key={article.numericId} article={article} />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
