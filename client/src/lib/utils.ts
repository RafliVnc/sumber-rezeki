import { EmployeeRole } from "@/type/enum/employee-role";
import { UserRole } from "@/type/enum/user-role";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const ConvertUserRole = (role: UserRole) => {
  switch (role) {
    case UserRole.SUPER_ADMIN:
      return "Super Admin";
    case UserRole.OWNER:
      return "Owner";
    case UserRole.WAREHOUSE_HEAD:
      return "Kepala Gudang";
    case UserRole.TREASURER:
      return "Bendahara";
  }
};

export const ConvertEmployeeRole = (role: EmployeeRole) => {
  switch (role) {
    case EmployeeRole.EMPLOYEE_WAREHOUSE_HEAD:
      return "Kepala Gudang";
    case EmployeeRole.SALES:
      return "Sales";
    case EmployeeRole.DRIVER:
      return "Supir";
    case EmployeeRole.HELPER:
      return "Kenek";
    case EmployeeRole.EMPLOYEE_TREASURER:
      return "Bendahara";
    case EmployeeRole.STAFF:
      return "Orang Gudang";
  }
};

export const ConvertEmployeeRoleBadge = (
  role: EmployeeRole
):
  | "default"
  | "blue"
  | "purple"
  | "green"
  | "sales"
  | "driver"
  | "helper"
  | "staff"
  | "secondary"
  | "destructive"
  | "outline"
  | null
  | undefined => {
  switch (role) {
    case EmployeeRole.EMPLOYEE_WAREHOUSE_HEAD:
      return "green";
    case EmployeeRole.SALES:
      return "sales";
    case EmployeeRole.DRIVER:
      return "driver";
    case EmployeeRole.HELPER:
      return "helper";
    case EmployeeRole.EMPLOYEE_TREASURER:
      return "purple";
    case EmployeeRole.STAFF:
      return "staff";
  }
};

export const formatIDR = (value: number | string) => {
  const num =
    typeof value === "string" ? Number(value.replace(/\D/g, "")) : value;
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(num || 0);
};

export const idrToNumber = (value: string): number => {
  if (!value) return 0;

  // Hapus "Rp", spasi, titik
  const cleaned = value
    .replace(/Rp/gi, "")
    .replace(/\s+/g, "")
    .replace(/\./g, "")
    .replace(/,/g, ".");

  const num = parseFloat(cleaned);
  return isNaN(num) ? 0 : num;
};
