import { ApiClient, ApiClientError } from "./api-client";
import { ApiResponse } from "@/types/api";

describe("ApiClient", () => {
  const mockFetch = jest.fn<Promise<Response>, [RequestInfo | URL, RequestInit | undefined]>();
  const client = new ApiClient({ baseUrl: "https://api.example.com", fetchImpl: mockFetch, getAuthToken: () => "token" });

  beforeEach(() => {
    mockFetch.mockReset();
    jest.clearAllTimers();
  });

  it("prefixes the base URL and attaches authorization headers", async () => {
    const responseBody: ApiResponse<{ ok: boolean }> = { data: { ok: true } };
    mockFetch.mockResolvedValue(
      new Response(JSON.stringify(responseBody), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    const result = await client.get<{ ok: boolean }>("/health");

    expect(mockFetch).toHaveBeenCalledWith("https://api.example.com/health", expect.objectContaining({
      method: "GET",
      headers: expect.objectContaining({ Authorization: "Bearer token", Accept: "application/json" }),
      signal: expect.any(AbortSignal),
    }));
    expect(result).toEqual(responseBody);
  });

  it("throws ApiClientError with parsed details on failure", async () => {
    mockFetch.mockResolvedValue(
      new Response(
        JSON.stringify({
          error: { status: 401, message: "unauthorized", details: [{ code: "AUTH001", message: "invalid" }] },
          traceId: "trace-1",
        }),
        {
          status: 401,
          headers: { "content-type": "application/json" },
        }
      )
    );

    await expect(client.get("/secure"))
      .rejects.toEqual(new ApiClientError("unauthorized", 401, [{ code: "AUTH001", message: "invalid" }], "trace-1"));
  });

  it("falls back to empty data when response has no JSON body", async () => {
    mockFetch.mockResolvedValue(new Response(null, { status: 204 }));

    const result = await client.delete("/resource/1");
    expect(result.data).toBeUndefined();
  });

  it("throws timeout error when request exceeds timeout", async () => {
    jest.useFakeTimers();
    const client = new ApiClient({
      baseUrl: "https://api.example.com",
      fetchImpl: mockFetch,
      timeout: 1000
    });

    mockFetch.mockImplementation(() => new Promise((resolve) => {
      setTimeout(() => resolve(new Response("{}", { status: 200 })), 2000);
    }));

    const promise = client.get("/slow");
    jest.advanceTimersByTime(1000);

    await expect(promise).rejects.toMatchObject({
      message: "Request timeout",
      status: 408,
      details: [{ field: "_request", code: "TIMEOUT", message: "Request timed out" }],
    });

    jest.useRealTimers();
  });

  it("clears timeout on successful response", async () => {
    const clearTimeoutSpy = jest.spyOn(global, "clearTimeout");
    const responseBody: ApiResponse<{ ok: boolean }> = { data: { ok: true } };

    mockFetch.mockResolvedValue(
      new Response(JSON.stringify(responseBody), {
        status: 200,
        headers: { "content-type": "application/json" },
      })
    );

    await client.get<{ ok: boolean }>("/health");

    expect(clearTimeoutSpy).toHaveBeenCalled();
    clearTimeoutSpy.mockRestore();
  });
});
