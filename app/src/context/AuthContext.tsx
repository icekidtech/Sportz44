import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import * as SecureStore from 'expo-secure-store';
import { api, setAuthTokenGetter } from '../api/client';
import type { User } from '../api/types';

const TOKEN_KEY = 'sportz44_token';

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<void>;
}

const AuthContext = createContext<AuthState>(null!);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setAuthTokenGetter(() => token);
  }, [token]);

  const fetchMe = useCallback(async (t: string) => {
    try {
      const res = await api.get<{ user: User }>('/api/auth/me', {
        headers: { Authorization: `Bearer ${t}` },
      });
      setUser(res.user);
    } catch {
      setUser(null);
    }
  }, []);

  useEffect(() => {
    (async () => {
      try {
        const stored = await SecureStore.getItemAsync(TOKEN_KEY);
        if (stored) {
          setToken(stored);
          await fetchMe(stored);
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [fetchMe]);

  const login = async (email: string, password: string) => {
    const res = await api.post<{ access_token: string; user: User }>('/api/auth/login', { email, password }, { skipAuth: true });
    await SecureStore.setItemAsync(TOKEN_KEY, res.access_token);
    setToken(res.access_token);
    setUser(res.user);
  };

  const register = async (username: string, email: string, password: string) => {
    const res = await api.post<{ access_token: string; user: User }>('/api/auth/register', { username, email, password }, { skipAuth: true });
    await SecureStore.setItemAsync(TOKEN_KEY, res.access_token);
    setToken(res.access_token);
    setUser(res.user);
  };

  const logout = async () => {
    try { await api.post('/api/auth/logout'); } catch {}
    await SecureStore.deleteItemAsync(TOKEN_KEY);
    setToken(null);
    setUser(null);
  };

  const refresh = async () => {
    const res = await api.post<{ access_token: string }>('/api/auth/refresh');
    await SecureStore.setItemAsync(TOKEN_KEY, res.access_token);
    setToken(res.access_token);
  };

  return (
    <AuthContext.Provider value={{ user, token, loading, login, register, logout, refresh }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
