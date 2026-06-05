import { apiRequest } from '@/lib/auth/api-client';

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

// Auth API — public endpoints do not attach/clear session on 401
export const authApi = {
  sendOtp: (data: { phone: string; name?: string }) =>
    apiRequest<SendOtpResponse>('/auth/send-otp', { method: 'POST', body: data, auth: false }),

  verifyOtp: (data: { phone: string; code: string; name?: string }) =>
    apiRequest<AuthResponse>('/auth/verify-otp', { method: 'POST', body: data, auth: false }),

  register: (data: { email: string; password: string; name: string }) =>
    apiRequest<AuthResponse>('/auth/register', { method: 'POST', body: data, auth: false }),

  login: (data: { email: string; password: string }) =>
    apiRequest<AuthResponse>('/auth/login', { method: 'POST', body: data, auth: false }),

  getCurrentUser: () => apiRequest<User>('/users/me'),

  updateUser: (data: { name?: string; bio?: string; avatar?: string; locale?: string }) =>
    apiRequest<User>('/users/me', { method: 'PUT', body: data }),
};

// User API
export const userApi = {
  getUser: (id: number | string) => apiRequest<User>(`/users/${id}`),

  getUserArticles: (id: number | string, page = 1, perPage = 10) =>
    apiRequest<PaginatedResponse<Article>>(
      `/users/${id}/articles?page=${page}&per_page=${perPage}`
    ),

  followUser: (id: number | string) =>
    apiRequest<{ following: boolean; message: string }>(`/users/${id}/follow`, {
      method: 'POST',
    }),
};

// Articles API
export const articlesApi = {
  getArticles: (params?: { page?: number; per_page?: number; tag?: string; author_id?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', String(params.page));
    if (params?.per_page) searchParams.set('per_page', String(params.per_page));
    if (params?.tag) searchParams.set('tag', params.tag);
    if (params?.author_id) searchParams.set('author_id', params.author_id);
    return apiRequest<PaginatedResponse<Article>>(`/articles?${searchParams}`);
  },

  getArticle: (id: number | string) => apiRequest<Article>(`/articles/${id}`),

  getTrendingArticles: (limit = 6) =>
    apiRequest<Article[]>(`/articles/trending?limit=${limit}`),

  searchArticles: (query: string, page = 1, perPage = 10) =>
    apiRequest<PaginatedResponse<Article>>(
      `/articles/search?q=${encodeURIComponent(query)}&page=${page}&per_page=${perPage}`
    ),

  createArticle: (data: {
    title: string;
    excerpt: string;
    content: string;
    image?: string;
    tags?: string[];
    is_member_only?: boolean;
  }) => apiRequest<Article>('/articles', { method: 'POST', body: data }),

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
  ) => apiRequest<Article>(`/articles/${id}`, { method: 'PUT', body: data }),

  deleteArticle: (id: number | string) =>
    apiRequest<{ message: string }>(`/articles/${id}`, { method: 'DELETE' }),

  likeArticle: (id: number | string) =>
    apiRequest<{ liked: boolean; message: string }>(`/articles/${id}/like`, { method: 'POST' }),

  bookmarkArticle: (id: number | string) =>
    apiRequest<{ bookmarked: boolean; message: string }>(`/articles/${id}/bookmark`, {
      method: 'POST',
    }),
};

// Comments API
export const commentsApi = {
  getComments: (articleId: number | string, page = 1, perPage = 20) =>
    apiRequest<PaginatedResponse<Comment>>(
      `/articles/${articleId}/comments?page=${page}&per_page=${perPage}`
    ),

  createComment: (articleId: number | string, content: string) =>
    apiRequest<Comment>(`/articles/${articleId}/comments`, {
      method: 'POST',
      body: { content },
    }),

  deleteComment: (articleId: number | string, commentId: number | string) =>
    apiRequest<{ message: string }>(`/articles/${articleId}/comments/${commentId}`, {
      method: 'DELETE',
    }),
};

// Bookmarks API
export const bookmarksApi = {
  getBookmarks: (page = 1, perPage = 10) =>
    apiRequest<PaginatedResponse<Article>>(`/bookmarks?page=${page}&per_page=${perPage}`),
};

import { clearAccessToken, getAccessToken, setAccessToken } from '@/lib/auth/session';

/** @deprecated Use setAccessToken/clearAccessToken/getAccessToken from @/lib/auth/session */
export const tokenManager = {
  setToken: setAccessToken,
  getToken: getAccessToken,
  removeToken: clearAccessToken,
};
