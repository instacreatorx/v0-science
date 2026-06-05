import { apiRequest } from '@/lib/auth/api-client';

export interface User {
  id: number;
  email?: string;
  phone?: string;
  name: string;
  bio: string;
  avatar: string;
  locale: string;
  role?: string;
  followers: number;
  created_at: string;
  updated_at: string;
  is_following?: boolean;
}

export interface Team {
  id: number;
  name: string;
  slug: string;
  bio?: string;
  avatar?: string;
  owner_id: number;
  owner?: User;
  verified_at?: string | null;
  is_verified?: boolean;
  created_at: string;
  updated_at: string;
}

export interface Article {
  id: number;
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  image: string;
  author_id: number;
  author?: User;
  team_id?: number | null;
  team?: Team | null;
  tags: string[];
  read_time: string;
  status: 'draft' | 'published' | 'archived' | string;
  is_member_only: boolean;
  likes_count: number;
  comments_count: number;
  published_at?: string | null;
  created_at: string;
  updated_at: string;
  liked_by_me?: boolean;
  bookmarked_by_me?: boolean;
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

export interface TeamMember {
  id: number;
  team_id: number;
  user_id: number;
  user?: User;
  role: string;
  created_at: string;
}

export interface TeamVerificationRequest {
  id: number;
  team_id: number;
  team?: Team;
  submitted_by: number;
  proof_text: string;
  proof_url?: string;
  status: string;
  rejection_reason?: string;
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
  refresh_token: string;
  user: User;
  is_new_user?: boolean;
}

export interface SendOtpResponse {
  message: string;
  expires_in: number;
}

export const authApi = {
  sendOtp: (data: { phone: string; name?: string }) =>
    apiRequest<SendOtpResponse>('/auth/send-otp', { method: 'POST', body: data, auth: false }),

  verifyOtp: (data: { phone: string; code: string; name?: string }) =>
    apiRequest<AuthResponse>('/auth/verify-otp', { method: 'POST', body: data, auth: false }),

  register: (data: { email: string; password: string; name: string }) =>
    apiRequest<AuthResponse>('/auth/register', { method: 'POST', body: data, auth: false }),

  login: (data: { email: string; password: string }) =>
    apiRequest<AuthResponse>('/auth/login', { method: 'POST', body: data, auth: false }),

  refresh: (refreshToken: string) =>
    apiRequest<AuthResponse>('/auth/refresh', {
      method: 'POST',
      body: { refresh_token: refreshToken },
      auth: false,
    }),

  logout: (refreshToken: string) =>
    apiRequest<{ message: string }>('/auth/logout', {
      method: 'POST',
      body: { refresh_token: refreshToken },
    }),

  getCurrentUser: () => apiRequest<User>('/users/me'),

  updateUser: (data: { name?: string; bio?: string; avatar?: string; locale?: string }) =>
    apiRequest<User>('/users/me', { method: 'PUT', body: data }),
};

export const userApi = {
  getUser: (id: number | string) => apiRequest<User>(`/users/${id}`),

  getUserArticles: (id: number | string, page = 1, perPage = 10) =>
    apiRequest<PaginatedResponse<Article>>(`/users/${id}/articles?page=${page}&per_page=${perPage}`),

  followUser: (id: number | string) =>
    apiRequest<{ following: boolean; message: string }>(`/users/${id}/follow`, { method: 'POST' }),
};

export const articlesApi = {
  list: (params?: { page?: number; per_page?: number; tag?: string; author_id?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', String(params.page));
    if (params?.per_page) searchParams.set('per_page', String(params.per_page));
    if (params?.tag) searchParams.set('tag', params.tag);
    if (params?.author_id) searchParams.set('author_id', params.author_id);
    return apiRequest<PaginatedResponse<Article>>(`/articles?${searchParams}`);
  },

  myArticles: (params?: { page?: number; per_page?: number; status?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', String(params.page));
    if (params?.per_page) searchParams.set('per_page', String(params.per_page));
    if (params?.status) searchParams.set('status', params.status);
    return apiRequest<PaginatedResponse<Article>>(`/articles/me?${searchParams}`);
  },

  feed: (page = 1, perPage = 10) =>
    apiRequest<PaginatedResponse<Article>>(`/feed?page=${page}&per_page=${perPage}`),

  getById: (id: number | string) => apiRequest<Article>(`/articles/${id}`),

  getBySlug: (slug: string) => apiRequest<Article>(`/articles/slug/${encodeURIComponent(slug)}`),

  getTrending: (limit = 6) => apiRequest<Article[]>(`/articles/trending?limit=${limit}`),

  search: (query: string, page = 1, perPage = 10) =>
    apiRequest<PaginatedResponse<Article>>(
      `/articles/search?q=${encodeURIComponent(query)}&page=${page}&per_page=${perPage}`
    ),

  create: (data: {
    title: string;
    excerpt?: string;
    content?: string;
    image?: string;
    tags?: string[];
    team_id?: number;
    is_member_only?: boolean;
  }) => apiRequest<Article>('/articles', { method: 'POST', body: data }),

  update: (
    id: number | string,
    data: {
      title?: string;
      excerpt?: string;
      content?: string;
      image?: string;
      tags?: string[];
      team_id?: number | null;
      is_member_only?: boolean;
    }
  ) => apiRequest<Article>(`/articles/${id}`, { method: 'PUT', body: data }),

  delete: (id: number | string) =>
    apiRequest<{ message: string }>(`/articles/${id}`, { method: 'DELETE' }),

  publish: (id: number | string) =>
    apiRequest<Article>(`/articles/${id}/publish`, { method: 'POST' }),

  unpublish: (id: number | string) =>
    apiRequest<Article>(`/articles/${id}/unpublish`, { method: 'POST' }),

  archive: (id: number | string) =>
    apiRequest<Article>(`/articles/${id}/archive`, { method: 'POST' }),

  like: (id: number | string) =>
    apiRequest<{ liked: boolean; message: string }>(`/articles/${id}/like`, { method: 'POST' }),

  bookmark: (id: number | string) =>
    apiRequest<{ bookmarked: boolean; message: string }>(`/articles/${id}/bookmark`, {
      method: 'POST',
    }),
};

export const commentsApi = {
  list: (articleId: number | string, page = 1, perPage = 20) =>
    apiRequest<PaginatedResponse<Comment>>(
      `/articles/${articleId}/comments?page=${page}&per_page=${perPage}`
    ),

  create: (articleId: number | string, content: string) =>
    apiRequest<Comment>(`/articles/${articleId}/comments`, {
      method: 'POST',
      body: { content },
    }),

  delete: (articleId: number | string, commentId: number | string) =>
    apiRequest<{ message: string }>(`/articles/${articleId}/comments/${commentId}`, {
      method: 'DELETE',
    }),
};

export const bookmarksApi = {
  list: (page = 1, perPage = 10) =>
    apiRequest<PaginatedResponse<Article>>(`/bookmarks?page=${page}&per_page=${perPage}`),
};

export const teamsApi = {
  create: (data: { name: string; slug: string; bio?: string; avatar?: string }) =>
    apiRequest<Team>('/teams', { method: 'POST', body: data }),

  mine: () => apiRequest<Team[]>('/teams/mine'),

  getBySlug: (slug: string) => apiRequest<Team>(`/teams/slug/${encodeURIComponent(slug)}`),

  update: (id: number, data: { name?: string; bio?: string; avatar?: string }) =>
    apiRequest<Team>(`/teams/${id}`, { method: 'PUT', body: data }),

  members: (id: number) => apiRequest<TeamMember[]>(`/teams/${id}/members`),

  addMember: (id: number, data: { user_id: number; role: string }) =>
    apiRequest<TeamMember>(`/teams/${id}/members`, { method: 'POST', body: data }),

  removeMember: (teamId: number, userId: number) =>
    apiRequest<{ message: string }>(`/teams/${teamId}/members/${userId}`, { method: 'DELETE' }),

  submitVerification: (id: number, data: { proof_text: string; proof_url?: string }) =>
    apiRequest<TeamVerificationRequest>(`/teams/${id}/verify-request`, {
      method: 'POST',
      body: data,
    }),
};

export const adminApi = {
  listVerificationRequests: (status = 'pending', page = 1) =>
    apiRequest<PaginatedResponse<TeamVerificationRequest>>(
      `/admin/verification-requests?status=${status}&page=${page}&per_page=20`
    ),

  approveVerification: (id: number) =>
    apiRequest<Team>(`/admin/verification-requests/${id}/approve`, { method: 'POST' }),

  rejectVerification: (id: number, reason: string) =>
    apiRequest<TeamVerificationRequest>(`/admin/verification-requests/${id}/reject`, {
      method: 'POST',
      body: { reason },
    }),
};
