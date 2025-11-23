"use client";

import React, { useMemo, useState } from "react";

import CalendarView from "@/components/features/calendar/CalendarView";
import ReservationDetail from "@/components/features/reservation/ReservationDetail";
import ReservationForm from "@/components/features/reservation/ReservationForm";
import { Button } from "@/components/ui/Button";
import { useEvents, UseEventsResult } from "@/hooks/useEvents";
import { Reservation } from "@/types/models";

export type EventsPageProps = {
  useEventsHook?: (params?: Parameters<typeof useEvents>[0]) => UseEventsResult;
};

export const EventsPage: React.FC<EventsPageProps> = ({ useEventsHook = useEvents }) => {
  const eventsState = useEventsHook({ autoLoad: true });
  const { events, isLoading, error, fetchEvents } = eventsState;
  const [selected, setSelected] = useState<Reservation | null>(null);
  const [lastCreated, setLastCreated] = useState<Reservation | null>(null);

  const sortedEvents = useMemo(() => {
    return [...events].sort((a, b) => a.startAt.localeCompare(b.startAt));
  }, [events]);

  const handleCreated = async (reservation: Reservation) => {
    setLastCreated(reservation);
    await fetchEvents();
  };

  return (
    <main className="flex flex-col gap-6" aria-label="予定管理">
      <header className="flex flex-col gap-2">
        <p className="text-sm font-semibold text-blue-700">予定管理</p>
        <h1 className="text-2xl font-bold text-gray-900">予約の作成と確認</h1>
        <p className="text-sm text-gray-700">会議やイベントを登録し、詳細を確認できます。</p>
      </header>

      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
        <section className="lg:col-span-1 rounded-lg border border-gray-200 bg-white p-5 shadow-sm" aria-label="予約作成フォーム">
          <h2 className="text-lg font-semibold text-gray-900">新規予約</h2>
          <p className="text-sm text-gray-700">開始・終了時間を指定して予約を作成します。</p>
          <div className="mt-4">
            <ReservationForm
              useEventsHook={() => eventsState}
              onCreated={handleCreated}
            />
          </div>
          {lastCreated && (
            <p className="mt-3 rounded-md bg-green-50 px-3 py-2 text-sm text-green-800" data-testid="create-success">
              「{lastCreated.title}」を作成しました。
            </p>
          )}
        </section>

        <section className="lg:col-span-2 rounded-lg border border-gray-200 bg-white p-5 shadow-sm" aria-label="予約一覧">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-lg font-semibold text-gray-900">登録済みの予約</h2>
              <p className="text-sm text-gray-700">クリックすると詳細を表示します。</p>
            </div>
            <Button variant="ghost" onClick={() => void fetchEvents()} isLoading={isLoading}>
              更新
            </Button>
          </div>
          {sortedEvents.length === 0 ? (
            <p className="mt-3 text-sm text-gray-700">表示できる予約がありません。</p>
          ) : (
            <ul className="mt-3 divide-y divide-gray-200" data-testid="reservation-list">
              {sortedEvents.map((reservation) => (
                <li
                  key={reservation.id}
                  className="flex cursor-pointer items-center justify-between py-3 hover:text-blue-700"
                  onClick={() => setSelected(reservation)}
                >
                  <div>
                    <p className="font-semibold">{reservation.title}</p>
                    <p className="text-xs text-gray-500">{reservation.startAt} 〜 {reservation.endAt}</p>
                  </div>
                  <span className="rounded-full bg-blue-50 px-2 py-1 text-xs font-semibold text-blue-700">{reservation.approvalStatus}</span>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

      <section className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm" aria-label="カレンダー">
        <CalendarView fetchOnMount useEventsHook={() => eventsState} />
      </section>

      {selected && (
        <ReservationDetail
          reservation={selected}
          isOpen={true}
          onClose={() => setSelected(null)}
          onEdit={() => void fetchEvents()}
          onCancel={() => void fetchEvents()}
        />
      )}
    </main>
  );
};

export default function Page() {
  return <EventsPage />;
}
