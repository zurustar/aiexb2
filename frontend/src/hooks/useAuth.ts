import { useCallback, useEffect, useMemo, useState } from "react";

import { ApiClient } from "@/lib/api-client";
import { AuthLoginParams, AuthManager, Session, createMemoryStorage } from "@/lib/auth";
import { Role, User } from "@/types/models";

const createDefaultAuthManager = (): AuthManager => {
  const storage = typeof window !== "undefined" ? window.localStorage : createMemoryStorage();
  let manager: AuthManager;
  const apiClient = new ApiClient({
    getAuthToken: () => manager?.loadSession()?.accessToken ?? null,
  });

  manager = new AuthManager(apiClient, storage as Storage);
  return manager;
};

const defaultAuthManager = createDefaultAuthManager();

export type UseAuthState = {
  session: Session | null;
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
};

export type UseAuthResult = UseAuthState & {
  login: (params: AuthLoginParams) => Promise<void>;
  logout: () => Promise<void>;
  hasRole: (role: Role | Role[]) => boolean;
  refresh: () => void;
  refreshToken: () => Promise<void>;
};

export const useAuth = (manager: AuthManager = defaultAuthManager): UseAuthResult => {
  const [session, setSession] = useState<Session | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadSession = useCallback(() => {
    const current = manager.loadSession();
    setSession(current);
    setIsLoading(false);
  }, [manager]);

  // Cross-tab session synchronization
  useEffect(() => {
    if (typeof window === "undefined") return;

    const handleStorageChange = (e: StorageEvent) => {
      // Only react to session key changes
      if (e.key === "esms.session" || e.key === null) {
        loadSession();
      }
    };

    window.addEventListener("storage", handleStorageChange);
    return () => window.removeEventListener("storage", handleStorageChange);
  }, [loadSession]);

  // Auto token refresh
  useEffect(() => {
    if (!session) return;

    const checkAndRefresh = async () => {
      if (manager.shouldRefreshToken()) {
        try {
          const refreshed = await manager.refreshToken();
          if (refreshed) {
            setSession(refreshed);
          } else {
            // Refresh failed, clear session
            setSession(null);
          }
        } catch (error) {
          console.error("Token refresh failed:", error);
          setSession(null);
        }
      }
    };

    // Check every 60 seconds
    const intervalId = setInterval(checkAndRefresh, 60000);

    // Also check immediately
    checkAndRefresh();

    return () => clearInterval(intervalId);
  }, [session, manager]);

  useEffect(() => {
    loadSession();
  }, [loadSession]);

  const login = useCallback(
    async (params: AuthLoginParams) => {
      setIsLoading(true);
      setError(null);
      try {
        const nextSession = await manager.login(params);
        setSession(nextSession);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Login failed");
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [manager]
  );

  const logout = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      await manager.logout();
      setSession(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Logout failed");
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [manager]);

  const hasRole = useCallback(
    (role: Role | Role[]) => {
      return manager.hasRole(role);
    },
    [manager]
  );

  const refreshToken = useCallback(async () => {
    try {
      const refreshed = await manager.refreshToken();
      if (refreshed) {
        setSession(refreshed);
      }
    } catch (error) {
      setError(error instanceof Error ? error.message : "Token refresh failed");
      throw error;
    }
  }, [manager]);

  const state = useMemo<UseAuthState>(
    () => ({
      session,
      user: session?.user ?? null,
      isAuthenticated: !!session && manager.isAuthenticated(),
      isLoading,
      error,
    }),
    [session, isLoading, error, manager]
  );

  return {
    ...state,
    login,
    logout,
    hasRole,
    refresh: loadSession,
    refreshToken,
  };
};

