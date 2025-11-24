import { ApiErrorDetail, ApiErrorResponse, ApiResponse } from "@/types/api";

type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export type ApiClientOptions = {
  baseUrl?: string;
  getAuthToken?: () => string | null;
  fetchImpl?: typeof fetch;
  defaultHeaders?: HeadersInit;
  timeout?: number; // milliseconds, default: 30000
};

export class ApiClientError extends Error {
  public readonly status: number;
  public readonly details?: ApiErrorDetail[];
  public readonly traceId?: string;

  constructor(message: string, status: number, details?: ApiErrorDetail[], traceId?: string) {
    super(message);
    this.name = "ApiClientError";
    this.status = status;
    this.details = details;
    this.traceId = traceId;
  }
}

export class ApiClient {
  private readonly baseUrl: string;
  private readonly getAuthToken?: () => string | null;
  private readonly fetchImpl: typeof fetch;
  private readonly defaultHeaders: HeadersInit;
  private readonly timeout: number;

  constructor(options: ApiClientOptions = {}) {
    this.baseUrl = options.baseUrl ?? "";
    this.getAuthToken = options.getAuthToken;
    this.fetchImpl = options.fetchImpl ?? fetch;
    this.defaultHeaders = options.defaultHeaders ?? {
      "Content-Type": "application/json",
      Accept: "application/json",
    };
    this.timeout = options.timeout ?? 30000; // 30 seconds default
  }

  async get<T>(path: string): Promise<ApiResponse<T>> {
    return this.request<T>(path, { method: "GET" });
  }

  async post<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(path, {
      method: "POST",
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  }

  async put<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(path, {
      method: "PUT",
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  }

  async patch<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(path, {
      method: "PATCH",
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  }

  async delete<T>(path: string): Promise<ApiResponse<T>> {
    return this.request<T>(path, { method: "DELETE" });
  }

  private buildUrl(path: string): string {
    if (!this.baseUrl) return path;
    if (path.startsWith("http")) return path;
    const joiner = path.startsWith("/") ? "" : "/";
    return `${this.baseUrl}${joiner}${path}`;
  }

  private buildHeaders(extra?: HeadersInit): HeadersInit {
    const token = this.getAuthToken?.();
    const headers: HeadersInit = { ...this.defaultHeaders, ...extra };
    if (token) {
      return { ...headers, Authorization: `Bearer ${token}` };
    }
    return headers;
  }

  private async parseError(response: Response): Promise<ApiClientError> {
    try {
      const payload = (await response.json()) as ApiErrorResponse;
      const message = payload.error?.message ?? response.statusText;
      const status = payload.error?.status ?? response.status;
      return new ApiClientError(message, status, payload.error?.details, payload.traceId);
    } catch (err) {
      return new ApiClientError(response.statusText, response.status);
    }
  }

  private async request<T>(path: string, init: RequestInit & { method: HttpMethod }): Promise<ApiResponse<T>> {
    const url = this.buildUrl(path);
    const headers = this.buildHeaders(init.headers);

    // Setup AbortController with timeout
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await this.fetchImpl(url, {
        ...init,
        headers,
        signal: controller.signal
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        throw await this.parseError(response);
      }

      const contentType = response.headers.get("content-type") ?? "";
      if (contentType.includes("application/json")) {
        return (await response.json()) as ApiResponse<T>;
      }

      // Fallback for empty body or unexpected content type
      return { data: (undefined as unknown) as T };
    } catch (error) {
      clearTimeout(timeoutId);

      // Handle AbortError (timeout)
      if (error instanceof Error && error.name === "AbortError") {
        throw new ApiClientError(
          "Request timeout",
          408,
          [{ field: "_request", code: "TIMEOUT", message: "Request timed out" }],
          undefined
        );
      }

      // Re-throw ApiClientError as-is
      if (error instanceof ApiClientError) {
        throw error;
      }

      // Wrap other errors
      throw new ApiClientError(
        error instanceof Error ? error.message : "Unknown error",
        500,
        undefined,
        undefined
      );
    }
  }
}
