import React from "react";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { ApprovalList } from "./ApprovalList";
import { Reservation } from "@/types/models";

type StubEventsHook = () => {
  events: Reservation[];
  isLoading: boolean;
  error: string | null;
  updateEvent: jest.Mock;
};

const events: Reservation[] = [
  {
    id: "p1",
    organizerId: "u1",
    title: "設備利用",
    description: "会議室予約",
    startAt: "2025-11-24T09:00:00Z",
    endAt: "2025-11-24T10:00:00Z",
    isPrivate: false,
    timezone: "Asia/Tokyo",
    approvalStatus: "PENDING",
    version: 1,
    createdAt: "2025-11-23T00:00:00Z",
    updatedAt: "2025-11-23T00:00:00Z",
  },
  {
    id: "c1",
    organizerId: "u2",
    title: "承認済み予約",
    description: "完了", 
    startAt: "2025-11-25T09:00:00Z",
    endAt: "2025-11-25T10:00:00Z",
    isPrivate: false,
    timezone: "Asia/Tokyo",
    approvalStatus: "CONFIRMED",
    version: 1,
    createdAt: "2025-11-23T00:00:00Z",
    updatedAt: "2025-11-23T00:00:00Z",
  },
];

const createStubHook = (updateEvent: jest.Mock): StubEventsHook => {
  return () => ({
    events,
    isLoading: false,
    error: null,
    updateEvent,
    createEvent: jest.fn(),
    fetchEvents: jest.fn(),
    deleteEvent: jest.fn(),
  });
};

describe("ApprovalList", () => {
  it("shows only pending approvals and handles approve/reject", async () => {
    const handleApprove = jest.fn();
    const handleReject = jest.fn();
    const updateEvent = jest.fn().mockResolvedValue(events[0]);

    render(
      <ApprovalList
        useEventsHook={createStubHook(updateEvent)}
        autoLoad={false}
        onApprove={handleApprove}
        onReject={handleReject}
      />
    );

    expect(screen.getByText("設備利用")).toBeInTheDocument();
    expect(screen.queryByText("承認済み予約")).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "承認" }));
    fireEvent.click(screen.getByRole("button", { name: "却下" }));

    await waitFor(() => expect(updateEvent).toHaveBeenCalledTimes(2));
    expect(handleApprove).toHaveBeenCalledWith("p1");
    expect(handleReject).toHaveBeenCalledWith("p1");
  });
});
