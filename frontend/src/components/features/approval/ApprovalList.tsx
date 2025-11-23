import React from "react";

import { Button } from "@/components/ui/Button";
import { useEvents } from "@/hooks/useEvents";
import { formatDateTime } from "@/lib/utils";
import { Reservation } from "@/types/models";

export type ApprovalListProps = {
  autoLoad?: boolean;
  useEventsHook?: typeof useEvents;
  onApprove?: (id: string) => void;
  onReject?: (id: string) => void;
  timezone?: string;
};

export const ApprovalList: React.FC<ApprovalListProps> = ({
  autoLoad = true,
  useEventsHook = useEvents,
  onApprove,
  onReject,
  timezone = "Asia/Tokyo",
}) => {
  const { events, isLoading, error, updateEvent } = useEventsHook({ autoLoad });
  const approvalRequests = events.filter((event) => event.approvalStatus === "PENDING");

  const handleApprove = async (reservation: Reservation) => {
    await updateEvent(reservation.id, { approvalStatus: "CONFIRMED" });
    onApprove?.(reservation.id);
  };

  const handleReject = async (reservation: Reservation) => {
    await updateEvent(reservation.id, { approvalStatus: "REJECTED" });
    onReject?.(reservation.id);
  };

  return (
    <section className="flex flex-col gap-3" aria-label="承認待ち一覧">
      <div>
        <h2 className="text-lg font-semibold text-gray-900">承認依頼</h2>
        <p className="text-sm text-gray-600">未処理の予約リクエストを確認します</p>
      </div>

      {isLoading && <p className="text-sm text-gray-600">読み込み中...</p>}
      {error && (
        <p className="text-sm text-red-600" role="alert">
          {error}
        </p>
      )}

      {approvalRequests.length === 0 && !isLoading ? (
        <p className="text-sm text-gray-600">承認待ちのリクエストはありません。</p>
      ) : (
        <ul className="flex flex-col gap-3">
          {approvalRequests.map((reservation) => (
            <li key={reservation.id} className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-semibold text-gray-900">{reservation.title}</p>
                  <p className="text-xs text-gray-600">
                    申請者: {reservation.organizer?.name ?? reservation.organizerId}
                  </p>
                </div>
                <span className="rounded-full bg-yellow-50 px-3 py-1 text-xs font-semibold text-yellow-700">PENDING</span>
              </div>
              <p className="mt-2 text-sm text-gray-700">
                {formatDateTime(reservation.startAt, "ja-JP", timezone)} 〜 {formatDateTime(reservation.endAt, "ja-JP", timezone)}
              </p>
              <div className="mt-3 flex items-center gap-2">
                <Button size="sm" variant="secondary" onClick={() => void handleApprove(reservation)}>
                  承認
                </Button>
                <Button size="sm" variant="ghost" onClick={() => void handleReject(reservation)}>
                  却下
                </Button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
};

export default ApprovalList;
