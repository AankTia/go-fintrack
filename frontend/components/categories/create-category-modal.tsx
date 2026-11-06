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

const categorySchema = z.object({
  name: z.string().min(2, "Nama kategori wajib diisi"),
  type: z.string().min(1, "Tipe kategori wajib dipilih"),
});

export function CreateCategoryModal() {
  const queryClient = useQueryClient();
  const form = useForm({
    resolver: zodResolver(categorySchema),
    defaultValues: { name: "", type: "EXPENSE" },
  });

  const mutation = useMutation({
    mutationFn: (values: any) => api.post("/categories", values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["categories"] });
      toast.add({ title: "Sukses", description: "Kategori berhasil ditambahkan" });
      form.reset();
    },
  });

  return (
    <Dialog>
      <DialogTrigger render={<Button>+ Tambah Kategori</Button>} />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Tambah Kategori Baru</DialogTitle>
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
                <FieldLabel htmlFor="name">Nama Kategori</FieldLabel>
                <Input {...field} id="name" type="input" placeholder="misal: Makanan, Transportasi" />
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
                <Input {...field} id="type" type="input" placeholder="INCOME / EXPENSE"/>
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
