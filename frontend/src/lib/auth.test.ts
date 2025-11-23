import { ApiClient } from "./api-client";
import { AuthManager, createMemoryStorage, Session } from "./auth";
import { User } from "@/types/models";

describe("AuthManager", () => {
  const mockFetch = jest.fn<Promise<Response>, [RequestInfo | URL, RequestInit | undefined]>();
  const apiClient = new ApiClient({ fetchImpl: mockFetch });
  const storage = createMemoryStorage();
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
  const session: Session = { accessToken: "token", user };
  const manager = new AuthManager(apiClient, storage);

  beforeEach(() => {
    storage.removeItem("esms.session");
    mockFetch.mockReset();
  });

  it("saves session when login succeeds", async () => {
    mockFetch.mockResolvedValue(
      new Response(JSON.stringify({ data: session }), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    await manager.login({ username: "user", password: "pass" });

    expect(manager.isAuthenticated()).toBe(true);
    expect(manager.currentUser()).toEqual(user);
  });

  it("clears session on logout regardless of API result", async () => {
    storage.setItem("esms.session", JSON.stringify(session));
    mockFetch.mockResolvedValue(new Response(null, { status: 500 }));

    await manager.logout();

    expect(manager.isAuthenticated()).toBe(false);
    expect(manager.currentUser()).toBeNull();
  });

  it("checks required role", () => {
    manager.saveSession(session);

    expect(manager.hasRole("ADMIN")).toBe(true);
    expect(manager.hasRole(["ADMIN", "GENERAL"])).toBe(true);
    expect(manager.hasRole("GENERAL")).toBe(false);
  });

  it("invalidates expired session", () => {
    manager.saveSession({ ...session, expiresAt: Date.now() - 1000 });
    expect(manager.isAuthenticated()).toBe(false);
  });
});
