'use client';

import { useCallback, useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { Header } from '@/components/blog/header';
import { ProtectedRoute } from '@/components/auth/protected-route';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { articlesApi, type Article } from '@/lib/api';

export default function EditDraftPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const [article, setArticle] = useState<Article | null>(null);
  const [title, setTitle] = useState('');
  const [excerpt, setExcerpt] = useState('');
  const [content, setContent] = useState('');
  const [tags, setTags] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [isPublishing, setIsPublishing] = useState(false);

  const load = useCallback(async () => {
    try {
      const data = await articlesApi.getById(params.id);
      setArticle(data);
      setTitle(data.title);
      setExcerpt(data.excerpt);
      setContent(data.content);
      setTags((data.tags ?? []).join(', '));
    } catch {
      toast.error('Failed to load draft');
      router.push('/me/stories');
    }
  }, [params.id, router]);

  useEffect(() => {
    void load();
  }, [load]);

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const updated = await articlesApi.update(params.id, {
        title,
        excerpt,
        content,
        tags: tags.split(',').map((t) => t.trim()).filter(Boolean),
      });
      setArticle(updated);
      toast.success('Draft saved');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to save');
    } finally {
      setIsSaving(false);
    }
  };

  const handlePublish = async () => {
    setIsPublishing(true);
    try {
      await articlesApi.update(params.id, {
        title,
        excerpt,
        content,
        tags: tags.split(',').map((t) => t.trim()).filter(Boolean),
      });
      const published = await articlesApi.publish(params.id);
      toast.success('Published!');
      router.push(`/article/${published.slug || published.id}`);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to publish');
    } finally {
      setIsPublishing(false);
    }
  };

  if (!article) {
    return (
      <ProtectedRoute>
        <div className="flex min-h-screen items-center justify-center">Loading...</div>
      </ProtectedRoute>
    );
  }

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-background">
        <Header />
        <main className="mx-auto max-w-2xl px-4 py-8">
          <div className="mb-6 flex items-center justify-between">
            <span className="rounded-full bg-secondary px-3 py-1 text-xs font-medium uppercase">
              {article.status}
            </span>
            <div className="flex gap-2">
              <Button variant="outline" onClick={handleSave} disabled={isSaving}>
                {isSaving ? 'Saving...' : 'Save'}
              </Button>
              <Button onClick={handlePublish} disabled={isPublishing}>
                {isPublishing ? 'Publishing...' : 'Publish'}
              </Button>
            </div>
          </div>

          <div className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="title">Title</Label>
              <Input id="title" value={title} onChange={(e) => setTitle(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="excerpt">Subtitle / excerpt</Label>
              <Input id="excerpt" value={excerpt} onChange={(e) => setExcerpt(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="tags">Tags (comma separated)</Label>
              <Input id="tags" value={tags} onChange={(e) => setTags(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="content">Content</Label>
              <Textarea
                id="content"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                rows={20}
                className="font-serif text-base leading-relaxed"
              />
            </div>
          </div>
        </main>
      </div>
    </ProtectedRoute>
  );
}
