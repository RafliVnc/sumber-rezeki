"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Pencil, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";

export const columns = ({
  handleEdit,
  handleDelete,
}: {
  handleEdit: (value: Factory) => void;
  handleDelete: (id: number, name: string) => void;
}): ColumnDef<Factory>[] => [
  {
    accessorKey: "name",
    header: "Nama",
  },
  {
    accessorKey: "phone",
    header: "No. HP",
  },
  {
    accessorKey: "description",
    header: "Deskripsi",
    cell: ({ row }) => {
      return (
        <div className="w-full">
          <p className="line-clamp-3 !whitespace-normal">
            {row.original.description || "-"}
          </p>
        </div>
      );
    },
  },
  {
    id: "actions",
    header: "Aksi",
    maxSize: 80,
    cell: ({ row }) => {
      const routes = row.original;

      return (
        <div className="flex gap-1 ">
          <Button
            variant="ghost"
            className="size-6"
            onClick={() => {
              handleEdit(routes);
            }}
          >
            <Pencil />
          </Button>
          <Button
            variant="ghost"
            className="size-6"
            onClick={() => {
              handleDelete(routes.id, routes.name);
            }}
          >
            <Trash2 className="text-red-500" />
          </Button>
        </div>
      );
    },
  },
];
