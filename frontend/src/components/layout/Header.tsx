import React, { ReactNode } from "react";

import { Button } from "@/components/ui/Button";
import { useAuth, UseAuthResult } from "@/hooks/useAuth";

export type HeaderNavItem = {
  label: string;
  href: string;
  active?: boolean;
};

export type HeaderProps = {
  title?: string;
  navItems?: HeaderNavItem[];
  actions?: ReactNode;
  onToggleSidebar?: () => void;
  auth?: Pick<UseAuthResult, "user" | "logout" | "isAuthenticated" | "isLoading">;
};

export const Header: React.FC<HeaderProps> = ({
  title = "Enterprise Scheduler",
  navItems = [],
  actions,
  onToggleSidebar,
  auth,
}) => {
  const defaultAuth = useAuth();
  const authState = auth ?? defaultAuth;
  const userName = authState.isAuthenticated ? authState.user?.name ?? "ユーザー" : "ゲスト";
  const roleLabel = authState.user?.role ? `Role: ${authState.user.role}` : "";

  return (
    <header className="flex items-center justify-between border-b border-gray-200 bg-white px-4 py-3 shadow-sm">
      <div className="flex items-center gap-3">
        {onToggleSidebar && (
          <button
            type="button"
            aria-label="メニューを開く"
            onClick={onToggleSidebar}
            className="rounded-md p-2 text-gray-600 transition hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            ☰
          </button>
        )}
        <div>
          <p className="text-lg font-semibold text-gray-900">{title}</p>
          <p className="text-xs text-gray-500">社内リソース管理システム</p>
        </div>
      </div>

      <nav className="hidden items-center gap-4 md:flex" aria-label="メインナビゲーション">
        {navItems.map((item) => (
          <a
            key={item.href}
            href={item.href}
            aria-current={item.active ? "page" : undefined}
            className={`text-sm font-medium transition hover:text-blue-600 ${item.active ? "text-blue-700 underline" : "text-gray-700"
              }`}
          >
            {item.label}
          </a>
        ))}
      </nav>

      <div className="flex items-center gap-3">
        {actions}
        <div className="text-right">
          <p className="text-sm font-medium text-gray-900">{userName}</p>
          {roleLabel && <p className="text-xs text-gray-500">{roleLabel}</p>}
        </div>
        <Button
          variant="secondary"
          size="sm"
          onClick={authState.isAuthenticated ? authState.logout : undefined}
          disabled={authState.isLoading}
        >
          {authState.isAuthenticated ? "ログアウト" : "ログイン"}
        </Button>
      </div>
    </header>
  );
};

export default Header;
