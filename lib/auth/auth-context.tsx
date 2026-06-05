'use client';

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react';
import { authApi, type User } from '@/lib/api';
import {
  AUTH_SESSION_EXPIRED_EVENT,
  clearAccessToken,
  getAccessToken,
  getRefreshToken,
  setSessionTokens,
} from '@/lib/auth/session';
import { isApiError } from '@/lib/auth/api-client';
import { syncLocaleFromProfile } from '@/lib/locale';

interface AuthContextValue {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (token: string, refreshToken: string, user: User) => void;
  logout: () => Promise<void>;
  refreshUser: () => Promise<User | null>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const refreshUser = useCallback(async (): Promise<User | null> => {
    const token = getAccessToken();
    if (!token) {
      setUser(null);
      return null;
    }

    try {
      const currentUser = await authApi.getCurrentUser();
      setUser(currentUser);
      return currentUser;
    } catch (error) {
      if (isApiError(error) && error.status === 401) {
        const refreshToken = getRefreshToken();
        if (refreshToken) {
          try {
            const response = await authApi.refresh(refreshToken);
            setSessionTokens(response.token, response.refresh_token);
            setUser(response.user);
            return response.user;
          } catch {
            clearAccessToken();
            setUser(null);
            return null;
          }
        }
        clearAccessToken();
        setUser(null);
        return null;
      }
      setUser(null);
      return null;
    }
  }, []);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      await refreshUser();
      if (!cancelled) setIsLoading(false);
    })();

    return () => {
      cancelled = true;
    };
  }, [refreshUser]);

  useEffect(() => {
    const onSessionExpired = () => {
      clearAccessToken();
      setUser(null);
    };

    window.addEventListener(AUTH_SESSION_EXPIRED_EVENT, onSessionExpired);
    return () => window.removeEventListener(AUTH_SESSION_EXPIRED_EVENT, onSessionExpired);
  }, []);

  useEffect(() => {
    if (user?.locale) {
      syncLocaleFromProfile(user.locale);
    }
  }, [user?.locale]);

  const login = useCallback((token: string, refreshToken: string, loggedInUser: User) => {
    setSessionTokens(token, refreshToken);
    setUser(loggedInUser);
    syncLocaleFromProfile(loggedInUser.locale);
  }, []);

  const logout = useCallback(async () => {
    const refreshToken = getRefreshToken();
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken);
      } catch {
        // Clear local session even if server logout fails
      }
    }
    clearAccessToken();
    setUser(null);
  }, []);

  const value = useMemo(
    () => ({
      user,
      isLoading,
      isAuthenticated: !!user,
      login,
      logout,
      refreshUser,
    }),
    [user, isLoading, login, logout, refreshUser]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
