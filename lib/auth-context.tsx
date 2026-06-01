'use client';

import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { authApi, tokenManager, type User } from '@/lib/api';
import { syncLocaleFromProfile } from '@/lib/locale';

interface AuthContextValue {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const refreshUser = useCallback(async () => {
    const token = tokenManager.getToken();
    if (!token) {
      setUser(null);
      return;
    }
    try {
      const currentUser = await authApi.getCurrentUser();
      setUser(currentUser);
    } catch {
      tokenManager.removeToken();
      setUser(null);
    }
  }, []);

  useEffect(() => {
    refreshUser().finally(() => setIsLoading(false));
  }, [refreshUser]);

  useEffect(() => {
    if (user?.locale) {
      syncLocaleFromProfile(user.locale);
    }
  }, [user?.locale]);

  const login = useCallback((token: string, loggedInUser: User) => {
    tokenManager.setToken(token);
    setUser(loggedInUser);
    syncLocaleFromProfile(loggedInUser.locale);
  }, []);

  const logout = useCallback(() => {
    tokenManager.removeToken();
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
