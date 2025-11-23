"use client";

import React, { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/Button";
import { useAuth, UseAuthResult } from "@/hooks/useAuth";

export type CallbackPageProps = {
  useAuthHook?: () => UseAuthResult;
};

export const CallbackPage: React.FC<CallbackPageProps> = ({ useAuthHook = useAuth }) => {
  const auth = useAuthHook();
  const params = useSearchParams();
  const [status, setStatus] = useState<"verifying" | "completed" | "failed">("verifying");

  const code = params.get("code");
  const state = params.get("state");

  useEffect(() => {
    auth.refresh();
  }, [auth]);

  useEffect(() => {
    if (auth.isLoading) return;
    if (auth.isAuthenticated) {
      setStatus("completed");
    } else if (auth.error) {
      setStatus("failed");
    }
  }, [auth.isAuthenticated, auth.isLoading, auth.error]);

  const statusMessage = useMemo(() => {
    if (status === "verifying") return "認証コードを検証しています...";
    if (status === "completed") return "ログインが完了しました。ダッシュボードへ遷移してください。";
    return auth.error ?? "セッションを確認できませんでした。再度ログインしてください。";
  }, [status, auth.error]);

  return (
    <main className="mx-auto flex max-w-2xl flex-col gap-6" aria-label="認証コールバック">
      <section className="rounded-xl border border-blue-100 bg-blue-50 p-5 text-blue-900 shadow-sm">
        <p className="text-sm font-semibold text-blue-700">外部IdPから復帰しました</p>
        <h1 className="text-2xl font-bold">セッションを確立しています</h1>
        <p className="mt-2 text-sm">セキュリティのため、ブラウザを閉じる前に必ずログアウトしてください。</p>
      </section>

      <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <dl className="grid grid-cols-3 gap-2 text-sm text-gray-800">
          <dt className="text-gray-500">code</dt>
          <dd className="col-span-2 break-all">{code ?? "(未受領)"}</dd>
          <dt className="text-gray-500">state</dt>
          <dd className="col-span-2 break-all">{state ?? "(未受領)"}</dd>
        </dl>

        <p className="mt-4 text-sm" data-testid="status-message">
          {statusMessage}
        </p>

        <div className="mt-4 flex flex-wrap gap-3">
          <Link href="/dashboard">
            <Button variant="secondary" disabled={status === "verifying"}>
              ダッシュボードへ移動
            </Button>
          </Link>
          <Link href="/login">
            <Button variant="ghost">ログイン画面へ戻る</Button>
          </Link>
        </div>
      </div>
    </main>
  );
};

export default function Page() {
  return <CallbackPage />;
}
