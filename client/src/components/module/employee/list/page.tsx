"use client";

import { DynamicBreadcrumb } from "@/components/ui/breadcrumb/dynamic-breadcrumb";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { DataTable } from "@/components/ui/table/data-table";
import { useTableData } from "@/hooks/use-table-data";
import { api } from "@/lib/api";
import { Plus, Search } from "lucide-react";
import React, { useState } from "react";
import { useConfirmationDialog } from "@/context/dialog-context";
import { toast } from "sonner";
import { useMutation } from "@tanstack/react-query";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { columns } from "./column";
import FormEmployee from "./form";

const fetchEmployee = async ({
  perPage,
  page,
  search,
}: {
  perPage: number;
  page: number;
  search?: string;
}) => {
  return await api<{ data: Employee[]; paging: PageMetadata }>({
    url: "employees",
    params: {
      perPage,
      page,
      ...(search && { search }),
    },
  });
};

const fetchDeleteMutation = async (id: number): Promise<boolean> => {
  return await api({ url: `employees/${id}`, method: "DELETE" });
};

export default function EmployeePage() {
  const [isOpen, setIsOpen] = useState(false);
  const [employee, setEmployee] = useState<Employee | undefined>(undefined);
  const { showConfirmation } = useConfirmationDialog();

  const {
    data,
    pageCount,
    isLoading,
    pagination,
    searchInput,
    setPagination,
    handleSearchChange,
    refetch,
  } = useTableData<Employee>({
    queryKey: "employees",
    queryFn: fetchEmployee,
  });

  const handleEdit = (employee: Employee) => {
    setIsOpen(true);
    setEmployee(employee);
  };

  const handleClose = () => {
    setIsOpen(false);
    setEmployee(undefined);
  };

  const deleteMutation = useMutation({
    mutationFn: fetchDeleteMutation,
    onSuccess: () => {
      refetch();
      toast.success(`Employee berhasil dihapus`);
    },
    onError: (error: Error) => {
      toast.error(error ? error.message : "Operasi gagal");
    },
  });

  const handleDelete = async (id: number, name: string) => {
    const result = await showConfirmation({
      title: `Menghapus Employee ${name}`,
      description: `Apakah anda ingin menghapus Employee ${name}?`,
      confirmText: "Hapus",
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
                  placeholder="Cari Karyawan"
                  value={searchInput}
                  onChange={(e) => handleSearchChange(e.target.value)}
                  className="pl-9 h-8"
                />
              </div>
            </div>

            <Button
              size="sm"
              className="h-8"
              onClick={() => {
                setIsOpen(true);
              }}
            >
              <Plus className="h-4 w-4 mr-2" />
              Tambah Karyawan
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
              {employee ? "Ubah Karyawan" : "Tambah Karyawan"}
            </SheetTitle>
          </SheetHeader>
          <SheetDescription />
          <FormEmployee handleClose={handleClose} id={employee?.id} />
        </SheetContent>
      </Sheet>
    </div>
  );
}
