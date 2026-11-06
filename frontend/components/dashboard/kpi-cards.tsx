"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ArrowUpCircle,
  ArrowDownCircle,
  Wallet,
  TrendingUp,
} from "lucide-react";

export function KPICards() {
  const { data: summary, isLoading } = useQuery({
    queryKey: ["dashboard-summary"],
    queryFn: async () => (await api.get("/dashboard/summary")).data,
  });

  if (isLoading)
    return (
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div className="h-24 bg-gray-100 animate-pulse rounded-xl" />
      </div>
    );

  const formatIDR = (val: number) =>
    new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
    }).format(val);

  const stats = [
    {
      title: "Total Balance",
      value: summary.total_balance,
      icon: Wallet,
      color: "text-blue-600",
    },
    {
      title: "Monthly Income",
      value: summary.monthly_income,
      icon: ArrowUpCircle,
      color: "text-green-600",
    },
    {
      title: "Monthly Expense",
      value: summary.monthly_expense,
      icon: ArrowDownCircle,
      color: "text-red-600",
    },
    {
      title: "Net Cash Flow",
      value: summary.net_cash_flow,
      icon: TrendingUp,
      color: "text-purple-600",
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
      {stats.map((s) => (
        <Card key={s.title}>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {s.title}
            </CardTitle>
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
