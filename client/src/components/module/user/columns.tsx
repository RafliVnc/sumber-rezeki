"use client";

import { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal, Pencil, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ConvertUserRole } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { UserRole } from "@/type/enum/user-role";
import { DataTableColumnHeader } from "@/components/ui/table/data-table-column-header";

export const columns = ({
  handleEdit,
  handleDelete,
}: {
  handleEdit: (value: User) => void;
  handleDelete: (id: string) => void;
}): ColumnDef<User>[] => [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "username",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="username" />
    ),
  },
  {
    accessorKey: "phone",
    header: "No. HP",
  },
  {
    accessorKey: "role",
    header: "Role",
    cell: ({ row }) => {
      return (
        <Badge
          variant={
            row.original.role === UserRole.OWNER
              ? "blue"
              : row.original.role === UserRole.WAREHOUSE_HEAD
                ? "green"
                : "purple"
          }
        >
          {ConvertUserRole(row.original.role)}
        </Badge>
      );
    },
  },
  {
    id: "actions",
    header: "Aksi",
    maxSize: 80,
    cell: ({ row }) => {
      const user = row.original;

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Aksi</DropdownMenuLabel>
            <DropdownMenuItem
              onClick={() => {
                handleEdit(user);
              }}
            >
              <Pencil className="size-4" />
              Ubah Pengguna
            </DropdownMenuItem>
            <DropdownMenuItem
              variant="destructive"
              onClick={() => {
                handleDelete(user.id);
              }}
            >
              <Trash2 className="size-4" /> Hapus Pengguna
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];
