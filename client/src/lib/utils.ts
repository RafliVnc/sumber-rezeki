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
