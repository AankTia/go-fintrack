package model

import "time"

type DashboardSummary struct {
	TotalBalance int64      `json:"total_balance"` // Total saldo semua akun
	TotalIncome  int64      `json:"total_income"`  // Pemasukan bulan ini
	TotalExpense int64      `json:"total_expense"` // Pengeluaran bulan ini
	NetCashFlow  int64      `json:"net_cash_flow"` // Income - Expense
	ReportMonth  time.Month `json:"report_month"`
	ReportYear   int        `json:"report_year"`
}

type CategoryBreakdown struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  int64   `json:"total_amount"`
	Percentage   float64 `json:"percentage"` // Persentase terhadap total pengeluaran bulanan
}
