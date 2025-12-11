"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Pencil, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  ConvertEmployeeRole,
  ConvertEmployeeRoleBadge,
  formatIDR,
} from "@/lib/utils";
import { Badge } from "@/components/ui/badge";

export const columns = ({
  handleEdit,
  handleDelete,
}: {
  handleEdit: (value: Employee) => void;
  handleDelete: (id: number, name: string) => void;
}): ColumnDef<Employee>[] => [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "role",
    header: "Role",
    cell: ({ row }) => {
      return (
        <Badge
          variant={ConvertEmployeeRoleBadge(row.original.role) || "purple"}
        >
          {row.original.role ? ConvertEmployeeRole(row.original.role) : "-"}
        </Badge>
      );
    },
  },
  {
    accessorKey: "salary",
    header: "Gaji",
    cell: ({ row }) => {
      return <p>{formatIDR(row.original.salary) || "-"}</p>;
    },
  },
  {
    id: "actions",
    header: "Aksi",
    cell: ({ row }) => {
      const employee = row.original;

      return (
        <div className="flex gap-1 ">
          <Button
            variant="ghost"
            className="size-6"
            onClick={() => {
              handleEdit(employee);
            }}
          >
            <Pencil />
          </Button>
          <Button
            variant="ghost"
            className="size-6"
            onClick={() => {
              handleDelete(employee.id, employee.name);
            }}
          >
            <Trash2 className="text-red-500" />
          </Button>
        </div>
      );
    },
  },
];
