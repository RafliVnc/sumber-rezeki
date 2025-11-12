"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Pencil, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";

export const columns = ({
  handleEdit,
  handleDelete,
}: {
  handleEdit: (value: Sales) => void;
  handleDelete: (id: number, name: string) => void;
}): ColumnDef<Sales>[] => [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "phone",
    header: "No. HP",
  },
  {
    id: "actions",
    header: "Aksi",
    cell: ({ row }) => {
      const sales = row.original;

      return (
        <div className="flex gap-1 ">
          <Button
            variant="ghost"
            className="size-6"
            onClick={() => {
              handleEdit(sales);
            }}
          >
            <Pencil />
          </Button>
          <Button
            variant="ghost"
            className="size-6"
            onClick={() => {
              handleDelete(sales.id, sales.name);
            }}
          >
            <Trash2 className="text-red-500" />
          </Button>
        </div>
      );
    },
  },
];
