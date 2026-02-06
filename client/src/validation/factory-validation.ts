import z from "zod";

export class FactoryValidation {
  static MAIN = z.object({
    name: z.string().min(1, { message: "Nama kilang dibutuhkan" }),
    phone: z
      .string()
      .refine((val) => /^\d+$/.test(val), {
        message: "No. HP harus angka numerik",
      })
      .refine((val) => val.length >= 11, {
        message: "No. HP minimal 11 karakter",
      }),
    dueDate: z
      .number({
        message: "Jatuh tempo dibutuhkan",
      })
      .min(1, { message: "Jatuh tempo dibutuhkan" }),
    description: z.string().optional(),
  });
}
