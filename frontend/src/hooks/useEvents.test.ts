import { act, renderHook, waitFor } from "@testing-library/react";

import { ApiClient } from "@/lib/api-client";
import { useEvents } from "./useEvents";
import { Reservation } from "@/types/models";

const mockFetch = jest.fn<Promise<Response>, [RequestInfo | URL, RequestInit | undefined]>();
const apiClient = new ApiClient({ fetchImpl: mockFetch });

const baseReservation: Reservation = {
  id: "res-1",
  organizerId: "org-1",
  title: "Weekly sync",
  description: "",
  startAt: "2025-12-01T00:00:00Z",
  endAt: "2025-12-01T01:00:00Z",
  timezone: "UTC",
  approvalStatus: "PENDING",
  version: 1,
  isPrivate: false,
  createdAt: "2025-11-01T00:00:00Z",
  updatedAt: "2025-11-01T00:00:00Z",
};

describe("useEvents", () => {
  beforeEach(() => {
    mockFetch.mockReset();
  });

  it("fetches events on autoLoad", async () => {
    mockFetch.mockResolvedValue(
      new Response(JSON.stringify({ data: [baseReservation] }), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    const { result } = renderHook(() => useEvents({ apiClient, autoLoad: true }));
    await waitFor(() => expect(result.current.events).toHaveLength(1));

    expect(result.current.events).toHaveLength(1);
    expect(result.current.events[0].id).toBe("res-1");
  });

  it("creates, updates, and deletes events", async () => {
    // initial fetch
    mockFetch.mockResolvedValueOnce(
      new Response(JSON.stringify({ data: [] }), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    const { result } = renderHook(() => useEvents({ apiClient, autoLoad: true }));
    await waitFor(() => expect(result.current.isLoading).toBe(false));

    // create
    mockFetch.mockResolvedValueOnce(
      new Response(JSON.stringify({ data: baseReservation }), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    await act(async () => {
      await result.current.createEvent({
        title: baseReservation.title,
        description: baseReservation.description,
        startAt: baseReservation.startAt,
        endAt: baseReservation.endAt,
        timezone: baseReservation.timezone,
      });
    });

    expect(result.current.events).toHaveLength(1);

    // update
    const updated = { ...baseReservation, title: "Updated title" };
    mockFetch.mockResolvedValueOnce(
      new Response(JSON.stringify({ data: updated }), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    await act(async () => {
      await result.current.updateEvent(updated.id, { title: updated.title });
    });

    expect(result.current.events[0].title).toBe("Updated title");

    // delete
    mockFetch.mockResolvedValueOnce(new Response(JSON.stringify({ data: null }), { status: 200 }));

    await act(async () => {
      await result.current.deleteEvent(updated.id);
    });

    expect(result.current.events).toHaveLength(0);
  });
});

