import React, { useMemo, useState } from "react";

import { DatePicker } from "@/components/ui/DatePicker";
import { useEvents } from "@/hooks/useEvents";
import { formatDateTime } from "@/lib/utils";
import { Reservation } from "@/types/models";

export type CalendarViewProps = {
  timezone?: string;
  fetchOnMount?: boolean;
  useEventsHook?: typeof useEvents;
};

export const CalendarView: React.FC<CalendarViewProps> = ({
  timezone = "Asia/Tokyo",
  fetchOnMount = true,
  useEventsHook = useEvents,
}) => {
  const [selectedDate, setSelectedDate] = useState<string>("");
  const { events, isLoading, error, fetchEvents } = useEventsHook({ autoLoad: fetchOnMount });

  const filteredEvents = useMemo(() => {
    if (!selectedDate) return events;
    const dateOnly = selectedDate.slice(0, 10);
    return events.filter((event) => event.startAt.startsWith(dateOnly));
  }, [events, selectedDate]);

  const groupedByDay = useMemo(() => {
    return filteredEvents.reduce<Record<string, Reservation[]>>((acc, event) => {
      const date = event.startAt.slice(0, 10);
      acc[date] = acc[date] ? [...acc[date], event] : [event];
      return acc;
    }, {});
  }, [filteredEvents]);

  const handleDateChange = (value: string) => {
    setSelectedDate(value);
    void fetchEvents();
  };

  const renderEvents = (eventsForDay: Reservation[]) => {
    return eventsForDay
      .sort((a, b) => a.startAt.localeCompare(b.startAt))
      .map((event) => (
        <li key={event.id} className="rounded-md border border-gray-200 bg-white px-4 py-3 shadow-sm">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-semibold text-gray-900">{event.title}</p>
              <p className="text-xs text-gray-600">{event.description}</p>
            </div>
            <span className="rounded-full bg-blue-50 px-3 py-1 text-xs font-semibold text-blue-700">{event.approvalStatus}</span>
          </div>
          <p className="mt-2 text-sm text-gray-700">
            {formatDateTime(event.startAt, "ja-JP", timezone)} 〜 {formatDateTime(event.endAt, "ja-JP", timezone)}
          </p>
        </li>
      ));
  };

  return (
    <section aria-label="カレンダー" className="flex flex-col gap-4">
      <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">カレンダー</h2>
          <p className="text-sm text-gray-600">予定を日付ごとに確認できます</p>
        </div>
        <DatePicker
          label="表示日"
          name="calendar-date"
          value={selectedDate}
          onChange={handleDateChange}
          helperText="日付を指定して予定を絞り込み"
        />
      </div>

      {isLoading && <p className="text-sm text-gray-600">読み込み中...</p>}
      {error && (
        <p className="text-sm text-red-600" role="alert">
          {error}
        </p>
      )}

      {filteredEvents.length === 0 && !isLoading ? (
        <p className="text-sm text-gray-600">表示する予定がありません。</p>
      ) : (
        <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
          {Object.entries(groupedByDay).map(([date, eventsForDay]) => (
            <div key={date} className="flex flex-col gap-2 rounded-lg border border-gray-200 bg-gray-50 p-3">
              <p className="text-sm font-semibold text-gray-800">{date}</p>
              <ul className="flex flex-col gap-2">{renderEvents(eventsForDay)}</ul>
            </div>
          ))}
        </div>
      )}
    </section>
  );
};

export default CalendarView;
