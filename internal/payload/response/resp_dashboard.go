package response

// Financial Overview
type RespFinancialOverview struct {
	CurrentBalance float64 `json:"current_balance"`
	MonthlyIncome  float64 `json:"monthly_income"`
	MonthlyExpense float64 `json:"monthly_expense"`
	TotalSavings   float64 `json:"total_savings"`
}

// Expense Analysis
type ChartDataset struct {
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	BorderColor     string    `json:"border_color,omitempty"`
	BackgroundColor string    `json:"background_color,omitempty"`
}

type RespIncomeVsExpense struct {
	Labels   []string       `json:"labels"`
	Datasets []ChartDataset `json:"datasets"`
}

type CategoryDistribution struct {
	Labels   []string `json:"labels"`
	Datasets []struct {
		Data            []float64 `json:"data"`
		BackgroundColor []string  `json:"background_color"`
	} `json:"datasets"`
}

type TopExpenses struct {
	Labels   []string `json:"labels"`
	Datasets []struct {
		Data            []float64 `json:"data"`
		BackgroundColor string    `json:"background_color"`
	} `json:"datasets"`
}

type RespDashboardCharts struct {
	IncomeVsExpense      RespIncomeVsExpense  `json:"income_vs_expense"`
	CategoryDistribution CategoryDistribution `json:"category_distribution"`
	TopExpenses          TopExpenses          `json:"top_expenses"`
}
