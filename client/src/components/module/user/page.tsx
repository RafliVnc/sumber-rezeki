"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search, Plus } from "lucide-react";
import { DataTableFacetedFilter } from "@/components/ui/table/data-table-faceted-filter";
import { DataTable } from "@/components/ui/table/data-table";
import { columns } from "./columns";
import { api } from "@/lib/api";
import { UserRole } from "@/type/enum/user-role";
import { useTableData } from "@/hooks/use-table-data";
import { useState } from "react";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import FormUser from "./form";
import { useConfirmationDialog } from "@/context/dialog-context";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

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

const fetchUsers = async ({
  perPage,
  page,
  search,
  filters,
}: {
  perPage: number;
  page: number;
  search?: string;
  filters?: UserFilters;
}) => {
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

const fetchDeleteMutation = async (id: string): Promise<boolean> => {
  return await api({ url: `users/${id}`, method: "DELETE" });
};

export default function UserPage() {
  const [isOpen, setIsOpen] = useState(false);
  const [user, setUser] = useState<User | undefined>(undefined);
  const { showConfirmation } = useConfirmationDialog();

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
    refetch,
  } = useTableData<User, UserFilters>({
    queryKey: "users",
    queryFn: fetchUsers,
  });

  // Convert Set to array for role filter
  const selectedRoles = new Set(filters.roles ?? []);

  const handleRoleFilterChange = (values: Set<string>) => {
    handleFilterChange(
      "roles",
      values.size > 0 ? Array.from(values) : undefined,
    );
  };

  const handleEdit = (user: User) => {
    setIsOpen(true);
    setUser(user);
  };

  const handleClose = () => {
    setIsOpen(false);
    setUser(undefined);
  };

  const deleteMutation = useMutation({
    mutationFn: fetchDeleteMutation,
    onSuccess: () => {
      refetch();
      toast.success(`Pengguna berhasil dihapus`);
    },
    onError: (error: Error) => {
      toast.error(error ? error.message : "Operasi gagal");
    },
  });

  const handleDelete = async (id: string) => {
    const result = await showConfirmation({
      title: "Apakah anda yakin?",
      description: "Aksi ini tidak dapat dibatalkan",
      confirmText: "Ya",
      cancelText: "Batal",
      variant: "destructive",
    });

    if (result.confirmed) {
      deleteMutation.mutate(id);
    }
  };

  return (
    <div>
      {/* Table card */}
      <Card>
        <CardContent>
          {/* Toolbar */}
          <div className="flex items-center justify-between gap-4 mb-4">
            <div className="flex items-center gap-2 flex-1">
              <div className="relative max-w-sm">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Cari pengguna"
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

            <Button
              size="sm"
              className="h-8"
              onClick={() => {
                setIsOpen(true);
              }}
            >
              <Plus className="size-4" />
              Tambah Pengguna
            </Button>
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

      {/* Sheet */}
      <Sheet open={isOpen} onOpenChange={handleClose}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>
              {user ? "Ubah Pengguna" : "Tambah Pengguna"}
            </SheetTitle>
          </SheetHeader>
          <SheetDescription />
          <FormUser handleClose={handleClose} user={user} />
        </SheetContent>
      </Sheet>
    </div>
  );
}
