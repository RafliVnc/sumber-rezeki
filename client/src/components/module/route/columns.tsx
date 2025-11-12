"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Pencil, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";

export const columns = ({
  handleEdit,
  handleDelete,
}: {
  handleEdit: (value: Route) => void;
  handleDelete: (id: number, name: string) => void;
}): ColumnDef<Route>[] => [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "description",
    header: "Deskripsi",
    cell: ({ row }) => {
      return <p className="w-[200px]">{row.original.description || "-"}</p>;
    },
  },
  {
    id: "actions",
    header: "Aksi",
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
