import { VehicleType } from "@/type/enum/vehicle-type";
import z from "zod";

export class VehicleValidation {
  static MAIN = z.object({
    plate: z.string().regex(/^[A-Z]{1,2}\s\d{1,4}\s[A-Z]{1,3}$/, {
      message: "Format plat harus seperti BM 2937 BA",
    }),
    type: z.enum(VehicleType, {
      message: "Tipe kendaraan harus PICKUP, TRONTON, atau TRUCK",
    }),
  });
}
