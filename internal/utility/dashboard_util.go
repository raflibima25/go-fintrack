package utility

import (
	"gorm.io/gorm"
)

type DashboardUtil struct {
	DB *gorm.DB
}

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
