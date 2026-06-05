import { cookies } from 'next/headers';
import type { Article, PaginatedResponse, User } from '@/lib/api';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

async function serverFetch<T>(path: string, init?: RequestInit): Promise<T | null> {
  const cookieStore = await cookies();
  const token = cookieStore.get('auth_token')?.value;

  const headers: Record<string, string> = {
    Accept: 'application/json',
    ...(init?.headers as Record<string, string>),
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  try {
    const response = await fetch(`${API_BASE}/api${path}`, {
      ...init,
      headers,
      cache: 'no-store',
    });
    if (!response.ok) return null;
    return response.json() as Promise<T>;
  } catch {
    return null;
  }
}

export async function fetchArticles(page = 1, perPage = 10) {
  return serverFetch<PaginatedResponse<Article>>(`/articles?page=${page}&per_page=${perPage}`);
}

export async function fetchTrending(limit = 6) {
  return serverFetch<Article[]>(`/articles/trending?limit=${limit}`);
}

export async function fetchArticleById(id: string) {
  if (/^\d+$/.test(id)) {
    return serverFetch<Article>(`/articles/${id}`);
  }
  return serverFetch<Article>(`/articles/slug/${encodeURIComponent(id)}`);
}

export async function fetchSearchArticles(query: string, page = 1) {
  return serverFetch<PaginatedResponse<Article>>(
    `/articles/search?q=${encodeURIComponent(query)}&page=${page}&per_page=10`
  );
}

export async function fetchUserArticles(userId: number, page = 1) {
  return serverFetch<PaginatedResponse<Article>>(
    `/users/${userId}/articles?page=${page}&per_page=10`
  );
}

export async function fetchCurrentUser() {
  return serverFetch<User>('/users/me');
}
