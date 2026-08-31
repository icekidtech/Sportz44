import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import { api, setUnauthorizedHandler } from '../api/client';
import type { User } from '../api/types';

interface AuthState {
  user: User | null;
  loading: boolean;
  login: (identifier: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthState>(null!);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchMe = useCallback(async () => {
    try {
      const res = await api.get<{ user: User }>('/api/auth/me');
      setUser(res.user);
    } catch {
      setUser(null);
    }
  }, []);

  useEffect(() => {
    setUnauthorizedHandler(() => setUser(null));
  }, []);

  useEffect(() => {
    (async () => {
      await fetchMe();
      setLoading(false);
    })();
  }, [fetchMe]);

  const login = async (identifier: string, password: string) => {
    const res = await api.post<{ user: User }>('/api/auth/login', { identifier, password });
    setUser(res.user);
  };

  const register = async (username: string, email: string, password: string) => {
    const res = await api.post<{ user: User }>('/api/auth/register', { username, email, password });
    setUser(res.user);
  };

  const logout = async () => {
    try {
      await api.post('/api/auth/logout');
    } catch {}
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
