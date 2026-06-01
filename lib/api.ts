const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

interface RequestOptions {
  method?: string;
  body?: object;
  headers?: Record<string, string>;
}

async function request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method: options.method || 'GET',
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || 'Request failed');
  }

  return response.json();
}

// Types
export interface User {
  id: number;
  email?: string;
  phone?: string;
  name: string;
  bio: string;
  avatar: string;
  locale: string;
  followers: number;
  created_at: string;
  updated_at: string;
}

export interface Article {
  id: number;
  title: string;
  excerpt: string;
  content: string;
  image: string;
  author_id: number;
  author?: User;
  tags: string[];
  read_time: string;
  is_member_only: boolean;
  likes_count: number;
  comments_count: number;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: number;
  article_id: number;
  user_id: number;
  user?: User;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface AuthResponse {
  token: string;
  user: User;
  is_new_user?: boolean;
}

export interface SendOtpResponse {
  message: string;
  expires_in: number;
}

// Auth API
export const authApi = {
  sendOtp: (data: { phone: string; name?: string }) =>
    request<SendOtpResponse>('/api/auth/send-otp', { method: 'POST', body: data }),

  verifyOtp: (data: { phone: string; code: string; name?: string }) =>
    request<AuthResponse>('/api/auth/verify-otp', { method: 'POST', body: data }),

  register: (data: { email: string; password: string; name: string }) =>
    request<AuthResponse>('/api/auth/register', { method: 'POST', body: data }),

  login: (data: { email: string; password: string }) =>
    request<AuthResponse>('/api/auth/login', { method: 'POST', body: data }),

  getCurrentUser: () => request<User>('/api/users/me'),

  updateUser: (data: { name?: string; bio?: string; avatar?: string; locale?: string }) =>
    request<User>('/api/users/me', { method: 'PUT', body: data }),
};

// User API
export const userApi = {
  getUser: (id: number | string) => request<User>(`/api/users/${id}`),

  getUserArticles: (id: number | string, page = 1, perPage = 10) =>
    request<PaginatedResponse<Article>>(`/api/users/${id}/articles?page=${page}&per_page=${perPage}`),

  followUser: (id: number | string) =>
    request<{ following: boolean; message: string }>(`/api/users/${id}/follow`, { method: 'POST' }),
};

// Articles API
export const articlesApi = {
  getArticles: (params?: { page?: number; per_page?: number; tag?: string; author_id?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', String(params.page));
    if (params?.per_page) searchParams.set('per_page', String(params.per_page));
    if (params?.tag) searchParams.set('tag', params.tag);
    if (params?.author_id) searchParams.set('author_id', params.author_id);
    return request<PaginatedResponse<Article>>(`/api/articles?${searchParams}`);
  },

  getArticle: (id: number | string) => request<Article>(`/api/articles/${id}`),

  getTrendingArticles: (limit = 6) => request<Article[]>(`/api/articles/trending?limit=${limit}`),

  searchArticles: (query: string, page = 1, perPage = 10) =>
    request<PaginatedResponse<Article>>(
      `/api/articles/search?q=${encodeURIComponent(query)}&page=${page}&per_page=${perPage}`
    ),

  createArticle: (data: {
    title: string;
    excerpt: string;
    content: string;
    image?: string;
    tags?: string[];
    is_member_only?: boolean;
  }) => request<Article>('/api/articles', { method: 'POST', body: data }),

  updateArticle: (
    id: number | string,
    data: {
      title?: string;
      excerpt?: string;
      content?: string;
      image?: string;
      tags?: string[];
      is_member_only?: boolean;
    }
  ) => request<Article>(`/api/articles/${id}`, { method: 'PUT', body: data }),

  deleteArticle: (id: number | string) =>
    request<{ message: string }>(`/api/articles/${id}`, { method: 'DELETE' }),

  likeArticle: (id: number | string) =>
    request<{ liked: boolean; message: string }>(`/api/articles/${id}/like`, { method: 'POST' }),

  bookmarkArticle: (id: number | string) =>
    request<{ bookmarked: boolean; message: string }>(`/api/articles/${id}/bookmark`, {
      method: 'POST',
    }),
};

// Comments API
export const commentsApi = {
  getComments: (articleId: number | string, page = 1, perPage = 20) =>
    request<PaginatedResponse<Comment>>(
      `/api/articles/${articleId}/comments?page=${page}&per_page=${perPage}`
    ),

  createComment: (articleId: number | string, content: string) =>
    request<Comment>(`/api/articles/${articleId}/comments`, { method: 'POST', body: { content } }),

  deleteComment: (articleId: number | string, commentId: number | string) =>
    request<{ message: string }>(`/api/articles/${articleId}/comments/${commentId}`, {
      method: 'DELETE',
    }),
};

// Bookmarks API
export const bookmarksApi = {
  getBookmarks: (page = 1, perPage = 10) =>
    request<PaginatedResponse<Article>>(`/api/bookmarks?page=${page}&per_page=${perPage}`),
};

// Token management
export const tokenManager = {
  setToken: (token: string) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token);
    }
  },

  getToken: () => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('auth_token');
    }
    return null;
  },

  removeToken: () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  },
};
