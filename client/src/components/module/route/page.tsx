"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { DataTable } from "@/components/ui/table/data-table";
import { useTableData } from "@/hooks/use-table-data";
import { api } from "@/lib/api";
import { Plus, Search } from "lucide-react";
import React, { useState } from "react";
import { columns } from "./columns";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import FormRoute from "./form";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useConfirmationDialog } from "@/context/dialog-context";

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

const fetchDeleteMutation = async (id: number): Promise<boolean> => {
  return await api({ url: `routes/${id}`, method: "DELETE" });
};

export default function RoutePage() {
  const [isOpen, setIsOpen] = useState(false);
  const [route, setRoute] = useState<Route | undefined>(undefined);
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
  } = useTableData<Route>({
    queryKey: "routes",
    queryFn: fetchRoutes,
  });

  const handleEdit = (routes: Route) => {
    setIsOpen(true);
    setRoute(routes);
  };

  const deleteMutation = useMutation({
    mutationFn: fetchDeleteMutation,
    onSuccess: () => {
      refetch();
      toast.success(`Rute berhasil dihapus`);
    },
    onError: (error: Error) => {
      toast.error(error ? error.message : "Operasi gagal");
    },
  });

  const handleDelete = async (id: number, name: string) => {
    const result = await showConfirmation({
      title: `Menghapus Rute ${name}`,
      description: `Apakah anda ingin menghapus rute ${name}?`,
      confirmText: "Hapus",
      cancelText: "Batal",
      variant: "destructive",
    });

    if (result.confirmed) {
      deleteMutation.mutate(id);
    }
  };

  const handleClose = () => {
    setIsOpen(false);
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
                  placeholder="Cari rute"
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
              <Plus className="size-4" />
              Tambah Rute
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
            <SheetTitle>{route ? "Ubah Rute" : "Tambah Rute"}</SheetTitle>
          </SheetHeader>
          <SheetDescription />
          <FormRoute handleClose={handleClose} route={route} />
        </SheetContent>
      </Sheet>
    </div>
  );
}
