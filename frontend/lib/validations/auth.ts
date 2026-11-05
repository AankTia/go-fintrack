import * as z from "zod";

export const registerSchema = z.object({
  email: z.email("email tidak valid"),
  password: z.string().min(6, "Password minimal 6 karakter"),
  full_name: z.string().min(2, "Nama lengkap harus diisi"),
});

export const loginSchema = z.object({
  email: z.email("Email tidak valid"),
  password: z.string().min(1, "Password wajib diisi"),
});
