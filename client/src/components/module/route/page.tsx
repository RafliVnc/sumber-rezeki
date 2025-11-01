"use client";

import { DynamicBreadcrumb } from "@/components/ui/breadcrumb/dynamic-breadcrumb";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { DataTable } from "@/components/ui/table/data-table";
import { useTableData } from "@/hooks/use-table-data";
import { api } from "@/lib/api";
import { Search } from "lucide-react";
import React from "react";
import { columns } from "./columns";

const fetchRoutes = async ({
  perPage,
  page,
  search,
}: {
  perPage: number;
  page: number;
  search?: string;
}) => {
  return await api<{ data: Route[]; paging: PageMetadata }>({
    url: "routes",
    params: {
      perPage,
      page,
      ...(search && { search }),
    },
  });
};

export default function RoutePage() {
  const {
    data,
    pageCount,
    isLoading,
    pagination,
    searchInput,
    setPagination,
    handleSearchChange,
    refetch,
  } = useTableData<Route>({
    queryKey: "routes",
    queryFn: fetchRoutes,
  });

  //   TODO: Handle edit and delete
  const handleEdit = (routes: Route) => {
    console.log(routes);
  };

  const handleDelete = async (id: number) => {
    console.log(id);
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-1">Rute</h1>
      <DynamicBreadcrumb />

      {/* Table card */}
      <Card className="mt-4">
        <CardContent>
          {/* Toolbar */}
          <div className="flex items-center justify-between gap-4 mb-4">
            <div className="flex items-center gap-2 flex-1">
              <div className="relative max-w-sm">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Cari rute"
                  value={searchInput}
                  onChange={(e) => handleSearchChange(e.target.value)}
                  className="pl-9 h-8"
                />
              </div>
            </div>

            {/* <Button
              size="sm"
              className="h-8"
              onClick={() => {
                setIsOpen(true);
              }}
            >
              <Plus className="h-4 w-4 mr-2" />
              Tambah Rute
            </Button> */}
          </div>

          <DataTable
            columns={columns({ handleEdit, handleDelete })}
            data={data}
            pageCount={pageCount}
            pagination={pagination}
            onPaginationChange={setPagination}
            isLoading={isLoading}
          />
        </CardContent>
      </Card>
    </div>
  );
}
