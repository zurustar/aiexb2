import { useCallback, useEffect, useMemo, useState } from "react";

import { ApiClient } from "@/lib/api-client";
import { PaginatedResponse } from "@/types/api";
import { Resource, ResourceType, Role } from "@/types/models";

type ResourceFilters = {
  keyword?: string;
  type?: ResourceType;
  requiredRole?: Role;
  capacity?: number;
};

export type UseResourcesOptions = {
  apiClient?: ApiClient;
  autoLoad?: boolean;
  initialFilters?: ResourceFilters;
};

export type UseResourcesResult = {
  resources: Resource[];
  isLoading: boolean;
  error: string | null;
  search: (filters?: ResourceFilters) => Promise<void>;
  checkAvailability: (id: string, startAt: string, endAt: string) => Promise<boolean>;
};

const defaultClient = new ApiClient();

export const useResources = (options: UseResourcesOptions = {}): UseResourcesResult => {
  const apiClient = options.apiClient ?? defaultClient;
  const [resources, setResources] = useState<Resource[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(!!options.autoLoad);
  const [error, setError] = useState<string | null>(null);

  const search = useCallback(
    async (filters: ResourceFilters = {}) => {
      setIsLoading(true);
      setError(null);
      try {
        const params = new URLSearchParams();
        if (filters.keyword) params.set("keyword", filters.keyword);
        if (filters.type) params.set("type", filters.type);
        if (filters.requiredRole) params.set("requiredRole", filters.requiredRole);
        if (filters.capacity) params.set("capacity", String(filters.capacity));
        const query = params.toString();
        const path = query ? `/api/v1/resources?${query}` : "/api/v1/resources";
        const response = await apiClient.get<PaginatedResponse<Resource>>(path);
        const list = Array.isArray(response.data) ? response.data : response.data.data;
        setResources(list);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to fetch resources");
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [apiClient]
  );

  const checkAvailability = useCallback(
    async (id: string, startAt: string, endAt: string) => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await apiClient.get<{ data: { available: boolean } }>(
          `/api/v1/resources/${id}/availability?startAt=${encodeURIComponent(startAt)}&endAt=${encodeURIComponent(endAt)}`
        );
        return response.data.available;
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to check availability");
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [apiClient]
  );

  useEffect(() => {
    if (options.autoLoad) {
      void search(options.initialFilters);
    }
  }, [options.autoLoad, options.initialFilters, search]);

  return useMemo(
    () => ({ resources, isLoading, error, search, checkAvailability }),
    [resources, isLoading, error, search, checkAvailability]
  );
};

