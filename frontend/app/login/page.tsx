"use client";
import * as z from "zod";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema } from "@/lib/validations/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { toast } from "@/components/ui/toast";
import api from "@/lib/api";
import { useRouter } from "next/navigation";
import axios from "axios"

type LoginFormValues = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

  async function onSubmit(values: LoginFormValues) {
    try {
      const response = await api.post("/auth/login", values);
      localStorage.setItem("token", response.data.token);
      toast.add({ title: "Welcome back!" });
      window.location.href = "/dashboard";
    } catch (error: unknown) {
      let message = "Kredensial salah";

      if (axios.isAxiosError(error)) {
        message = error.response?.data?.message || message;
      }

      toast.add({
        type: "error",
        title: "Login Gagal",
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
        <h1 className="text-2xl font-bold text-center">Login</h1>

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
          Login
        </Button>

        <p className="text-center text-sm">
          Belum punya akun?{" "}
          <a href="/register" className="text-blue-500">
            Register
          </a>
        </p>
      </form>
    </div>
  );
}
