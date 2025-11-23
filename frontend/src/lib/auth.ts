import { ApiClient } from "./api-client";
import { Role, User } from "@/types/models";

export type Session = {
  accessToken: string;
  refreshToken?: string;
  user: User;
  expiresAt?: number;
};

export type AuthLoginParams = {
  username: string;
  password: string;
};

export interface StorageLike {
  getItem(key: string): string | null;
  setItem(key: string, value: string): void;
  removeItem(key: string): void;
}

class MemoryStorage implements StorageLike {
  private store = new Map<string, string>();

  getItem(key: string): string | null {
    return this.store.get(key) ?? null;
  }

  setItem(key: string, value: string): void {
    this.store.set(key, value);
  }

  removeItem(key: string): void {
    this.store.delete(key);
  }
}

const DEFAULT_SESSION_KEY = "esms.session";

export class AuthManager {
  private readonly apiClient: ApiClient;
  private readonly storage: StorageLike;
  private readonly sessionKey: string;

  constructor(apiClient: ApiClient, storage: StorageLike = new MemoryStorage(), sessionKey: string = DEFAULT_SESSION_KEY) {
    this.apiClient = apiClient;
    this.storage = storage;
    this.sessionKey = sessionKey;
  }

  loadSession(): Session | null {
    const raw = this.storage.getItem(this.sessionKey);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as Session;
    } catch (error) {
      this.storage.removeItem(this.sessionKey);
      return null;
    }
  }

  saveSession(session: Session): void {
    this.storage.setItem(this.sessionKey, JSON.stringify(session));
  }

  clearSession(): void {
    this.storage.removeItem(this.sessionKey);
  }

  async login(params: AuthLoginParams): Promise<Session> {
    const result = await this.apiClient.post<Session>("/api/v1/auth/login", params);
    this.saveSession(result.data);
    return result.data;
  }

  async logout(): Promise<void> {
    try {
      await this.apiClient.post<null>("/api/v1/auth/logout");
    } finally {
      this.clearSession();
    }
  }

  isAuthenticated(): boolean {
    const session = this.loadSession();
    if (!session) return false;
    if (session.expiresAt && Date.now() > session.expiresAt) {
      this.clearSession();
      return false;
    }
    return true;
  }

  currentUser(): User | null {
    return this.loadSession()?.user ?? null;
  }

  hasRole(required: Role | Role[]): boolean {
    const user = this.currentUser();
    if (!user) return false;
    if (Array.isArray(required)) {
      return required.includes(user.role);
    }
    return user.role === required;
  }
}

export const createMemoryStorage = () => new MemoryStorage();
