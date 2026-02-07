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
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import FormVehicle from "./form";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useConfirmationDialog } from "@/context/dialog-context";

type VehiclePageProps = {
  queryKey: string;
  types: VehicleType[];
};

const fetchVehicles =
  (types: VehicleType[]) =>
  async ({
    perPage,
    page,
    search,
  }: {
    perPage: number;
    page: number;
    search?: string;
  }) => {
    return await api<{ data: Vehicle[]; paging: PageMetadata }>({
      url: "vehicles",
      params: {
        perPage,
        page,
        ...(search && { search }),
        types,
      },
    });
  };

export default function VehiclePage({ queryKey, types }: VehiclePageProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [vehicle, setVehicle] = useState<Vehicle>();
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
  } = useTableData<Vehicle>({
    queryKey: queryKey,
    queryFn: fetchVehicles(types),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) =>
      api({ url: `vehicles/${id}`, method: "DELETE" }),
    onSuccess: () => {
      refetch();
      toast.success("Kendaraan berhasil dihapus");
    },
  });

  const handleDelete = async (id: number, name: string) => {
    const result = await showConfirmation({
      title: `Menghapus Kendaraan ${name}`,
      description: `Apakah anda ingin menghapus kendaraan ${name}?`,
      confirmText: "Hapus",
      cancelText: "Batal",
      variant: "destructive",
    });

    if (result.confirmed) deleteMutation.mutate(id);
  };

  const handleClose = () => {
    setIsOpen(false);
    setVehicle(undefined);
  };

  return (
    <div>
      <Card>
        <CardContent>
          <div className="flex items-center justify-between mb-4">
            <div className="relative max-w-sm">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4" />
              <Input
                placeholder="Cari kendaraan"
                value={searchInput}
                onChange={(e) => handleSearchChange(e.target.value)}
                className="pl-9 h-8"
              />
            </div>

            <Button size="sm" className="h-8" onClick={() => setIsOpen(true)}>
              <Plus className="size-4" />
              Tambah Kendaraan
            </Button>
          </div>

          <DataTable
            columns={columns({
              handleEdit: (v) => {
                setVehicle(v);
                setIsOpen(true);
              },
              handleDelete,
            })}
            data={data}
            pageCount={pageCount}
            pagination={pagination}
            onPaginationChange={setPagination}
            isLoading={isLoading}
          />
        </CardContent>
      </Card>

      <Sheet open={isOpen} onOpenChange={handleClose}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>
              {vehicle ? "Ubah Kendaraan" : "Tambah Kendaraan"}
            </SheetTitle>
          </SheetHeader>
          <FormVehicle handleClose={handleClose} vehicle={vehicle} />
        </SheetContent>
      </Sheet>
    </div>
  );
}
