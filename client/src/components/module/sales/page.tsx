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
import { columns } from "./columns";
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
import FormSales from "./form";

const fetchSales = async ({
  perPage,
  page,
  search,
}: {
  perPage: number;
  page: number;
  search?: string;
}) => {
  return await api<{ data: Sales[]; paging: PageMetadata }>({
    url: "sales",
    params: {
      perPage,
      page,
      ...(search && { search }),
    },
  });
};

const fetchDeleteMutation = async (id: number): Promise<boolean> => {
  return await api({ url: `sales/${id}`, method: "DELETE" });
};

export default function SalesPage() {
  const [isOpen, setIsOpen] = useState(false);
  const [sales, setSales] = useState<Sales | undefined>(undefined);
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
  } = useTableData<Sales>({
    queryKey: "saless",
    queryFn: fetchSales,
  });

  const handleEdit = (sales: Sales) => {
    setIsOpen(true);
    setSales(sales);
  };

  const handleClose = () => {
    setIsOpen(false);
    setSales(undefined);
  };

  const deleteMutation = useMutation({
    mutationFn: fetchDeleteMutation,
    onSuccess: () => {
      refetch();
      toast.success(`Sales berhasil dihapus`);
    },
    onError: (error: Error) => {
      toast.error(error ? error.message : "Operasi gagal");
    },
  });

  const handleDelete = async (id: number) => {
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
      <h1 className="text-2xl font-bold mb-1">Sales</h1>
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
                  placeholder="Cari sales"
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
              Tambah Sales
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
            <SheetTitle>{sales ? "Ubah Sales" : "Tambah Sales"}</SheetTitle>
          </SheetHeader>
          <SheetDescription />
          <FormSales handleClose={handleClose} sales={sales} />
        </SheetContent>
      </Sheet>
    </div>
  );
}
