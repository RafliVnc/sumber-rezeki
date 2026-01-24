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

export const ConvertStatusLabel = (status: AttendanceStatus | null): string => {
  const statusLabels: Record<AttendanceStatus, string> = {
    [AttendanceStatus.PRESENT]: "Hadir",
    [AttendanceStatus.LEAVE]: "Izin",
    [AttendanceStatus.SICK]: "Sakit",
    [AttendanceStatus.ABSENT]: "Tidak Hadir",
  };

  return status ? statusLabels[status] : "-";
};

export const ConvertEmployeeRoleBadge = (
  role: EmployeeRole,
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

export const getDayName = (date: Date): string => {
  const days = ["Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"];
  return days[date.getDay()];
};

export const getSundayOfWeek = (date: Date): Date => {
  const d = new Date(date);
  const day = d.getDay();
  const diff = d.getDate() - day;
  return new Date(d.setDate(diff));
};

export const formatDate = (date: Date): string => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
};
