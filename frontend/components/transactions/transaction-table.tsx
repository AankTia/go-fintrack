"use client";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

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
              <TableCell>
                {new Date(tx.date).toLocaleDateString("id-ID")}
              </TableCell>
              <TableCell>{tx.category.name}</TableCell>
              <TableCell>{tx.description}</TableCell>
              <TableCell
                className={`text-right font-medium ${tx.type === "expense" ? "text-red-500" : "text-green-500"}`}
              >
                {tx.type === "expense" ? "-" : "+"}{" "}
                {new Intl.NumberFormat("id-ID").format(tx.amount)}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
