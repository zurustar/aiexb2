import React from "react";

import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { formatDateTime } from "@/lib/utils";
import { Reservation } from "@/types/models";

export type ReservationDetailProps = {
  reservation: Reservation;
  isOpen: boolean;
  onClose: () => void;
  onEdit?: (reservation: Reservation) => void;
  onCancel?: (reservation: Reservation) => void;
  timezone?: string;
};

export const ReservationDetail: React.FC<ReservationDetailProps> = ({
  reservation,
  isOpen,
  onClose,
  onEdit,
  onCancel,
  timezone = "Asia/Tokyo",
}) => {
  const handleEdit = () => onEdit?.(reservation);
  const handleCancel = () => onCancel?.(reservation);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={reservation.title}
      footer={
        <div className="flex items-center justify-end gap-2">
          {onCancel && (
            <Button variant="ghost" onClick={handleCancel}>
              キャンセル
            </Button>
          )}
          {onEdit && (
            <Button variant="secondary" onClick={handleEdit}>
              編集
            </Button>
          )}
        </div>
      }
    >
      <div className="flex flex-col gap-3 text-sm text-gray-800">
        <div className="flex items-center gap-2">
          <span className="rounded-full bg-blue-50 px-3 py-1 text-xs font-semibold text-blue-700">
            {reservation.approvalStatus}
          </span>
          {reservation.isPrivate && <span className="text-xs text-gray-500">非公開</span>}
        </div>
        <dl className="grid grid-cols-3 gap-2">
          <dt className="text-gray-500">開始</dt>
          <dd className="col-span-2 font-medium" data-testid="start-at">
            {formatDateTime(reservation.startAt, "ja-JP", timezone)}
          </dd>
          <dt className="text-gray-500">終了</dt>
          <dd className="col-span-2 font-medium" data-testid="end-at">
            {formatDateTime(reservation.endAt, "ja-JP", timezone)}
          </dd>
          <dt className="text-gray-500">主催者</dt>
          <dd className="col-span-2">{reservation.organizer?.name ?? reservation.organizerId}</dd>
          <dt className="text-gray-500">説明</dt>
          <dd className="col-span-2 whitespace-pre-wrap">{reservation.description}</dd>
        </dl>
      </div>
    </Modal>
  );
};

export default ReservationDetail;
