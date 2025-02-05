package service

import (
	"errors"
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
