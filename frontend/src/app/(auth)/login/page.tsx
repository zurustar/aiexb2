"use client";

import React, { useState } from "react";
import Link from "next/link";

import { Button } from "@/components/ui/Button";
import { useAuth, UseAuthResult } from "@/hooks/useAuth";

export type LoginPageProps = {
  useAuthHook?: () => UseAuthResult;
};

export const LoginPage: React.FC<LoginPageProps> = ({ useAuthHook = useAuth }) => {
  const auth = useAuthHook();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [statusMessage, setStatusMessage] = useState<string | null>(null);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatusMessage(null);
    try {
      await auth.login({ username, password });
      setStatusMessage("ログインに成功しました。ダッシュボードへ移動できます。");
    } catch (err) {
      setStatusMessage(auth.error ?? "認証に失敗しました。");
    }
  };

  return (
    <main className="mx-auto flex max-w-xl flex-col gap-6" aria-label="ログインページ">
      <section className="flex flex-col gap-2">
        <p className="text-sm font-semibold text-blue-700">Enterprise Scheduler</p>
        <h1 className="text-2xl font-bold text-gray-900">サインイン</h1>
        <p className="text-sm text-gray-700">社内アカウントでログインしてください。</p>
      </section>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4 rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <label className="flex flex-col gap-1 text-sm font-medium text-gray-800" htmlFor="username">
          ユーザー名
          <input
            id="username"
            name="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            className="rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200"
            placeholder="user@example.com"
          />
        </label>

        <label className="flex flex-col gap-1 text-sm font-medium text-gray-800" htmlFor="password">
          パスワード
          <input
            id="password"
            name="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200"
          />
        </label>

        {auth.error && (
          <p role="alert" className="text-sm text-red-600">
            {auth.error}
          </p>
        )}

        {statusMessage && !auth.error && <p className="text-sm text-green-700">{statusMessage}</p>}

        <div className="flex items-center justify-between">
          <div className="text-sm text-gray-600">
            パスワードを忘れた場合は管理者へお問い合わせください。
          </div>
          <Button type="submit" isLoading={auth.isLoading} disabled={auth.isLoading}>
            サインイン
          </Button>
        </div>
      </form>

      <div className="rounded-xl border border-gray-200 bg-white p-4 text-sm text-gray-700 shadow-sm">
        <p className="font-semibold">シングルサインオンを利用中ですか？</p>
        <p className="mt-1">IdPでの認証後、この画面には戻らず自動的にコールバックページへ遷移します。</p>
        <Link href="/callback" className="mt-2 inline-block font-semibold text-blue-700">
          コールバックページを開く
        </Link>
      </div>
    </main>
  );
};

export default function Page() {
  return <LoginPage />;
}
