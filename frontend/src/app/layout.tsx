import type { Metadata } from "next";

import "../styles/globals.css";

import AppLayout from "@/components/layout/AppLayout";

export const metadata: Metadata = {
  title: "Enterprise Scheduler",
  description: "社内リソース・予約管理のためのフロントエンド",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body className="min-h-screen bg-gray-50 text-gray-900">
        <AppLayout>{children}</AppLayout>
      </body>
    </html>
  );
}
