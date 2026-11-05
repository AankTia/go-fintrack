
"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Wallet } from "lucide-react";

export function AccountList() {
  const {data: accounts, isLoading} = useQuery({
    queryKey: ["accounts"],
    queryFn: async () => {
      const res = await api.get("/accounts");
      return res.data;
    }
  });

  if (isLoading) return <p>Loading accounts...</p>

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