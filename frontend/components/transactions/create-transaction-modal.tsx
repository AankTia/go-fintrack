"use client";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { toast } from "@/components/ui/toast";

const txSchema = z.object({
  account_id: z.string().min(1, "Pilih akun"),
  category_id: z.string().min(1, "Pilih kategori"),
  amount: z.coerce.number().min(1),
  type: z.enum(["income", "expense"]),
  description: z.string().optional(),
  date: z.string().min(1),
});

export function CreateTransactionModal() {
  const queryClient = useQueryClient();
  const form = useForm({
    resolver: zodResolver(txSchema),
    defaultValues: {
      type: "expense",
      date: new Date().toISOString().split("T")[0],
    },
  });

  // Fetch data untuk dropdown
  const { data: accounts } = useQuery({
    queryKey: ["accounts"],
    queryFn: async () => (await api.get("/accounts")).data,
  });
  const { data: categories } = useQuery({
    queryKey: ["categories"],
    queryFn: async () => (await api.get("/categories")).data,
  });

  const mutation = useMutation({
    mutationFn: (values: any) => api.post("/transactions", values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
      queryClient.invalidateQueries({ queryKey: ["accounts"] }); // Update saldo akun
      toast.add({ title: "Sukses", description: "Transaksi dicatat" });
      form.reset();
    },
  });

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="destructive">+ Catat Transaksi</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Tambah Transaksi</DialogTitle>
        </DialogHeader>
        <form
          onSubmit={form.handleSubmit((v) => mutation.mutate(v))}
          className="space-y-4"
        >
          {/* <FormField
            control={form.control}
            name="type"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Tipe</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="income">Pemasukan</SelectItem>
                    <SelectItem value="expense">Pengeluaran</SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          /> */}

          {/* <FormField
            control={form.control}
            name="account_id"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Akun</FormLabel>
                <Select onValueChange={field.onChange}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Pilih Akun" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {accounts?.map((a: any) => (
                      <SelectItem key={a.id} value={a.id}>
                        {a.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          /> */}

          {/* <FormField
            control={form.control}
            name="category_id"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Kategori</FormLabel>
                <Select onValueChange={field.onChange}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Pilih Kategori" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {categories?.map((c: any) => (
                      <SelectItem key={c.id} value={c.id}>
                        {c.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          /> */}

          <Controller
            control={form.control}
            name="amount"
            render={({ field, fieldState }) => (
              <Field invalid={fieldState.invalid}>
                <FieldLabel htmlFor="amount">Nominal</FieldLabel>
                <Input
                  {...field}
                  id="amount"
                  type="number"
                  value={0}
                />
                <FieldError>{fieldState.error?.message}</FieldError>
              </Field>
            )}
          />

          <Controller
            control={form.control}
            name="description"
            render={({ field, fieldState }) => (
              <Field invalid={fieldState.invalid}>
                <FieldLabel htmlFor="description">Keterangan</FieldLabel>
                <Input
                  {...field}
                  id="description"
                />
                <FieldError>{fieldState.error?.message}</FieldError>
              </Field>
            )}
          />

          <Button
            type="submit"
            className="w-full"
            disabled={mutation.isPending}
          >
            Simpan Transaksi
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
