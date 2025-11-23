import React, { FormEvent, useState } from "react";

import { Button } from "@/components/ui/Button";
import { DatePicker } from "@/components/ui/DatePicker";
import { EventPayload, useEvents } from "@/hooks/useEvents";
import { Reservation } from "@/types/models";

export type ReservationFormProps = {
  initialValues?: Partial<EventPayload>;
  timezoneOptions?: string[];
  onCreated?: (reservation: Reservation) => void;
  useEventsHook?: typeof useEvents;
};

export const ReservationForm: React.FC<ReservationFormProps> = ({
  initialValues,
  timezoneOptions = ["Asia/Tokyo", "UTC"],
  onCreated,
  useEventsHook = useEvents,
}) => {
  const { createEvent, isLoading, error } = useEventsHook({ autoLoad: false });
  const [form, setForm] = useState<Required<Pick<EventPayload, "title" | "description" | "startAt" | "endAt" | "timezone">> & {
    isPrivate: boolean;
  }>(
    {
      title: initialValues?.title ?? "",
      description: initialValues?.description ?? "",
      startAt: initialValues?.startAt ?? "",
      endAt: initialValues?.endAt ?? "",
      timezone: initialValues?.timezone ?? "Asia/Tokyo",
      isPrivate: initialValues?.isPrivate ?? false,
    }
  );
  const [validationError, setValidationError] = useState<string | null>(null);

  const updateField = <K extends keyof typeof form>(key: K, value: (typeof form)[K]) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setValidationError(null);

    if (!form.title || !form.startAt || !form.endAt) {
      setValidationError("必須項目を入力してください。");
      return;
    }

    if (new Date(form.startAt) >= new Date(form.endAt)) {
      setValidationError("終了日時は開始日時より後に設定してください。");
      return;
    }

    const payload: EventPayload = {
      title: form.title,
      description: form.description,
      startAt: form.startAt,
      endAt: form.endAt,
      timezone: form.timezone,
      isPrivate: form.isPrivate,
    };

    const created = await createEvent(payload);
    onCreated?.(created);
  };

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit} aria-label="予約フォーム">
      <div>
        <label className="block text-sm font-medium text-gray-800" htmlFor="title">
          タイトル
        </label>
        <input
          id="title"
          name="title"
          value={form.title}
          onChange={(e) => updateField("title", e.target.value)}
          className="mt-1 w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200"
          placeholder="例: プロジェクトキックオフ"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-800" htmlFor="description">
          説明
        </label>
        <textarea
          id="description"
          name="description"
          value={form.description}
          onChange={(e) => updateField("description", e.target.value)}
          className="mt-1 w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200"
          rows={3}
        />
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <DatePicker
          label="開始日時"
          name="startAt"
          value={form.startAt}
          onChange={(value) => updateField("startAt", value)}
          required
        />
        <DatePicker
          label="終了日時"
          name="endAt"
          value={form.endAt}
          onChange={(value) => updateField("endAt", value)}
          required
        />
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <label className="flex flex-col gap-1 text-sm font-medium text-gray-800" htmlFor="timezone">
          タイムゾーン
          <select
            id="timezone"
            name="timezone"
            value={form.timezone}
            onChange={(e) => updateField("timezone", e.target.value)}
            className="rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200"
          >
            {timezoneOptions.map((tz) => (
              <option key={tz} value={tz}>
                {tz}
              </option>
            ))}
          </select>
        </label>

        <label className="flex items-center gap-2 text-sm font-medium text-gray-800" htmlFor="isPrivate">
          <input
            id="isPrivate"
            type="checkbox"
            checked={form.isPrivate}
            onChange={(e) => updateField("isPrivate", e.target.checked)}
            className="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          非公開予約
        </label>
      </div>

      {validationError && (
        <p className="text-sm text-red-600" role="alert">
          {validationError}
        </p>
      )}
      {error && (
        <p className="text-sm text-red-600" role="alert">
          {error}
        </p>
      )}

      <div className="flex items-center justify-end gap-2">
        <Button type="submit" isLoading={isLoading} disabled={isLoading}>
          予約を作成
        </Button>
      </div>
    </form>
  );
};

export default ReservationForm;
