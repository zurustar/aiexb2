import { renderHook, act, waitFor } from "@testing-library/react";

import { ApiClient } from "@/lib/api-client";
import { AuthLoginParams, AuthManager, Session, createMemoryStorage } from "@/lib/auth";
import { useAuth } from "./useAuth";
import { User } from "@/types/models";

const mockFetch = jest.fn<Promise<Response>, [RequestInfo | URL, RequestInit | undefined]>();
const apiClient = new ApiClient({ fetchImpl: mockFetch });
const storage = createMemoryStorage();
const authManager = new AuthManager(apiClient, storage);

const user: User = {
  id: "1",
  sub: "sub",
  email: "user@example.com",
  name: "User",
  role: "ADMIN",
  penaltyScore: 0,
  isActive: true,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
};

const session: Session = {
  accessToken: "token",
  user,
  expiresAt: Date.now() + 1000 * 60,
};

describe("useAuth", () => {
  beforeEach(() => {
    storage.removeItem("esms.session");
    mockFetch.mockReset();
  });

  it("loads session on mount", async () => {
    storage.setItem("esms.session", JSON.stringify(session));

    const { result } = renderHook(() => useAuth(authManager));
    await waitFor(() => expect(result.current.isAuthenticated).toBe(true));

    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.user).toEqual(user);
  });

  it("logs in and updates state", async () => {
    mockFetch.mockResolvedValue(
      new Response(JSON.stringify({ data: session }), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    const { result } = renderHook(() => useAuth(authManager));
    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await act(async () => {
      await result.current.login({ username: "user", password: "pass" } satisfies AuthLoginParams);
    });

    expect(result.current.user).toEqual(user);
    expect(result.current.isAuthenticated).toBe(true);
    expect(storage.getItem("esms.session")).not.toBeNull();
  });

  it("logs out and clears session", async () => {
    storage.setItem("esms.session", JSON.stringify(session));
    mockFetch.mockResolvedValue(new Response(null, { status: 200 }));

    const { result } = renderHook(() => useAuth(authManager));
    await waitFor(() => expect(result.current.isAuthenticated).toBe(true));

    await act(async () => {
      await result.current.logout();
    });

    expect(result.current.user).toBeNull();
    expect(result.current.isAuthenticated).toBe(false);
    expect(storage.getItem("esms.session")).toBeNull();
  });

  it("exposes role helper", async () => {
    storage.setItem("esms.session", JSON.stringify(session));
    const { result } = renderHook(() => useAuth(authManager));
    await waitFor(() => expect(result.current.isAuthenticated).toBe(true));

    expect(result.current.hasRole("ADMIN")).toBe(true);
    expect(result.current.hasRole(["GENERAL", "ADMIN"])).toBe(true);
    expect(result.current.hasRole("GENERAL")).toBe(false);
  });
});

