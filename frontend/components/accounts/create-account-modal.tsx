"use client";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { toast } from "@/components/ui/toast";

const accountSchema = z.object({
  name: z.string().min(2, "Nama akun wajib diisi"),
  type: z.string().min(1, "Tipe akun wajib dipilih"),
  initial_balance: z.coerce.number().min(0),
});

export function CreateAccountModal() {
  const queryClient = useQueryClient();
  const form = useForm({
    resolver: zodResolver(accountSchema),
    defaultValues: { name: "", type: "CASH", initial_balance: 0 },
  });

  const mutation = useMutation({
    mutationFn: (values: any) => api.post("/accounts", values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      toast.add({ title: "Sukses", description: "Akun berhasil ditambahkan" });
      form.reset();
    },
  });

  return (
    <Dialog>
      <DialogTrigger render={<Button>+ Tambah Akun</Button>} />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Tambah Akun Baru</DialogTitle>
        </DialogHeader>
        <form
          onSubmit={form.handleSubmit((v) => mutation.mutate(v))}
          className="space-y-4"
        >
          <Controller
            control={form.control}
            name="name"
            render={({ field, fieldState }) => (
              <Field invalid={fieldState.invalid}>
                <FieldLabel htmlFor="name">Nama Akun</FieldLabel>
                <Input {...field} id="name" type="input" />
                <FieldError>{fieldState.error?.message}</FieldError>
              </Field>
            )}
          />

          <Controller
            control={form.control}
            name="type"
            render={({ field, fieldState }) => (
              <Field invalid={fieldState.invalid}>
                <FieldLabel htmlFor="type">Tipe</FieldLabel>
                <Input {...field} id="type" type="input" placeholder="CASH / BANK"/>
                <FieldError>{fieldState.error?.message}</FieldError>
              </Field>
            )}
          />

          <Controller
            control={form.control}
            name="initial_balance"
            render={({ field, fieldState }) => (
              <Field invalid={fieldState.invalid}>
                <FieldLabel htmlFor="initial_balance">Saldo Awal</FieldLabel>
                <Input {...field} id="initial_balance" type="number" value={0}/>
                <FieldError>{fieldState.error?.message}</FieldError>
              </Field>
            )}
          />
          
          <Button
            type="submit"
            className="w-full"
            disabled={mutation.isPending}
          >
            Simpan
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
