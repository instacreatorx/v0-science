'use client';

import { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { toast } from 'sonner';
import { Header } from '@/components/blog/header';
import { ProtectedRoute } from '@/components/auth/protected-route';
import { Button } from '@/components/ui/button';
import { articlesApi, type Article } from '@/lib/api';
import { articleHref } from '@/lib/article-utils';

type Tab = 'draft' | 'published' | 'archived';

export default function MyStoriesPage() {
  const [tab, setTab] = useState<Tab>('draft');
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const response = await articlesApi.myArticles({ status: tab, per_page: 50 });
      setArticles(response.data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to load stories');
    } finally {
      setLoading(false);
    }
  }, [tab]);

  useEffect(() => {
    void load();
  }, [load]);

  const handlePublish = async (id: number) => {
    try {
      await articlesApi.publish(id);
      toast.success('Published');
      void load();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to publish');
    }
  };

  const handleArchive = async (id: number) => {
    try {
      await articlesApi.archive(id);
      toast.success('Archived');
      void load();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to archive');
    }
  };

  const tabs: { key: Tab; label: string }[] = [
    { key: 'draft', label: 'Drafts' },
    { key: 'published', label: 'Published' },
    { key: 'archived', label: 'Archived' },
  ];

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-background">
        <Header />
        <main className="mx-auto max-w-3xl px-4 py-8">
          <div className="mb-8 flex items-center justify-between">
            <h1 className="font-serif text-3xl font-bold">My stories</h1>
            <Button asChild>
              <Link href="/write">New story</Link>
            </Button>
          </div>

          <div className="mb-6 flex gap-2 border-b border-border">
            {tabs.map(({ key, label }) => (
              <button
                key={key}
                type="button"
                onClick={() => setTab(key)}
                className={`border-b-2 px-4 py-2 text-sm font-medium transition-colors ${
                  tab === key
                    ? 'border-primary text-foreground'
                    : 'border-transparent text-muted-foreground hover:text-foreground'
                }`}
              >
                {label}
              </button>
            ))}
          </div>

          {loading ? (
            <p className="text-muted-foreground">Loading...</p>
          ) : articles.length === 0 ? (
            <p className="text-muted-foreground">No {tab} stories yet.</p>
          ) : (
            <ul className="divide-y divide-border">
              {articles.map((article) => (
                <li key={article.id} className="flex items-center justify-between gap-4 py-4">
                  <div>
                    <Link
                      href={tab === 'draft' ? `/write/${article.id}` : articleHref(article)}
                      className="font-serif text-lg font-bold hover:underline"
                    >
                      {article.title || 'Untitled'}
                    </Link>
                    <p className="text-sm text-muted-foreground">{article.excerpt || 'No excerpt'}</p>
                  </div>
                  <div className="flex shrink-0 gap-2">
                    {tab === 'draft' && (
                      <>
                        <Button size="sm" variant="outline" asChild>
                          <Link href={`/write/${article.id}`}>Edit</Link>
                        </Button>
                        <Button size="sm" onClick={() => void handlePublish(article.id)}>
                          Publish
                        </Button>
                      </>
                    )}
                    {tab === 'published' && (
                      <Button size="sm" variant="outline" onClick={() => void handleArchive(article.id)}>
                        Archive
                      </Button>
                    )}
                  </div>
                </li>
              ))}
            </ul>
          )}
        </main>
      </div>
    </ProtectedRoute>
  );
}
