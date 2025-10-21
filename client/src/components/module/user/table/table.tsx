"use client";

import { DataTable } from "@/components/ui/table/data-table";
import { columns } from "./columns";
import { api } from "@/lib/api";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search, Plus, X } from "lucide-react";
import { DataTableFacetedFilter } from "@/components/ui/table/data-table-faceted-filter";
import { UserRole } from "@/type/enum/user-role";
import { useTableData } from "@/hooks/use-table-data";

// Options untuk role filter
const roleOptions = [
  { label: "Owner", value: UserRole.OWNER },
  { label: "Bendahara", value: UserRole.TREASURER },
  { label: "Kepala Gudang", value: UserRole.WAREHOUSE_HEAD },
];

interface UserFilters {
  roles?: string[];
  status?: string;
  department?: string;
}

const fetchUsers = async ({ perPage, page, search, filters }: any) => {
  return await api<{ data: User[]; paging: PageMetadata }>({
    url: "users",
    params: {
      perPage,
      page,
      ...(search && { search }),
      ...(filters?.roles && { roles: filters.roles }),
      ...(filters?.status && { status: filters.status }),
    },
  });
};

export default function TableUser() {
  const {
    data,
    pageCount,
    isLoading,
    pagination,
    searchInput,
    filters,
    setPagination,
    handleSearchChange,
    handleFilterChange,
  } = useTableData<User, UserFilters>({
    queryKey: "users",
    queryFn: fetchUsers,
  });

  // Convert Set to array untuk role filter
  const selectedRoles = new Set(filters.roles ?? []);

  const handleRoleFilterChange = (values: Set<string>) => {
    handleFilterChange(
      "roles",
      values.size > 0 ? Array.from(values) : undefined
    );
  };

  return (
    <div className="container mx-auto">
      {/* Toolbar */}
      <div className="flex items-center justify-between gap-4 mb-4">
        <div className="flex items-center gap-2 flex-1">
          <div className="relative max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search users..."
              value={searchInput}
              onChange={(e) => handleSearchChange(e.target.value)}
              className="pl-9 h-8"
            />
          </div>

          <DataTableFacetedFilter
            title="Role"
            options={roleOptions}
            selectedValues={selectedRoles}
            onSelectedValuesChange={handleRoleFilterChange}
          />
        </div>

        <Button size="sm" className="h-8">
          <Plus className="h-4 w-4 mr-2" />
          Tambah Pengguna
        </Button>
      </div>

      {/* Table */}
      <DataTable
        columns={columns}
        data={data}
        pageCount={pageCount}
        pagination={pagination}
        onPaginationChange={setPagination}
        isLoading={isLoading}
      />
    </div>
  );
}
