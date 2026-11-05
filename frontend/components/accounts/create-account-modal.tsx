"use client";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useToast } from "@/components/ui/use-toast";

const accountSchema = z.object({
  name: z.string().min(2, "Nama akun wajib diisi"),
  type: z.string().min(1, "Tipe akun wajib dipilih"),
  initial_balance: z.coerce.number().min(0),
});

export function CreateAccountModal() {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const form = useForm({ resolver: zodResolver(accountSchema), defaultValues: { name: "", type: "CASH", initial_balance: 0 } });

  const mutation = useMutation({
    mutationFn: (values: any) => api.post("/accounts", values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      toast({ title: "Sukses", description: "Akun berhasil ditambahkan" });
      form.reset();
    },
  });

  return (
    <Dialog>
      <DialogTrigger asChild><Button>+ Tambah Akun</Button></DialogTrigger>
      <DialogContent>
        <DialogHeader><DialogTitle>Tambah Akun Baru</DialogTitle></DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit((v) => mutation.mutate(v))} className="space-y-4">
            <FormField control={form.control} name="name" render={({ field }) => (
              <FormItem><FormLabel>Nama Akun</FormLabel><FormControl><Input {...field} placeholder="Contoh: Bank BCA" /></FormControl><FormMessage /></FormItem>
            )} />
            <FormField control={form.control} name="type" render={({ field }) => (
              <FormItem><FormLabel>Tipe</FormLabel><FormControl><Input {...field} placeholder="CASH / BANK" /></FormControl><FormMessage /></FormItem>
            )} />
            <FormField control={form.control} name="initial_balance" render={({ field }) => (
              <FormItem><FormLabel>Saldo Awal</FormLabel><FormControl><Input type="number" {...field} /></FormControl><FormMessage /></FormItem>
            )} />
            <Button type="submit" className="w-full" disabled={mutation.isPending}>Simpan</Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}