import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";

import { EventsPage } from "@/app/events/page";
import { UseEventsResult } from "@/hooks/useEvents";
import { Reservation } from "@/types/models";

const baseReservation: Reservation = {
  id: "res-1",
  organizerId: "user-1",
  title: "キックオフ",
  description: "説明",
  startAt: "2025-12-01T10:00",
  endAt: "2025-12-01T11:00",
  timezone: "Asia/Tokyo",
  isPrivate: false,
  approvalStatus: "PENDING",
  version: 1,
  createdAt: "2025-11-01T00:00:00Z",
  updatedAt: "2025-11-01T00:00:00Z",
};

describe("E2E: 予約作成フロー", () => {
  it("creates a reservation and shows detail trigger", async () => {
    const fetchEvents = jest.fn();
    const createEvent = jest.fn(async () => ({ ...baseReservation, id: "res-2", title: "新規予約" }));

    const mockState: UseEventsResult = {
      events: [baseReservation],
      isLoading: false,
      error: null,
      fetchEvents,
      createEvent,
      updateEvent: jest.fn(),
      deleteEvent: jest.fn(),
    };

    const useEventsHook = () => mockState;

    render(<EventsPage useEventsHook={useEventsHook} />);

    fireEvent.change(screen.getByLabelText("タイトル"), { target: { value: "新規予約" } });
    fireEvent.change(screen.getByTestId("startAt"), { target: { value: "2025-12-02T09:00" } });
    fireEvent.change(screen.getByTestId("endAt"), { target: { value: "2025-12-02T10:00" } });
    fireEvent.click(screen.getByText("予約を作成"));

    await waitFor(() => expect(createEvent).toHaveBeenCalled());
    expect(fetchEvents).toHaveBeenCalled();
    expect(screen.getByTestId("reservation-list")).toBeInTheDocument();

    const list = screen.getByTestId("reservation-list");
    fireEvent.click(within(list).getByText(baseReservation.title));
    await waitFor(() => expect(screen.getByTestId("start-at")).toHaveTextContent(/2025\/12\/01/));
  });
});
