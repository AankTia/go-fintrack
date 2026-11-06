"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";

export function BudgetList() {
  // Ambil bulan saat ini dalam format YYYY-MM
  const currentPeriod = new Date().toISOString().slice(0, 7);

  const { data: budgets, isLoading } = useQuery({
     queryKey: ["budgets", currentPeriod],
    queryFn: async () => (await api.get(`/budgets?period=${currentPeriod}`)).data,
  });

  if (isLoading) return <p>Loading budgets...</p>;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {budgets?.map((b: any) => {
        const percentage = b.amount > 0 ? (b.current_spending / b.amount) * 100 : 0;
        const isOver = b.is_overlimit || percentage > 100;

        return (
          <Card key={b.id}>
            <CardContent className="pt-6">
              <div className="flex justify-between mb-2">
                <span className="font-medium">{b.category.name}</span>
                <span
                  className={`text-sm ${isOver ? "text-red-500 font-bold" : "text-muted-foreground"}`}
                >
                  {new Intl.NumberFormat("id-ID").format(b.current_spending)} /{" "}
                  {new Intl.NumberFormat("id-ID").format(b.amount)}
                </span>
              </div>
              <Progress
                value={percentage}
                className={isOver ? "bg-red-100" : ""}
              />
              {isOver && (
                <p className="text-xs text-red-500 mt-1 text-right">
                  Budget Terlampaui!
                </p>
              )}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
