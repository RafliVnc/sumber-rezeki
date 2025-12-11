import { EmployeeRole } from "@/type/enum/employee-role";
import z from "zod";

export class EmployeeValidation {
  static MAIN = z
    .object({
      name: z.string().min(1, { message: "Nama dibutuhkan" }),
      salary: z.number().min(1, { message: "Gaji dibutuhkan" }),
      role: z.enum(EmployeeRole, { message: "Role dibutuhkan" }),
      supervisorId: z
        .number()
        .min(1, { message: "Supervisor dibutuhkan" })
        .optional(),
      phone: z
        .string()
        .refine((val) => /^\d+$/.test(val), {
          message: "No. HP harus angka numerik",
        })
        .refine((val) => val.length >= 11, {
          message: "No. HP minimal 11 karakter",
        })
        .optional(),
      routeIds: z
        .array(
          z
            .number({ message: "Rute dibutuhkan" })
            .min(1, { message: "Rute dibutuhkan" })
        )
        .optional(),
    })
    .refine(
      (data) => {
        // Helper dan Driver harus ada supervisorId
        if (
          data.role === EmployeeRole.HELPER ||
          data.role === EmployeeRole.DRIVER
        ) {
          return data.supervisorId !== undefined && data.supervisorId > 0;
        }
        return true;
      },
      {
        message: "Sales dibutuhkan",
        path: ["supervisorId"],
      }
    )
    .refine(
      (data) => {
        // Sales harus ada phone
        if (data.role === EmployeeRole.SALES) {
          return data.phone !== undefined && data.phone.length >= 11;
        }
        return true;
      },
      {
        message: "No. HP dibutuhkan untuk Sales",
        path: ["phone"],
      }
    )
    .refine(
      (data) => {
        // Sales harus ada routeIds dan tidak kosong
        if (data.role === EmployeeRole.SALES) {
          return data.routeIds !== undefined && data.routeIds.length > 0;
        }
        return true;
      },
      {
        message: "Rute dibutuhkan untuk Sales",
        path: ["routeIds"],
      }
    );
}
