package request

type CreateTransactionRequest struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Type        string  `json:"type" binding:"required,oneof=income expense"`
	Description string  `json:"description"`
	Date        string  `json:"date" binding:"required"`
}

type UpdateTransactionRequest struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Type        string  `json:"type" binding:"required,oneof=income expense"`
	Description string  `json:"description"`
	Date        string  `json:"date" binding:"required"`
}

type TransactionFilter struct {
	StartDate  string `form:"start_date"` // format 2006-01-02
	EndDate    string `form:"end_date"`   // format 2006-01-02
	CategoryID uint   `form:"category_id"`
	Type       string `form:"type" binding:"omitempty,oneof=income expense"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
}
