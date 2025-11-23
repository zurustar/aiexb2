"use client";

import Link from "next/link";

import { Button } from "@/components/ui/Button";
import { useAuth } from "@/hooks/useAuth";

const featureCards = [
  {
    title: "ダッシュボード",
    description: "承認待ちや本日の予定を一目で把握",
    href: "/dashboard",
  },
  {
    title: "予定管理",
    description: "会議やイベントの作成・詳細確認を実施",
    href: "/events",
  },
  {
    title: "リソース管理",
    description: "会議室や備品の空き状況を検索",
    href: "/resources",
  },
];

export default function HomePage() {
  const { isAuthenticated, user } = useAuth();
  const greeting = isAuthenticated ? `${user?.name ?? "ユーザー"} さん、ようこそ！` : "ゲストとして閲覧しています";

  return (
    <main className="flex flex-col gap-8" aria-label="トップページ">
      <section className="rounded-xl bg-gradient-to-r from-blue-600 to-blue-500 p-8 text-white shadow">
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="space-y-2">
            <p className="text-sm text-blue-100">Enterprise Scheduler</p>
            <h1 className="text-3xl font-bold leading-tight">社内の予定とリソースを一元管理</h1>
            <p className="max-w-2xl text-blue-50">{greeting}</p>
          </div>
          <div className="flex gap-3">
            <Link href="/login" aria-label="ログインページへ">
              <Button variant="secondary">ログイン</Button>
            </Link>
            <Link href="/dashboard" aria-label="ダッシュボードへ">
              <Button>ダッシュボードを開く</Button>
            </Link>
          </div>
        </div>
      </section>

      <section className="grid grid-cols-1 gap-4 md:grid-cols-3" aria-label="主要機能への導線">
        {featureCards.map((card) => (
          <Link key={card.href} href={card.href} className="group block h-full">
            <article className="flex h-full flex-col gap-3 rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition group-hover:-translate-y-0.5 group-hover:shadow">
              <h2 className="text-lg font-semibold text-gray-900">{card.title}</h2>
              <p className="flex-1 text-sm text-gray-700">{card.description}</p>
              <span className="text-sm font-semibold text-blue-700">もっと見る →</span>
            </article>
          </Link>
        ))}
      </section>

      <section className="grid grid-cols-1 gap-4 md:grid-cols-2" aria-label="運用メモ">
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <h3 className="text-base font-semibold text-gray-900">今日のタスク</h3>
          <ul className="mt-3 list-disc space-y-2 pl-5 text-sm text-gray-700">
            <li>承認待ちのリクエストを確認</li>
            <li>会議室の空き状況をチェック</li>
            <li>重要イベントをダッシュボードへ固定</li>
          </ul>
        </div>
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <h3 className="text-base font-semibold text-gray-900">ガイド</h3>
          <p className="mt-2 text-sm text-gray-700">
            予約や承認の操作はダッシュボードから行えます。予約作成は「予定管理」、リソース確認は「リソース管理」に移動してください。
          </p>
        </div>
      </section>
    </main>
  );
}
