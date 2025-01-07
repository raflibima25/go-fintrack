package response

import "time"

type TransactionResponse struct {
	ID          uint      `json:"id"`
	CategoryID  uint      `json:"category_id"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TransactionSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Summary      TransactionSummary    `json:"summary"`
	Pagination   Pagination            `json:"pagination"`
}
