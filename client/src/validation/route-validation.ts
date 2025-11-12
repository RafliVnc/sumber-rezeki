import z from "zod";

export class RouteValidation {
  static MAIN = z.object({
    name: z.string().min(1, { message: "Nama dibutuhkan" }),
    description: z.string({ error: "Deskripsi berupa teks" }).optional(),
  });
}
