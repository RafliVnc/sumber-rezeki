import z from "zod";

export class SalesValidation {
  static MAIN = z.object({
    name: z.string().min(1, { message: "Nama dibutuhkan" }),
    phone: z
      .string()
      .refine((val) => /^\d+$/.test(val), {
        message: "No. HP harus angka numerik",
      })
      .refine((val) => val.length >= 11, {
        message: "No. HP minimal 11 karakter",
      }),
    routeIds: z
      .array(
        z
          .number({ error: "Route dibutuhkan" })
          .min(1, { message: "Route dibutuhkan" })
      )
      .optional(),
  });
}
