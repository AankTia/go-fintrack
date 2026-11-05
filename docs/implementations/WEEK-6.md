# WEEK 6: Frontend Integration

## Day 16: Next.js Setup & Auth UI

Saya menyarankan struktur: `go-fintrack/frontend` agar terpisah dari kode Go.

### 🛠️ Langkah 1: Inisialisasi Project & Dependencies

Jalankan perintah ini di terminal (di root project):

#### 1. Create Next.js App
```bash
npx create-next-app@latest frontend --typescript --tailwind --eslint --app

# Masuk ke folder frontend
cd frontend 
```

#### 2. Install Core Dependencies
```bash
npm install axios @tanstack/react-query zod react-hook-form @hookform/resolvers lucide-react
```

#### 3. Init shadcn/ui
```bash
npx shadcn@latest init
# Pilih: Default style, Slate color, CSS variables: Yes
```

#### 4. Install UI Components yang dibutuhkan

```bash
npx shadcn@latest add button input card form toast dialog
```

———

### 🛠️ Langkah 2: Konfigurasi Environment & API Client

#### 1. Buat file .env.local

Buat di /frontend/.env.local:

NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

#### 2. Buat API Client Wrapper

Buat file di /frontend/lib/api.ts:

import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor: Otomatis tambah Bearer Token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response Interceptor: Handle 401 Unauthorized
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;

———

### 🛠️ Langkah 3: Implementasi Auth UI & Logic

#### 1. Schema Validasi (Zod)

Buat file di /frontend/lib/validations/auth.ts:

import * as z from "zod";

export const registerSchema = z.object({
  email: z.string().email("Email tidak valid"),
  password: z.string().min(6, "Password minimal 6 karakter"),
  full_name: z.string().min(2, "Nama lengkap harus diisi"),
});

export const loginSchema = z.object({
  email: z.string().email("Email tidak valid"),
  password: z.string().min(1, "Password wajib diisi"),
});

#### 2. Halaman Register (/frontend/app/register/page.tsx)

"use client";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { registerSchema } from "@/lib/validations/auth";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import api from "@/lib/api";
import { useRouter } from "next/navigation";

export default function RegisterPage() {
  const { toast } = useToast();
  const router = useRouter();
  const form = useForm({
    resolver: zodResolver(registerSchema),
    defaultValues: { email: "", password: "", full_name: "" },
  });

  async function onSubmit(values: any) {
    try {
      await api.post("/auth/register", values);
      toast({ title: "Berhasil!", description: "Akun berhasil dibuat. Silakan login." });
      router.push("/login");
    } catch (error: any) {
      toast({ variant: "destructive", title: "Error", description: error.response?.data?.message || "Terjadi kesalahan" });
    }
  }

  return (
    <div className="flex justify-center items-center min-h-screen">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="w-full max-w-md space-y-4 p-6 border rounded-lg shadow">
          <h1 className="text-2xl font-bold text-center">Register</h1>
          <FormField control={form.control} name="full_name" render={({ field }) => (
            <FormItem><FormLabel>Nama Lengkap</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
          )} />
          <FormField control={form.control} name="email" render={({ field }) => (
            <FormItem><FormLabel>Email</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
          )} />
          <FormField control={form.control} name="password" render={({ field }) => (
            <FormItem><FormLabel>Password</FormLabel><FormControl><Input type="password" {...field} /></FormControl><FormMessage /></FormItem>
          )} />
          <Button type="submit" className="w-full">Daftar</Button>
          <p className="text-center text-sm">Sudah punya akun? <a href="/login" className="text-blue-500">Login</a></p>
        </form>
      </Form>
    </div>
  );
}

#### 3. Halaman Login (/frontend/app/login/page.tsx)

"use client";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema } from "@/lib/validations/auth";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import api from "@/lib/api";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const { toast } = useToast();
  const router = useRouter();
  const form = useForm({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

  async function onSubmit(values: any) {
    try {
      const response = await api.post("/auth/login", values);
      localStorage.setItem('token', response.data.token);
      toast({ title: "Welcome back!" });
      router.push("/dashboard");
    } catch (error: any) {
      toast({ variant: "destructive", title: "Login Gagal", description: error.response?.data?.message || "Kredensial salah" });
    }
  }

  return (
    <div className="flex justify-center items-center min-h-screen">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="w-full max-w-md space-y-4 p-6 border rounded-lg shadow">
          <h1 className="text-2xl font-bold text-center">Login</h1>
          <FormField control={form.control} name="email" render={({ field }) => (
            <FormItem><FormLabel>Email</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
          )} />
          <FormField control={form.control} name="password" render={({ field }) => (
            <FormItem><FormLabel>Password</FormLabel><FormControl><Input type="password" {...field} /></FormControl><FormMessage /></FormItem>
          )} />
          <Button type="submit" className="w-full">Login</Button>
          <p className="text-center text-sm">Belum punya akun? <a href="/register" className="text-blue-500">Register</a></p>
        </form>
      </Form>
    </div>
  );
}

———

### 🛠️ Langkah 4: Middleware Client-Side (Protection)

Buat file di /frontend/middleware.ts (Root of frontend):
Catatan: Karena kita menggunakan localStorage untuk token (client-side), protection yang paling simpel untuk MVP adalah membuat Higher Order Component atau mengecek di layout.

Namun untuk Next.js App Router, buat file /frontend/app/dashboard/layout.tsx:

"use client";
import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
    }
  }, [router]);

  return <>{children}</>;
}

### 🏁 Verifikasi Akhir

1. Jalankan backend Go Anda.
2. Jalankan frontend: npm run dev.
3. Test alur: Register → Login → Akses /dashboard.
4. Coba hapus token di Application tab browser → Refresh halaman → Harus terlempar ke /login.

## Day 17: Accounts & Transactions Frontend

Untuk mengimplementasikan Day 17, kita akan membangun inti dari manajemen keuangan: Akun (tempat uang berada) dan Transaksi (aliran uang).

Kita akan menggunakan TanStack Query (useQuery & useMutation) untuk mengelola state server agar data otomatis terupdate tanpa perlu refresh halaman manual.

### 🛠️ Langkah 1: Setup Query Client

Sebelum masuk ke fitur, kita harus membungkus aplikasi dengan QueryClientProvider.

Edit file /frontend/app/layout.tsx:

"use client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "@/components/ui/toaster";
import "./globals.css";

const queryClient = new QueryClient();

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <QueryClientProvider client={queryClient}>
          {children}
          <Toaster />
        </QueryClientProvider>
      </body>
    </html>
  );
}

———

### 🛠️ Langkah 2: Account Management UI

Kita akan membuat tampilan saldo dalam bentuk Card dan modal untuk menambah akun.

#### 1. Account List Component

Buat file /frontend/components/accounts/account-list.tsx:

"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Wallet } from "lucide-react";

export function AccountList() {
  const { data: accounts, isLoading } = useQuery({
    queryKey: ["accounts"],
    queryFn: async () => {
      const res = await api.get("/accounts");
      return res.data;
    },
  });

  if (isLoading) return <p>Loading accounts...</p>;

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
      {accounts?.map((acc: any) => (
        <Card key={acc.id}>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">{acc.name}</CardTitle>
            <Wallet className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR" }).format(acc.balance)}
            </div>
            <p className="text-xs text-muted-foreground">{acc.type}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

#### 2. Create Account Modal

Buat file /frontend/components/accounts/create-account-modal.tsx:

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

———

### 🛠️ Langkah 3: Transaction List & Form

#### 1. Transaction Table Component

Buat file /frontend/components/transactions/transaction-table.tsx:

"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

export function TransactionTable() {
  const { data: transactions, isLoading } = useQuery({
    queryKey: ["transactions"],
    queryFn: async () => {
      const res = await api.get("/transactions");
      return res.data;
    },
  });

  if (isLoading) return <p>Loading transactions...</p>;

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Tanggal</TableHead>
            <TableHead>Kategori</TableHead>
            <TableHead>Keterangan</TableHead>
            <TableHead className="text-right">Nominal</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {transactions?.map((tx: any) => (
            <TableRow key={tx.id}>
              <TableCell>{new Date(tx.date).toLocaleDateString("id-ID")}</TableCell>
              <TableCell>{tx.category.name}</TableCell>
              <TableCell>{tx.description}</TableCell>
              <TableCell className={`text-right font-medium ${tx.type === 'expense' ? 'text-red-500' : 'text-green-500'}`}>
                {tx.type === 'expense' ? '-' : '+'} {new Intl.NumberFormat("id-ID").format(tx.amount)}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

#### 2. Transaction Form Modal

Buat file /frontend/components/transactions/create-transaction-modal.tsx:

"use client";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useToast } from "@/components/ui/use-toast";

const txSchema = z.object({
  account_id: z.string().min(1, "Pilih akun"),
  category_id: z.string().min(1, "Pilih kategori"),
  amount: z.coerce.number().min(1),
  type: z.enum(["income", "expense"]),
  description: z.string().optional(),
  date: z.string().min(1),
});

export function CreateTransactionModal() {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const form = useForm({ resolver: zodResolver(txSchema), defaultValues: { type: "expense", date: new Date().toISOString().split('T')[0] } });

  // Fetch data untuk dropdown
  const { data: accounts } = useQuery({ queryKey: ["accounts"], queryFn: async () => (await api.get("/accounts")).data });
  const { data: categories } = useQuery({ queryKey: ["categories"], queryFn: async () => (await api.get("/categories")).data });

  const mutation = useMutation({
    mutationFn: (values: any) => api.post("/transactions", values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
      queryClient.invalidateQueries({ queryKey: ["accounts"] }); // Update saldo akun
      toast({ title: "Sukses", description: "Transaksi dicatat" });
      form.reset();
    },
  });

  return (
    <Dialog>
      <DialogTrigger asChild><Button variant="destructive">+ Catat Transaksi</Button></DialogTrigger>
      <DialogContent>
        <DialogHeader><DialogTitle>Tambah Transaksi</DialogTitle></DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit((v) => mutation.mutate(v))} className="space-y-4">
            <FormField control={form.control} name="type" render={({ field }) => (
              <FormItem><FormLabel>Tipe</FormLabel><Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl><SelectTrigger><SelectValue /></SelectTrigger></FormControl>
                <SelectContent><SelectItem value="income">Pemasukan</SelectItem><SelectItem value="expense">Pengeluaran</SelectItem></SelectContent>
              </Select><FormMessage /></FormItem>
            )} />
            <FormField control={form.control} name="account_id" render={({ field }) => (
              <FormItem><FormLabel>Akun</FormLabel><Select onValueChange={field.onChange}>
                <FormControl><SelectTrigger><SelectValue placeholder="Pilih Akun" /></SelectTrigger></FormControl>
                <SelectContent>{accounts?.map((a: any) => <SelectItem key={a.id} value={a.id}>{a.name}</SelectItem>)}</SelectContent>
              </Select><FormMessage /></FormItem>
            )} />
            <FormField control={form.control} name="category_id" render={({ field }) => (
              <FormItem><FormLabel>Kategori</FormLabel><Select onValueChange={field.onChange}>
                <FormControl><SelectTrigger><SelectValue placeholder="Pilih Kategori" /></SelectTrigger></FormControl>
                <SelectContent>{categories?.map((c: any) => <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>)}</SelectContent>
              </Select><FormMessage /></FormItem>
            )} />
            <FormField control={form.control} name="amount" render={({ field }) => (
              <FormItem><FormLabel>Nominal</FormLabel><FormControl><Input type="number" {...field} /></FormControl><FormMessage /></FormItem>
            )} />
            <FormField control={form.control} name="description" render={({ field }) => (
              <FormItem><FormLabel>Keterangan</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
            )} />
            <Button type="submit" className="w-full" disabled={mutation.isPending}>Simpan Transaksi</Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

———

### 🛠️ Langkah 4: Integrasi ke Dashboard Page

Sekarang kita satukan semua komponen di satu halaman.

Edit file /frontend/app/dashboard/page.tsx:

import { AccountList } from "@/components/accounts/account-list";
import { CreateAccountModal } from "@/components/accounts/create-account-modal";
import { TransactionTable } from "@/components/transactions/transaction-table";
import { CreateTransactionModal } from "@/components/transactions/create-transaction-modal";

export default function DashboardPage() {
  return (
    <div className="p-8 max-w-6xl mx-auto space-y-8">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Financial Dashboard</h1>
        <div className="flex gap-2">
          <CreateAccountModal />
          <CreateTransactionModal />
        </div>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-4">My Accounts</h2>
        <AccountList />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-4">Recent Transactions</h2>
        <TransactionTable />
      </section>
    </div>
  );
}

### 🏁 Verifikasi Akhir (Day 17)

1. Buka /dashboard.
2. Tambah Akun Baru → Pastikan Card akun muncul otomatis.
3. Catat Transaksi Pengeluaran → Pilih akun yang baru dibuat.
4. Cek Saldo Akun → Saldo harus berkurang otomatis (karena invalidateQueries(["accounts"])).
5. Cek Tabel Transaksi → Transaksi baru harus muncul di baris teratas.

## Day 18: Dashboard & Budget Frontend

Untuk menyelesaikan Day 18, kita akan fokus pada visualisasi data menggunakan Recharts dan implementasi kontrol budget. Ini adalah bagian yang memberikan "insight" kepada pengguna tentang kesehatan keuangan mereka.

### 🛠️ Langkah 1: Install Dependencies Visualisasi

Jalankan perintah ini di folder /frontend:

npm install recharts
# Tambahkan komponen shadcn untuk Progress bar
npx shadcn-ui@latest add progress

———

### 🛠️ Langkah 2: Dashboard Summary & Charts

Kita akan membuat komponen visualisasi yang mengambil data dari /api/v1/dashboard.

#### 1. KPI Cards Component

Buat file /frontend/components/dashboard/kpi-cards.tsx:

"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowUpCircle, ArrowDownCircle, Wallet, TrendingUp } from "lucide-react";

export function KPICards() {
  const { data: summary, isLoading } = useQuery({
    queryKey: ["dashboard-summary"],
    queryFn: async () => (await api.get("/dashboard/summary")).data,
  });

  if (isLoading) return <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8"><div className="h-24 bg-gray-100 animate-pulse rounded-xl" /></div>;

  const formatIDR = (val: number) => new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR" }).format(val);

  const stats = [
    { title: "Total Balance", value: summary.total_balance, icon: Wallet, color: "text-blue-600" },
    { title: "Monthly Income", value: summary.monthly_income, icon: ArrowUpCircle, color: "text-green-600" },
    { title: "Monthly Expense", value: summary.monthly_expense, icon: ArrowDownCircle, color: "text-red-600" },
    { title: "Net Cash Flow", value: summary.net_cash_flow, icon: TrendingUp, color: "text-purple-600" },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
      {stats.map((s) => (
        <Card key={s.title}>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">{s.title}</CardTitle>
            <s.icon className={`h-4 w-4 ${s.color}`} />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatIDR(s.value)}</div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

#### 2. Cash Flow Chart Component

Buat file /frontend/components/dashboard/cashflow-chart.tsx:

"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from "recharts";

export function CashFlowChart() {
  const { data: history } = useQuery({
    queryKey: ["dashboard-history"],
    queryFn: async () => (await api.get("/dashboard/history")).data,
  });

  return (
    <Card className="col-span-1 md:col-span-2">
      <CardHeader><CardTitle>Cash Flow Tren (6 Bulan Terakhir)</CardTitle></CardHeader>
      <CardContent className="h-[300px]">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={history}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="month" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Bar dataKey="income" fill="#22c55e" name="Pemasukan" />
            <Bar dataKey="expense" fill="#ef4444" name="Pengeluaran" />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}

———

### 🛠️ Langkah 3: Budget Tracking UI

Kita akan membuat list progress bar untuk setiap kategori budget.

#### 1. Budget Progress List

Buat file /frontend/components/budget/budget-list.tsx:

"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";

export function BudgetList() {
  const { data: budgets, isLoading } = useQuery({
    queryKey: ["budgets"],
    queryFn: async () => (await api.get("/budgets")).data,
  });

  if (isLoading) return <p>Loading budgets...</p>;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {budgets?.map((b: any) => {
        const percentage = (b.current_spending / b.amount) * 100;
        const isOver = b.is_overlimit || percentage > 100;

        return (
          <Card key={b.id}>
            <CardContent className="pt-6">
              <div className="flex justify-between mb-2">
                <span className="font-medium">{b.category.name}</span>
                <span className={`text-sm ${isOver ? 'text-red-500 font-bold' : 'text-muted-foreground'}`}>
                  {new Intl.NumberFormat("id-ID").format(b.current_spending)} / {new Intl.NumberFormat("id-ID").format(b.amount)}
                </span>
              </div>
              <Progress value={percentage} className={isOver ? "bg-red-100" : ""} />
              {isOver && <p className="text-xs text-red-500 mt-1 text-right">Budget Terlampaui!</p>}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}

#### 2. Budget Management Modal

Buat file /frontend/components/budget/manage-budget-modal.tsx:

"use client";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useToast } from "@/components/ui/use-toast";

const budgetSchema = z.object({
  category_id: z.string().min(1, "Pilih kategori"),
  amount: z.coerce.number().min(1, "Budget minimal 1"),
});

export function ManageBudgetModal() {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const form = useForm({ resolver: zodResolver(budgetSchema), defaultValues: { amount: 0 } });

  const { data: categories } = useQuery({ queryKey: ["categories"], queryFn: async () => (await api.get("/categories")).data });

  const mutation = useMutation({
    mutationFn: (values: any) => api.post("/budgets", values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["budgets"] });
      toast({ title: "Budget Diupdate", description: "Limit budget berhasil disimpan" });
      form.reset();
    },
  });

  return (
    <Dialog>
      <DialogTrigger asChild><Button variant="outline">+ Set Budget</Button></DialogTrigger>
      <DialogContent>
        <DialogHeader><DialogTitle>Atur Budget Bulanan</DialogTitle></DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit((v) => mutation.mutate(v))} className="space-y-4">
            <FormField control={form.control} name="category_id" render={({ field }) => (
              <FormItem><FormLabel>Kategori</FormLabel><Select onValueChange={field.onChange}>
                <FormControl><SelectTrigger><SelectValue placeholder="Pilih Kategori" /></SelectTrigger></FormControl>
                <SelectContent>{categories?.map((c: any) => <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>)}</SelectContent>
              </Select><FormMessage /></FormItem>
            )} />
            <FormField control={form.control} name="amount" render={({ field }) => (
              <FormItem><FormLabel>Limit Nominal</FormLabel><FormControl><Input type="number" {...field} /></FormControl><FormMessage /></FormItem>
            )} />
            <Button type="submit" className="w-full" disabled={mutation.isPending}>Simpan Budget</Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

———

### 🛠️ Langkah 4: Integrasi ke Dashboard Page

Sekarang kita update /frontend/app/dashboard/page.tsx untuk menampilkan semua visualisasi ini.

import { KPICards } from "@/components/dashboard/kpi-cards";
import { CashFlowChart } from "@/components/dashboard/cashflow-chart";
import { BudgetList } from "@/components/budget/budget-list";
import { ManageBudgetModal } from "@/components/budget/manage-budget-modal";
import { AccountList } from "@/components/accounts/account-list";
import { CreateAccountModal } from "@/components/accounts/create-account-modal";
import { TransactionTable } from "@/components/transactions/transaction-table";
import { CreateTransactionModal } from "@/components/transactions/create-transaction-modal";

export default function DashboardPage() {
  return (
    <div className="p-8 max-w-7xl mx-auto space-y-10">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight">Financial Overview</h1>
        <div className="flex gap-2">
          <CreateAccountModal />
          <CreateTransactionModal />
          <ManageBudgetModal />
        </div>
      </div>

      {/* Section 1: High Level Summary */}
      <section className="space-y-4">
        <KPICards />
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <CashFlowChart />
          <div className="col-span-1">
            <h2 className="text-xl font-semibold mb-4">Monthly Budget</h2>
            <BudgetList />
          </div>
        </div>
      </section>

      {/* Section 2: Detail Data */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-1">
          <h2 className="text-xl font-semibold mb-4">My Accounts</h2>
          <AccountList />
        </div>
        <div className="lg:col-span-2">
          <h2 className="text-xl font-semibold mb-4">Recent Transactions</h2>
          <TransactionTable />
        </div>
      </div>
    </div>
  );
}

### 🏁 Verifikasi Akhir (Day 18)

1. Cek KPI Cards: Apakah Total Balance, Income, dan Expense muncul benar?
2. Cek Chart: Apakah grafik batang menampilkan data pemasukan & pengeluaran per bulan?
3. Set Budget: Tambahkan budget untuk kategori "Makanan" → cek apakah muncul di BudgetList.
4. Test Overlimit: Buat transaksi pengeluaran kategori "Makanan" yang melebihi budget → pastikan Progress Bar berubah menjadi merah dan muncul peringatan "Budget Terlampaui!".
5. Responsivitas: Kecilkan browser ke ukuran mobile, pastikan grid berubah menjadi 1 kolom.