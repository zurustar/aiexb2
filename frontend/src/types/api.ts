export type Pagination = {
  page: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
};

export type ApiResponse<T> = {
  data: T;
  meta?: Pagination;
  message?: string;
};

export type ApiErrorDetail = {
  field?: string;
  code: string;
  message: string;
};

export type ApiErrorResponse = {
  error: {
    status: number;
    message: string;
    details?: ApiErrorDetail[];
  };
  traceId?: string;
};

export type PaginatedResponse<T> = ApiResponse<T[]> & { meta: Pagination };
