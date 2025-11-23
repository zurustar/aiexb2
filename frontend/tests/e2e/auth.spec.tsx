import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { CallbackPage } from "@/app/(auth)/callback/page";
import { LoginPage } from "@/app/(auth)/login/page";
import { UseAuthResult } from "@/hooks/useAuth";

const useSearchParamsMock = jest.fn();

jest.mock("next/navigation", () => ({
  __esModule: true,
  useSearchParams: () => useSearchParamsMock(),
}));

type PartialAuth = Partial<Pick<UseAuthResult, "login" | "logout" | "hasRole" | "refresh">> & {
  user?: UseAuthResult["user"];
  session?: UseAuthResult["session"];
  isAuthenticated?: boolean;
  isLoading?: boolean;
  error?: string | null;
};

const createAuthState = (overrides: PartialAuth = {}): UseAuthResult => ({
  user: null,
  session: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
  login: jest.fn(),
  logout: jest.fn(),
  hasRole: jest.fn().mockReturnValue(false),
  refresh: jest.fn(),
  ...overrides,
});

describe("E2E: 認証フロー", () => {
  beforeEach(() => {
    useSearchParamsMock.mockReturnValue({ get: () => null });
  });

  it("handles login form submission", async () => {
    const login = jest.fn().mockResolvedValue(undefined);
    const mockAuth = createAuthState({ login });
    const useAuthHook = () => mockAuth;

    render(<LoginPage useAuthHook={useAuthHook} />);

    fireEvent.change(screen.getByLabelText("ユーザー名"), { target: { value: "demo" } });
    fireEvent.change(screen.getByLabelText("パスワード"), { target: { value: "secret" } });
    fireEvent.click(screen.getByRole("button", { name: "サインイン" }));

    await waitFor(() => expect(login).toHaveBeenCalledWith({ username: "demo", password: "secret" }));
    expect(screen.getByText(/ログインに成功/)).toBeInTheDocument();
  });

  it("shows callback status when authenticated", async () => {
    useSearchParamsMock.mockReturnValue({
      get: (key: string) => (key === "code" ? "auth-code" : "state-token"),
    });

    const mockAuth = createAuthState({ isAuthenticated: true, refresh: jest.fn(), isLoading: false });
    const useAuthHook = () => mockAuth;

    render(<CallbackPage useAuthHook={useAuthHook} />);

    expect(useSearchParamsMock).toHaveBeenCalled();
    await waitFor(() => expect(screen.getByTestId("status-message")).toHaveTextContent(/ログインが完了/));
    expect(screen.getByText("auth-code")).toBeInTheDocument();
    expect(screen.getByText("state-token")).toBeInTheDocument();
  });
});
