import React from "react";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { ReservationForm } from "./ReservationForm";
import { Reservation } from "@/types/models";

type StubEventsHook = () => {
  createEvent: jest.Mock;
  isLoading: boolean;
  error: string | null;
};

const sampleReservation: Reservation = {
  id: "r1",
  organizerId: "u1",
  title: "ミーティング",
  description: "仕様調整",
  startAt: "2025-11-24T09:00:00Z",
  endAt: "2025-11-24T10:00:00Z",
  isPrivate: false,
  timezone: "Asia/Tokyo",
  approvalStatus: "PENDING",
  version: 1,
  createdAt: "2025-11-23T00:00:00Z",
  updatedAt: "2025-11-23T00:00:00Z",
};

const createStubHook = (createEvent: jest.Mock): StubEventsHook => {
  return () => ({
    createEvent,
    isLoading: false,
    error: null,
    events: [],
    fetchEvents: jest.fn(),
    updateEvent: jest.fn(),
    deleteEvent: jest.fn(),
  });
};

describe("ReservationForm", () => {
  it("submits payload and notifies callback", async () => {
    const handleCreated = jest.fn();
    const createEvent = jest.fn().mockResolvedValue(sampleReservation);

    render(<ReservationForm useEventsHook={createStubHook(createEvent)} onCreated={handleCreated} />);

    fireEvent.change(screen.getByLabelText("タイトル"), { target: { value: "ミーティング" } });
    fireEvent.change(screen.getByLabelText("説明"), { target: { value: "仕様調整" } });
    fireEvent.change(screen.getByLabelText("開始日時"), { target: { value: "2025-11-24T09:00" } });
    fireEvent.change(screen.getByLabelText("終了日時"), { target: { value: "2025-11-24T10:00" } });
    fireEvent.change(screen.getByLabelText("タイムゾーン"), { target: { value: "Asia/Tokyo" } });
    fireEvent.click(screen.getByLabelText("非公開予約"));

    fireEvent.submit(screen.getByLabelText("予約フォーム"));

    await waitFor(() => expect(createEvent).toHaveBeenCalled());
    expect(createEvent).toHaveBeenCalledWith(
      expect.objectContaining({
        title: "ミーティング",
        description: "仕様調整",
        startAt: "2025-11-24T09:00",
        endAt: "2025-11-24T10:00",
        isPrivate: true,
      })
    );
    expect(handleCreated).toHaveBeenCalledWith(sampleReservation);
  });

  it("shows validation error when end time is before start time", () => {
    const createEvent = jest.fn();

    render(<ReservationForm useEventsHook={createStubHook(createEvent)} />);

    fireEvent.change(screen.getByLabelText("タイトル"), { target: { value: "テスト" } });
    fireEvent.change(screen.getByLabelText("開始日時"), { target: { value: "2025-11-24T10:00" } });
    fireEvent.change(screen.getByLabelText("終了日時"), { target: { value: "2025-11-24T09:00" } });

    fireEvent.submit(screen.getByLabelText("予約フォーム"));

    expect(screen.getByText("終了日時は開始日時より後に設定してください。")).toBeInTheDocument();
    expect(createEvent).not.toHaveBeenCalled();
  });
});
