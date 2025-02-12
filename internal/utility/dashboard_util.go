package utility

import (
	"time"

	"gorm.io/gorm"
)

type DashboardUtil struct {
	DB *gorm.DB
}

// Financial Overview
func (u *DashboardUtil) CalculateCurrentBalance(userID uint) (float64, error) {
	var balance float64
	err := u.DB.Table("transactions").
		Select("COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 0)").
		Where("user_id = ?", userID).
		Row().
		Scan(&balance)
	return balance, err
}

func (u *DashboardUtil) CalculateMonthlyIncome(userID uint, startOfMonth string) (float64, error) {
	var income float64
	err := u.DB.Table("transactions").
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND type = 'income' AND date >= ?", userID, startOfMonth).
		Row().
		Scan(&income)
	return income, err
}

func (u *DashboardUtil) CalculateMonthlyExpense(userID uint, startOfMonth string) (float64, error) {
	var expense float64
	err := u.DB.Table("transactions").
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND type = 'expense' AND date >= ?", userID, startOfMonth).
		Row().
		Scan(&expense)
	return expense, err
}

func (u *DashboardUtil) CalculateTotalSavings(userID uint) (float64, error) {
	var savings float64
	err := u.DB.Table("transactions").
		Select("COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 0)").
		Where("user_id = ?", userID).
		Row().
		Scan(&savings)
	return savings, err
}

// Expense Analysis
func (u *DashboardUtil) GetLastSixMonthsData(userID uint) ([]string, []float64, []float64, error) {
	var labels []string
	var incomeData []float64
	var expenseData []float64

	// Get current time
	now := time.Now().UTC()

	// Generate last 6 months
	for i := 5; i >= 0; i-- {
		date := now.AddDate(0, -i, 0)
		startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

		// Get monthly income
		var income float64
		err := u.DB.Table("transactions").
			Select("COALESCE(SUM(amount), 0)").
			Where("user_id = ? AND type = 'income' AND deleted_at IS NULL AND date BETWEEN ? AND ?",
				userID, startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02")).
			Row().
			Scan(&income)
		if err != nil {
			return nil, nil, nil, err
		}

		// Get monthly expense
		var expense float64
		err = u.DB.Table("transactions").
			Select("COALESCE(SUM(amount), 0)").
			Where("user_id = ? AND type = 'expense' AND deleted_at IS NULL AND date BETWEEN ? AND ?",
				userID, startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02")).
			Row().
			Scan(&expense)
		if err != nil {
			return nil, nil, nil, err
		}

		labels = append(labels, date.Format("Jan"))
		incomeData = append(incomeData, income)
		expenseData = append(expenseData, expense)
	}

	return labels, incomeData, expenseData, nil
}

func (u *DashboardUtil) GetCategoryDistribution(userID uint) ([]string, []float64, error) {
	type CategoryTotal struct {
		Category string  `gorm:"column:category_name"`
		Total    float64 `gorm:"column:total"`
	}

	var results []CategoryTotal

	err := u.DB.Table("transactions").
		Select("categories.name as category_name, COALESCE(SUM(transactions.amount), 0) as total").
		Joins("LEFT JOIN categories ON transactions.category_id = categories.id").
		Where("transactions.user_id = ? AND transactions.type = 'expense' AND transactions.deleted_at IS NULL AND categories.deleted_at IS NULL", userID).
		Group("categories.name").
		Order("total DESC").
		Find(&results).Error
	if err != nil {
		return nil, nil, err
	}

	var labels []string
	var data []float64

	for _, result := range results {
		labels = append(labels, result.Category)
		data = append(data, result.Total)
	}

	return labels, data, nil
}

func (u *DashboardUtil) GetTopExpenseCategories(userID uint, limit int) ([]string, []float64, error) {
	type CategoryTotal struct {
		Category string  `gorm:"column:category_name"`
		Total    float64 `gorm:"column:total"`
	}

	var results []CategoryTotal

	err := u.DB.Table("transactions").
		Select("categories.name as category_name, COALESCE(SUM(transactions.amount), 0) as total").
		Joins("LEFT JOIN categories ON transactions.category_id = categories.id").
		Where("transactions.user_id = ? AND transactions.type = 'expense' AND transactions.deleted_at IS NULL AND categories.deleted_at IS NULL", userID).
		Group("categories.name").
		Order("total DESC").
		Limit(limit).
		Find(&results).Error
	if err != nil {
		return nil, nil, err
	}

	var labels []string
	var data []float64

	for _, result := range results {
		labels = append(labels, result.Category)
		data = append(data, result.Total)
	}

	return labels, data, nil
}
