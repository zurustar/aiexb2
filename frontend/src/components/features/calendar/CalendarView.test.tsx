import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";

import { CalendarView } from "./CalendarView";
import { Reservation } from "@/types/models";

const sampleEvents: Reservation[] = [
  {
    id: "1",
    organizerId: "u1",
    title: "週次MTG",
    description: "開発チーム定例",
    startAt: "2025-11-24T09:00:00Z",
    endAt: "2025-11-24T10:00:00Z",
    isPrivate: false,
    timezone: "Asia/Tokyo",
    approvalStatus: "CONFIRMED",
    version: 1,
    createdAt: "2025-11-23T00:00:00Z",
    updatedAt: "2025-11-23T00:00:00Z",
  },
  {
    id: "2",
    organizerId: "u2",
    title: "1on1",
    description: "フォロー面談",
    startAt: "2025-11-25T11:00:00Z",
    endAt: "2025-11-25T11:30:00Z",
    isPrivate: true,
    timezone: "Asia/Tokyo",
    approvalStatus: "PENDING",
    version: 1,
    createdAt: "2025-11-23T00:00:00Z",
    updatedAt: "2025-11-23T00:00:00Z",
  },
];

type StubEventsHook = () => {
  events: Reservation[];
  isLoading: boolean;
  error: string | null;
  fetchEvents: jest.Mock;
};

const createStubHook = (overrides?: Partial<ReturnType<StubEventsHook>>) => {
  const fetchEvents = jest.fn();
  return () => ({
    events: sampleEvents,
    isLoading: false,
    error: null,
    fetchEvents,
    createEvent: jest.fn(),
    updateEvent: jest.fn(),
    deleteEvent: jest.fn(),
    ...(overrides ?? {}),
  });
};

describe("CalendarView", () => {
  it("shows events grouped by date", () => {
    render(<CalendarView useEventsHook={createStubHook()} fetchOnMount={false} />);

    expect(screen.getByText("週次MTG")).toBeInTheDocument();
    expect(screen.getByText("1on1")).toBeInTheDocument();
    expect(screen.getByText("2025-11-24")).toBeInTheDocument();
    expect(screen.getByText("2025-11-25")).toBeInTheDocument();
  });

  it("refreshes events when date is changed", () => {
    const fetchEvents = jest.fn();
    const stubHook = createStubHook({ fetchEvents });

    render(<CalendarView useEventsHook={stubHook} fetchOnMount={false} />);

    const dateInput = screen.getByLabelText("表示日");
    fireEvent.change(dateInput, { target: { value: "2025-11-24T00:00" } });

    expect(fetchEvents).toHaveBeenCalled();
  });
});
