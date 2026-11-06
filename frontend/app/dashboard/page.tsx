import { AccountList } from "@/components/accounts/account-list";
import { CreateAccountModal } from "@/components/accounts/create-account-modal";
import { TransactionTable } from "@/components/transactions/transaction-table";
import { CreateTransactionModal } from "@/components/transactions/create-transaction-modal";
import { ManageBudgetModal } from "@/components/budget/manage-budget-modal";
import { KPICards } from "@/components/dashboard/kpi-cards";
import { CashFlowChart } from "@/components/dashboard/cashflow-chart";
import { BudgetList } from "@/components/budget/budget-list";

export default function DashboardPage() {
  return (
    <div className="p-8 max-w-7xl mx-auto space-y-10">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight">
          Financial Overview
        </h1>
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
