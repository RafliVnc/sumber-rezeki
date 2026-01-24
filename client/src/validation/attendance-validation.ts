import { AttendanceStatus } from "@/type/enum/attendance-status";
import z from "zod";

export class AttendanceValidation {
  // Validation untuk single date attendance
  static DATE_ATTENDANCE = (totalEmployees: number) =>
    z.discriminatedUnion("action", [
      z.object({
        action: z.literal("update"),
        date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/),
        employees: z
          .array(
            z.object({
              id: z.number().positive(),
              status: z.enum(AttendanceStatus),
            }),
          )
          .length(totalEmployees),
      }),

      z.object({
        action: z.literal("delete"),
        date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/),
        employees: z.array(z.any()).length(0),
      }),
    ]);

  // Validation untuk batch update (bisa 1 hari atau lebih)
  static BATCH = (totalEmployees: number) =>
    z.object({
      attendances: z
        .array(AttendanceValidation.DATE_ATTENDANCE(totalEmployees))
        .min(1, "Minimal harus ada satu tanggal"),
    });

  // Form schema untuk UI - fleksibel, tidak wajib lengkap
  static UI_FORM = z.object({
    dates: z.record(
      z.string(), // date string
      z.object({
        isActive: z.boolean(), // apakah tanggal ini aktif di form
        employees: z.array(
          z.object({
            id: z.number(),
            name: z.string(),
            status: z.nativeEnum(AttendanceStatus).nullable(),
          }),
        ),
      }),
    ),
  });
}

// Type inference untuk TypeScript
export type DateAttendanceInput = z.infer<
  ReturnType<typeof AttendanceValidation.DATE_ATTENDANCE>
>;
export type BatchAttendanceInput = z.infer<
  ReturnType<typeof AttendanceValidation.BATCH>
>;
export type UIFormInput = z.infer<typeof AttendanceValidation.UI_FORM>;
