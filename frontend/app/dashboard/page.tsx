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