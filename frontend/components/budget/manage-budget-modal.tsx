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

const budgetSchema = z.object({
  category_id: z.string().min(1, "Pilih kategori"),
  amount: z.coerce.number().min(1, "Budget minimal 1"),
});

export function ManageBudgetModal() {
  const queryClient = useQueryClient();
  const form = useForm({
    resolver: zodResolver(budgetSchema),
    defaultValues: { amount: 0 },
  });

  const { data: categories } = useQuery({
    queryKey: ["categories"],
    queryFn: async () => (await api.get("/categories")).data,
  });

  const mutation = useMutation({
    mutationFn: (values: any) => api.post("/budgets", values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["budgets"] });
      toast.add({
        title: "Budget Diupdate",
        description: "Limit budget berhasil disimpan",
      });
      form.reset();
    },
  });

  return (
    <Dialog>
      <DialogTrigger render={<Button variant="outline">+ Set Budget</Button>} />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Atur Budget Bulanan</DialogTitle>
        </DialogHeader>

        <form
          onSubmit={form.handleSubmit((v) => mutation.mutate(v))}
          className="space-y-4"
        >
          <Controller
            control={form.control}
            name="category_id"
            render={({ field, fieldState }) => (
              <Field invalid={fieldState.invalid}>
                <FieldLabel>Kategori</FieldLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Pilih Kategory" />
                  </SelectTrigger>
                  <SelectContent>
                    {categories?.map((a: any) => (
                      <SelectItem key={a.id} value={a.id}>
                        {a.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FieldError>{fieldState.error?.message}</FieldError>
              </Field>
            )}
          />

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
                  value={field.value ?? ""}
                  onChange={(e) => {
                    const val = e.target.value;
                    // Kirim null atau undefined jika kosong agar Zod bisa menangkap error .min(1)
                    field.onChange(val === "" ? "" : val);
                  }}
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
            Simpan Budget
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
