import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";

import { Header } from "./Header";
import { Role } from "@/types/models";

describe("Header", () => {
  const authState = {
    user: { name: "管理者", role: "ADMIN" as Role },
    isAuthenticated: true,
    isLoading: false,
    logout: jest.fn(),
  };

  it("renders navigation items and user info", () => {
    render(
      <Header
        title="予約管理"
        navItems={[
          { label: "ダッシュボード", href: "/dashboard", active: true },
          { label: "リソース", href: "/resources" },
        ]}
        auth={authState}
      />
    );

    expect(screen.getByText("予約管理")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "ダッシュボード" })).toHaveAttribute("href", "/dashboard");
    expect(screen.getByText("管理者")).toBeInTheDocument();
    expect(screen.getByText(/Role: ADMIN/)).toBeInTheDocument();
  });

  it("triggers logout when authenticated user clicks the button", () => {
    const logout = jest.fn();

    render(<Header auth={{ ...authState, logout }} />);

    fireEvent.click(screen.getByRole("button", { name: "ログアウト" }));
    expect(logout).toHaveBeenCalledTimes(1);
  });

  it("calls toggle handler when menu button is clicked", () => {
    const handleToggle = jest.fn();

    render(<Header navItems={[]} auth={authState} onToggleSidebar={handleToggle} />);

    fireEvent.click(screen.getByLabelText("メニューを開く"));
    expect(handleToggle).toHaveBeenCalledTimes(1);
  });
});
