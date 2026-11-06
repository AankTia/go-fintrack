"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";

export function CashFlowChart() {
  const { data: history } = useQuery({
    queryKey: ["dashboard-history"],
    queryFn: async () => (await api.get("/dashboard/history")).data,
  });

  return (
    <Card className="col-span-1 md:col-span-2">
      <CardHeader>
        <CardTitle>Cash Flow Tren (6 Bulan Terakhir)</CardTitle>
      </CardHeader>
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
