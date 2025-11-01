import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useDebounce } from "@/hooks/use-debounce";

interface UseTableDataOptions<TData, TFilters extends Record<string, any>> {
  queryKey: string;
  queryFn: (params: {
    perPage: number;
    page: number;
    search?: string;
    filters?: TFilters;
  }) => Promise<{ data: TData[]; paging: PageMetadata }>;
  initialPageSize?: number;
  debounceMs?: number;
}

export function useTableData<TData, TFilters extends Record<string, any> = {}>({
  queryKey,
  queryFn,
  initialPageSize = 10,
  debounceMs = 500,
}: UseTableDataOptions<TData, TFilters>) {
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: initialPageSize,
  });

  const [searchInput, setSearchInput] = useState("");
  const [filters, setFilters] = useState<TFilters>({} as TFilters);

  const debouncedSearch = useDebounce(searchInput, debounceMs);

  // Build query params
  const queryParams = {
    perPage: pagination.pageSize,
    page: pagination.pageIndex + 1,
    ...(debouncedSearch && { search: debouncedSearch }),
    ...(Object.keys(filters).length > 0 && { filters }),
  };

  const { data, isFetching, isError, error, refetch } = useQuery({
    queryKey: [
      queryKey,
      pagination.pageIndex,
      pagination.pageSize,
      debouncedSearch,
      filters,
    ],
    queryFn: () => queryFn(queryParams),
    placeholderData: (previousData) => previousData,
  });

  // Helper to reset to first page
  const resetToFirstPage = () => {
    setPagination((prev) => ({ ...prev, pageIndex: 0 }));
  };

  // Handlers
  const handleSearchChange = (value: string) => {
    setSearchInput(value);
    resetToFirstPage();
  };

  const handleFilterChange = (key: keyof TFilters, value: any) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    resetToFirstPage();
  };

  const handleFiltersChange = (newFilters: Partial<TFilters>) => {
    setFilters((prev) => ({ ...prev, ...newFilters }));
    resetToFirstPage();
  };

  return {
    // Data
    data: data?.data ?? [],
    pageCount: data?.paging?.totalPages ?? 0,
    totalItems: data?.paging?.totalItems ?? 0,
    isLoading: isFetching,
    isError,
    error,

    // State
    pagination,
    searchInput,
    filters,

    // Actions
    refetch,
    setPagination,
    handleSearchChange,
    handleFilterChange,
    handleFiltersChange,
  };
}
