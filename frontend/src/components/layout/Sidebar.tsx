import React, { ReactNode } from "react";

import { useAuth, UseAuthResult } from "@/hooks/useAuth";
import { Role } from "@/types/models";

export type SidebarItem = {
  label: string;
  href: string;
  icon?: ReactNode;
  badge?: string | number;
  roles?: Role[];
  active?: boolean;
};

export type SidebarProps = {
  items: SidebarItem[];
  collapsed?: boolean;
  onItemSelect?: (item: SidebarItem) => void;
  auth?: Pick<UseAuthResult, "hasRole" | "isAuthenticated">;
};

export const Sidebar: React.FC<SidebarProps> = ({ items, collapsed = false, onItemSelect, auth }) => {
  const authState = auth ?? useAuth();

  const isAllowed = (item: SidebarItem): boolean => {
    if (!item.roles || item.roles.length === 0) return true;
    if (!authState.isAuthenticated) return false;
    return item.roles.some((role) => authState.hasRole(role));
  };

  const visibleItems = items.filter(isAllowed);

  return (
    <aside
      className={`flex h-full flex-col gap-2 border-r border-gray-200 bg-gray-50 px-3 py-4 ${collapsed ? "w-16" : "w-60"}`}
      aria-label="サイドバー"
    >
      {visibleItems.map((item) => (
        <a
          key={item.href}
          href={item.href}
          onClick={() => onItemSelect?.(item)}
          className={`flex items-center justify-between rounded-md px-3 py-2 text-sm font-medium transition hover:bg-white hover:shadow ${
            item.active ? "bg-white text-blue-700 shadow" : "text-gray-800"
          }`}
        >
          <span className="flex items-center gap-2">
            {item.icon}
            <span className={collapsed ? "sr-only" : ""}>{item.label}</span>
          </span>
          {!collapsed && item.badge !== undefined && (
            <span className="rounded-full bg-blue-100 px-2 py-0.5 text-xs font-semibold text-blue-700" aria-label="badge">
              {item.badge}
            </span>
          )}
        </a>
      ))}
      {visibleItems.length === 0 && (
        <p className="text-sm text-gray-500">表示できるメニューがありません。</p>
      )}
    </aside>
  );
};

export default Sidebar;
