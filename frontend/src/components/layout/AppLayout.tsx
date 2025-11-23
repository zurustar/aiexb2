"use client";

import React, { useMemo, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";

import Footer from "@/components/layout/Footer";
import Header from "@/components/layout/Header";
import Sidebar from "@/components/layout/Sidebar";

const isActivePath = (pathname: string | null, href: string): boolean => {
  if (!pathname) return false;
  if (href === "/") return pathname === "/";
  return pathname === href || pathname.startsWith(`${href}/`);
};

export type AppLayoutProps = {
  children: React.ReactNode;
};

export const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const pathname = usePathname();
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);

  const showNavigation = pathname !== "/login" && pathname !== "/callback";

  const headerNavItems = useMemo(
    () => [
      { label: "トップ", href: "/", active: isActivePath(pathname, "/") },
      { label: "ダッシュボード", href: "/dashboard", active: isActivePath(pathname, "/dashboard") },
      { label: "予定管理", href: "/events", active: isActivePath(pathname, "/events") },
      { label: "リソース管理", href: "/resources", active: isActivePath(pathname, "/resources") },
    ],
    [pathname]
  );

  const sidebarItems = useMemo(
    () => [
      { label: "トップ", href: "/", active: isActivePath(pathname, "/") },
      { label: "ダッシュボード", href: "/dashboard", badge: "新着", active: isActivePath(pathname, "/dashboard") },
      { label: "予定管理", href: "/events", active: isActivePath(pathname, "/events") },
      { label: "承認一覧", href: "/dashboard#approvals", roles: ["MANAGER", "ADMIN"], active: pathname?.includes("approvals") },
      { label: "リソース管理", href: "/resources", active: isActivePath(pathname, "/resources") },
    ],
    [pathname]
  );

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900">
      {showNavigation && (
        <Header
          navItems={headerNavItems}
          onToggleSidebar={() => setIsSidebarCollapsed((prev) => !prev)}
          title="Enterprise Scheduler"
        />
      )}

      <div className="mx-auto flex min-h-[calc(100vh-120px)] max-w-7xl flex-col px-3 md:px-6">
        <div className={`flex flex-1 gap-4 ${showNavigation ? "py-6" : "py-10"}`}>
          {showNavigation && <Sidebar items={sidebarItems} collapsed={isSidebarCollapsed} />}

          <main className="flex-1 rounded-xl bg-white p-6 shadow-sm ring-1 ring-gray-200">{children}</main>
        </div>
      </div>

      {showNavigation ? (
        <Footer
          companyName="ESMS"
          version="0.1.0"
          links={[
            { label: "ヘルプセンター", href: "https://example.com/help" },
            { label: "利用規約", href: "https://example.com/terms" },
            { label: "プライバシー", href: "https://example.com/privacy" },
          ]}
        />
      ) : (
        <div className="py-6 text-center text-sm text-gray-500">
          <Link href="/">トップへ戻る</Link>
        </div>
      )}
    </div>
  );
};

export default AppLayout;
