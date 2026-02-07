"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Pencil, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import { ConvertToPlate } from "@/lib/utils";

export const columns = ({
  handleEdit,
  handleDelete,
}: {
  handleEdit: (value: Vehicle) => void;
  handleDelete: (id: number, name: string) => void;
}): ColumnDef<Vehicle>[] => [
  {
    accessorKey: "plate",
    header: "Plat",
    cell: ({ row }) => {
      const vehicle = row.original;

      return (
        <div className="flex items-center gap-2">
          <span>{ConvertToPlate(vehicle.plate)}</span>
        </div>
      );
    },
  },
  {
    accessorKey: "type",
    header: "Jenis",
  },
  {
    id: "actions",
    header: "Aksi",
    maxSize: 80,
    cell: ({ row }) => {
      const vehicle = row.original;

      return (
        <div className="flex gap-1 ">
          <Button
            variant="ghost"
            className="size-6"
            onClick={() => {
              handleEdit(vehicle);
            }}
          >
            <Pencil />
          </Button>
          <Button
            variant="ghost"
            className="size-6"
            onClick={() => {
              handleDelete(vehicle.id, ConvertToPlate(vehicle.plate));
            }}
          >
            <Trash2 className="text-red-500" />
          </Button>
        </div>
      );
    },
  },
];
