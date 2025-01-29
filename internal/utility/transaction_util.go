package utility

import (
	"fmt"
	"go-fintrack/internal/payload/entity"
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/payload/response"
	"time"

	"gorm.io/gorm"
)

type TransactionUtil struct {
	DB *gorm.DB
}

func (u *TransactionUtil) BuildFilterQuery(query *gorm.DB, filter request.TransactionFilter) *gorm.DB {
	newQuery := query

	// filter tanggal
	if filter.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", filter.StartDate); err == nil {
			newQuery = newQuery.Where("date >= ?", startDate)
		}
	}

	if filter.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", filter.EndDate); err == nil {
			newQuery = newQuery.Where("date <= ?", endDate)
		}
	}

	// filter kategori
	if filter.CategoryID != 0 {
		newQuery = newQuery.Where("category_id = ?", filter.CategoryID)
	}

	// filter tipe transaksi
	if filter.Type != "" {
		newQuery = newQuery.Where("type = ?", filter.Type)
	}

	return newQuery
}

func (u *TransactionUtil) CalculateTransactionSummary(baseQuery *gorm.DB, filter request.TransactionFilter) (*response.TransactionSummary, error) {
	var totalIncome, totalExpense float64

	// hitung total income
	incomeQuery := baseQuery.Session(&gorm.Session{})
	if err := incomeQuery.Model(&entity.Transaction{}).
		Where("type = ?", "income").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total income: %v", err)
	}

	// hitung total expense
	expenseQuery := baseQuery.Session(&gorm.Session{})
	if err := expenseQuery.Model(&entity.Transaction{}).
		Where("type = ?", "expense").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total expense: %v", err)
	}

	return &response.TransactionSummary{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
	}, nil
}
