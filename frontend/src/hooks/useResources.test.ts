import { act, renderHook } from "@testing-library/react";

import { ApiClient } from "@/lib/api-client";
import { useResources } from "./useResources";
import { Resource } from "@/types/models";

const mockFetch = jest.fn<Promise<Response>, [RequestInfo | URL, RequestInit | undefined]>();
const apiClient = new ApiClient({ fetchImpl: mockFetch });

const resources: Resource[] = [
  {
    id: "room-1",
    name: "Room A",
    type: "MEETING_ROOM",
    capacity: 8,
    location: "Floor 1",
    equipment: {},
    requiredRole: null,
    isActive: true,
    createdAt: "2025-11-01T00:00:00Z",
    updatedAt: "2025-11-01T00:00:00Z",
  },
];

describe("useResources", () => {
  beforeEach(() => {
    mockFetch.mockReset();
  });

  it("searches resources on demand", async () => {
    mockFetch.mockResolvedValue(
      new Response(JSON.stringify({ data: resources }), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    const { result } = renderHook(() => useResources({ apiClient }));

    await act(async () => {
      await result.current.search({ keyword: "Room" });
    });

    expect(result.current.resources).toHaveLength(1);
    expect(result.current.resources[0].name).toBe("Room A");
  });

  it("checks availability", async () => {
    mockFetch.mockResolvedValueOnce(
      new Response(JSON.stringify({ data: resources }), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );
    const { result } = renderHook(() => useResources({ apiClient }));

    await act(async () => {
      await result.current.search();
    });

    mockFetch.mockResolvedValueOnce(
      new Response(JSON.stringify({ data: { available: true } }), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    const available = await result.current.checkAvailability("room-1", "2025-12-01", "2025-12-01");
    expect(available).toBe(true);
  });
});

