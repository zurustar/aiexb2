"use client";

import React, { useMemo } from "react";
import Link from "next/link";

import ApprovalList from "@/components/features/approval/ApprovalList";
import CalendarView from "@/components/features/calendar/CalendarView";
import { Button } from "@/components/ui/Button";
import { useEvents, UseEventsResult } from "@/hooks/useEvents";

export type DashboardPageProps = {
  useEventsHook?: (params?: Parameters<typeof useEvents>[0]) => UseEventsResult;
};

export const DashboardPage: React.FC<DashboardPageProps> = ({ useEventsHook = useEvents }) => {
  const { events, isLoading, error, fetchEvents } = useEventsHook({ autoLoad: true });

  const stats = useMemo(() => {
    const pending = events.filter((event) => event.approvalStatus === "PENDING").length;
    const confirmed = events.filter((event) => event.approvalStatus === "CONFIRMED").length;
    const upcoming = events.slice(0, 5);
    return { pending, confirmed, upcoming };
  }, [events]);

  return (
    <main className="flex flex-col gap-6" aria-label="ダッシュボード">
      <header className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <p className="text-sm font-semibold text-blue-700">ダッシュボード</p>
          <h1 className="text-2xl font-bold text-gray-900">本日の状況</h1>
          <p className="text-sm text-gray-700">承認依頼と予定のサマリーを確認できます。</p>
        </div>
        <div className="flex gap-3">
          <Link href="/events">
            <Button variant="secondary">予定を作成</Button>
          </Link>
          <Button variant="ghost" onClick={() => void fetchEvents() } isLoading={isLoading}>
            最新の状態に更新
          </Button>
        </div>
      </header>

      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}

      <section className="grid grid-cols-1 gap-4 md:grid-cols-3" aria-label="ダッシュボード統計">
        <article className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
          <p className="text-sm text-gray-600">承認待ち</p>
          <p className="text-3xl font-bold text-yellow-600" data-testid="pending-count">
            {stats.pending}
          </p>
        </article>
        <article className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
          <p className="text-sm text-gray-600">確定済み</p>
          <p className="text-3xl font-bold text-green-700">{stats.confirmed}</p>
        </article>
        <article className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
          <p className="text-sm text-gray-600">登録済みイベント</p>
          <p className="text-3xl font-bold text-blue-700">{events.length}</p>
        </article>
      </section>

      <section className="grid grid-cols-1 gap-4 lg:grid-cols-2" aria-label="予定と承認">
        <div className="rounded-lg border border-gray-200 bg-gray-50 p-4 shadow-inner">
          <CalendarView fetchOnMount useEventsHook={useEventsHook} />
        </div>
        <div id="approvals" className="rounded-lg border border-gray-200 bg-gray-50 p-4 shadow-inner">
          <ApprovalList useEventsHook={useEventsHook} />
        </div>
      </section>

      <section aria-label="直近の予定" className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900">直近の予定</h2>
          <Link href="/events" className="text-sm font-semibold text-blue-700">
            すべて見る
          </Link>
        </div>
        {stats.upcoming.length === 0 ? (
          <p className="mt-2 text-sm text-gray-700">表示できる予定がありません。</p>
        ) : (
          <ul className="mt-2 divide-y divide-gray-200 text-sm text-gray-800" data-testid="upcoming-list">
            {stats.upcoming.map((event) => (
              <li key={event.id} className="flex items-center justify-between py-2">
                <div className="flex flex-col">
                  <span className="font-semibold">{event.title}</span>
                  <span className="text-xs text-gray-500">{event.startAt} 〜 {event.endAt}</span>
                </div>
                <span className="rounded-full bg-blue-50 px-2 py-1 text-xs font-semibold text-blue-700">{event.approvalStatus}</span>
              </li>
            ))}
          </ul>
        )}
      </section>
    </main>
  );
};

export default function Page() {
  return <DashboardPage />;
}
