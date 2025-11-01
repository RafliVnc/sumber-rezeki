import { UserRole } from "@/type/enum/user-role";
import z from "zod";

export class UserValidation {
  static BASE = z.object({
    name: z.string().min(1, { message: "Nama dibutuhkan" }),
    username: z.string().min(1, { message: "Username dibutuhkan" }),
    phone: z
      .string()
      .refine((val) => /^\d+$/.test(val), {
        message: "No. HP harus angka numerik",
      })
      .refine((val) => val.length >= 11, {
        message: "No. HP minimal 11 karakter",
      }),
    role: z.enum(UserRole, { message: "Role dibutuhkan" }),
  });

  static CREATE = this.BASE.extend({
    password: z.string().min(8, { message: "Kata sandi minimal 8 karakter" }),
  });

  static UPDATE = this.BASE.partial();
}
