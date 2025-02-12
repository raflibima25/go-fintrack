package service

import (
	"errors"
	"fmt"
	"go-fintrack/internal/payload/response"
	"go-fintrack/internal/utility"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DashboardService struct {
	DB            *gorm.DB
	dashboardUtil *utility.DashboardUtil
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{
		DB:            db,
		dashboardUtil: &utility.DashboardUtil{DB: db},
	}
}

func (s *DashboardService) GetFinancialOverview(userID uint) (*response.RespFinancialOverview, error) {
	logrus.Info("Getting financial overview for user: ", userID)

	var overview response.RespFinancialOverview
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 4)

	// get current balance
	wg.Add(1)
	go func() {
		defer wg.Done()
		balance, err := s.dashboardUtil.CalculateCurrentBalance(userID)
		if err != nil {
			logrus.Errorf("Failed to calculate current balance: %v", err)
			errChan <- err
			return
		}
		mu.Lock()
		overview.CurrentBalance = balance
		mu.Unlock()
	}()

	// get monthly income
	wg.Add(1)
	go func() {
		defer wg.Done()
		startOfMonth := time.Now().UTC().Format("2006-01-01")
		income, err := s.dashboardUtil.CalculateMonthlyIncome(userID, startOfMonth)
		if err != nil {
			logrus.Errorf("Failed to calculate monthly income: %v", err)
			errChan <- err
			return
		}
		mu.Lock()
		overview.MonthlyIncome = income
		mu.Unlock()
	}()

	// get monthly expense
	wg.Add(1)
	go func() {
		defer wg.Done()
		startOfMonth := time.Now().UTC().Format("2006-01-01")
		expense, err := s.dashboardUtil.CalculateMonthlyExpense(userID, startOfMonth)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		overview.MonthlyExpense = expense
		mu.Unlock()
	}()

	// get total savings
	wg.Add(1)
	go func() {
		defer wg.Done()
		savings, err := s.dashboardUtil.CalculateTotalSavings(userID)
		if err != nil {
			logrus.Errorf("Failed to calculate total savings: %v", err)
			errChan <- err
			return
		}
		mu.Lock()
		overview.TotalSavings = savings
		mu.Unlock()
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, errors.New("failed to get financial overview")
		}
	}

	logrus.Info("Successfully retrieved financial overview")
	return &overview, nil
}

func (s *DashboardService) GetDashboardCharts(userID uint) (*response.RespDashboardCharts, error) {
	logrus.Info("Getting dashboard charts for user: ", userID)

	var charts response.RespDashboardCharts
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 3)

	defer close(errChan)

	// get income vs expense
	wg.Add(1)
	go func() {
		defer wg.Done()
		labels, incomeData, expenseData, err := s.dashboardUtil.GetLastSixMonthsData(userID)
		if err != nil {
			logrus.Errorf("Failed to get income vs expense data: %v", err)
			errChan <- err
			return
		}

		mu.Lock()
		charts.IncomeVsExpense = response.RespIncomeVsExpense{
			Labels: labels,
			Datasets: []response.ChartDataset{
				{
					Label:           "Income",
					Data:            incomeData,
					BorderColor:     "#10B981",
					BackgroundColor: "rgba(16, 185, 129, 0.1)",
				},
				{
					Label:           "Expense",
					Data:            expenseData,
					BorderColor:     "#EF4444",
					BackgroundColor: "rgba(239, 68, 68, 0.1)",
				},
			},
		}
		mu.Unlock()
	}()

	// get category distribution
	wg.Add(1)
	go func() {
		defer wg.Done()
		labels, data, err := s.dashboardUtil.GetCategoryDistribution(userID)
		if err != nil {
			logrus.Errorf("Failed to get category distribution data: %v", err)
			errChan <- err
			return
		}

		mu.Lock()
		charts.CategoryDistribution = response.CategoryDistribution{
			Labels: labels,
			Datasets: []struct {
				Data            []float64 `json:"data"`
				BackgroundColor []string  `json:"background_color"`
			}{
				{
					Data: data,
					BackgroundColor: []string{
						"#10B981", "#3B82F6", "#F59E0B",
						"#6366F1", "#EC4899", "#8B5CF6",
					},
				},
			},
		}
		mu.Unlock()
	}()

	// get top expenses
	wg.Add(1)
	go func() {
		defer wg.Done()
		labels, data, err := s.dashboardUtil.GetTopExpenseCategories(userID, 5)
		if err != nil {
			logrus.Errorf("Failed to get top expenses data: %v", err)
			errChan <- err
			return
		}

		mu.Lock()
		charts.TopExpenses = response.TopExpenses{
			Labels: labels,
			Datasets: []struct {
				Data            []float64 `json:"data"`
				BackgroundColor string    `json:"background_color"`
			}{
				{
					Data:            data,
					BackgroundColor: "#3B82F6",
				},
			},
		}
		mu.Unlock()
	}()

	wg.Wait()

	select {
	case err := <-errChan:
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboard charts: %v", err)
		}
	default:
	}

	logrus.Info("Successfully retrieved dashboard charts")
	return &charts, nil
}
