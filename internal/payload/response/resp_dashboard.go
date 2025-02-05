package response

type RespFinancialOverview struct {
	CurrentBalance float64 `json:"current_balance"`
	MonthlyIncome  float64 `json:"monthly_income"`
	MonthlyExpense float64 `json:"monthly_expense"`
	TotalSavings   float64 `json:"total_savings"`
}
