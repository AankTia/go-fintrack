"use client";
import * as z from "zod";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { registerSchema } from "@/lib/validations/auth";
import { Button } from "@/components/ui/button";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { toast } from "@/components/ui/toast";
import api from "@/lib/api";
import { useRouter } from "next/navigation";
import axios from "axios";

type RegisterFormValues = z.infer<typeof registerSchema>;

export default function RegisterPage() {
  const router = useRouter();
  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { email: "", password: "", full_name: "" },
  });

  async function onSubmit(values: RegisterFormValues) {
    try {
      await api.post("/auth/register", values);
      toast.add({
        title: "Berhasil",
        description: "Akun berhasil dibuat, silahkan login.",
      });
      router.push("/login");
    } catch (error: unknown) {
      let message = "Terjadi kesalahan";

      if (axios.isAxiosError(error)) {
        message = error.response?.data?.message || message;
      }
      toast.add({
        type: "error",
        title: "Error",
        description: message,
      });
    }
  }

  return (
    <div className="flex justify-center items-center min-h-screen">
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="w-full max-w-md space-y-4 p-6 border rounded-lg shadow"
      >
        <h1 className="text-2xl font-bold text-center">Register</h1>

        <Controller
          control={form.control}
          name="full_name"
          render={({ field, fieldState }) => (
            <Field invalid={fieldState.invalid}>
              <FieldLabel htmlFor="full_name">Full Name</FieldLabel>
              <Input {...field} id="full_name" type="input" />
              <FieldError>{fieldState.error?.message}</FieldError>
            </Field>
          )}
        />

        <Controller
          control={form.control}
          name="email"
          render={({ field, fieldState }) => (
            <Field invalid={fieldState.invalid}>
              <FieldLabel htmlFor="email">Email</FieldLabel>
              <Input {...field} id="email" type="email" />
              <FieldError>{fieldState.error?.message}</FieldError>
            </Field>
          )}
        />

        <Controller
          control={form.control}
          name="password"
          render={({ field, fieldState }) => (
            <Field invalid={fieldState.invalid}>
              <FieldLabel htmlFor="password">Password</FieldLabel>
              <Input {...field} id="password" type="password" />
              <FieldError>{fieldState.error?.message}</FieldError>
            </Field>
          )}
        />

        <Button type="submit" className="w-full">
          Daftar
        </Button>
        <p className="text-center text-sm">
          Sudah punya akun?{" "}
          <a href="/login" className="text-blue-500">
            Login
          </a>
        </p>
      </form>
    </div>
  );
}
