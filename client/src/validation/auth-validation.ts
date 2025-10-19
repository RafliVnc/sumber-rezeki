import { z } from "zod";

export class AuthValidation {
  static MAIN = z.object({
    username: z.string().min(1, { error: "Username dibutuhkan" }),
    password: z.string().min(1, { error: "Kata sandi dibutuhkan" }),
  });
}
