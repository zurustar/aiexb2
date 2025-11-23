import { useCallback, useEffect, useMemo, useState } from "react";

import { ApiClient } from "@/lib/api-client";
import { PaginatedResponse } from "@/types/api";
import { Reservation } from "@/types/models";

export type EventFilters = {
  page?: number;
  pageSize?: number;
  organizerId?: string;
};

export type EventPayload = Pick<Reservation, "title" | "description" | "startAt" | "endAt" | "timezone"> &
  Partial<Pick<Reservation, "rrule" | "isPrivate" | "approvalStatus">> & {
    resourceIds?: string[];
    participantIds?: string[];
  };

export type UseEventsOptions = {
  apiClient?: ApiClient;
  autoLoad?: boolean;
  initialFilters?: EventFilters;
};

export type UseEventsResult = {
  events: Reservation[];
  isLoading: boolean;
  error: string | null;
  fetchEvents: (filters?: EventFilters) => Promise<void>;
  createEvent: (payload: EventPayload) => Promise<Reservation>;
  updateEvent: (id: string, payload: Partial<EventPayload>) => Promise<Reservation>;
  deleteEvent: (id: string) => Promise<void>;
};

const defaultClient = new ApiClient();

export const useEvents = (options: UseEventsOptions = {}): UseEventsResult => {
  const apiClient = options.apiClient ?? defaultClient;
  const [events, setEvents] = useState<Reservation[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(!!options.autoLoad);
  const [error, setError] = useState<string | null>(null);

  const fetchEvents = useCallback(
    async (filters: EventFilters = {}) => {
      setIsLoading(true);
      setError(null);
      try {
        const params = new URLSearchParams();
        if (filters.page) params.set("page", String(filters.page));
        if (filters.pageSize) params.set("pageSize", String(filters.pageSize));
        if (filters.organizerId) params.set("organizerId", filters.organizerId);
        const query = params.toString();
        const path = query ? `/api/v1/reservations?${query}` : "/api/v1/reservations";
        const response = await apiClient.get<PaginatedResponse<Reservation>>(path);
        const list = Array.isArray(response.data) ? response.data : response.data.data;
        setEvents(list);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to fetch events");
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [apiClient]
  );

  const createEvent = useCallback(
    async (payload: EventPayload) => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await apiClient.post<Reservation>("/api/v1/reservations", payload);
        setEvents((prev) => [...prev, response.data]);
        return response.data;
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to create event");
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [apiClient]
  );

  const updateEvent = useCallback(
    async (id: string, payload: Partial<EventPayload>) => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await apiClient.patch<Reservation>(`/api/v1/reservations/${id}`, payload);
        setEvents((prev) => prev.map((event) => (event.id === id ? response.data : event)));
        return response.data;
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to update event");
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [apiClient]
  );

  const deleteEvent = useCallback(
    async (id: string) => {
      setIsLoading(true);
      setError(null);
      try {
        await apiClient.delete<null>(`/api/v1/reservations/${id}`);
        setEvents((prev) => prev.filter((event) => event.id !== id));
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to delete event");
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [apiClient]
  );

  useEffect(() => {
    if (options.autoLoad) {
      void fetchEvents(options.initialFilters);
    }
  }, [options.autoLoad, options.initialFilters, fetchEvents]);

  return useMemo(
    () => ({
      events,
      isLoading,
      error,
      fetchEvents,
      createEvent,
      updateEvent,
      deleteEvent,
    }),
    [events, isLoading, error, fetchEvents, createEvent, updateEvent, deleteEvent]
  );
};

