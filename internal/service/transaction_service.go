package service

import (
	"errors"
	"github.com/sirupsen/logrus"
	"go-manajemen-keuangan/internal/payload/entity"
	"go-manajemen-keuangan/internal/payload/request"
	"go-manajemen-keuangan/internal/payload/response"
	"gorm.io/gorm"
	"time"
)

type TransactionService struct {
	DB *gorm.DB
}

func (s *TransactionService) GetTransactionByUser(userID uint) (*response.TransactionListResponse, error) {
	var transactions []entity.Transaction

	if err := s.DB.Preload("Category").
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&transactions).Error; err != nil {
		return nil, errors.New("failed to get transactions")
	}

	// transform ke response format
	transactionResponses := make([]response.TransactionResponse, len(transactions))
	var totalIncome, totalExpense float64

	for i, tx := range transactions {
		if tx.Type == "income" {
			totalIncome += tx.Amount
		} else {
			totalExpense += tx.Amount
		}

		transactionResponses[i] = response.TransactionResponse{
			ID:          tx.ID,
			CategoryID:  tx.CategoryID,
			Category:    tx.Category.Name,
			Amount:      tx.Amount,
			Type:        tx.Type,
			Description: tx.Description,
			Date:        tx.Date,
			CreatedAt:   tx.CreatedAt,
			UpdatedAt:   tx.UpdatedAt,
		}
	}

	return &response.TransactionListResponse{
		Transaction: transactionResponses,
		Summary: response.TransactionSummary{
			TotalIncome:  totalIncome,
			TotalExpense: totalExpense,
			Balance:      totalIncome - totalExpense,
		},
	}, nil
}

func (s *TransactionService) CreateTransaction(userID uint, req request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	// validasi category
	var category entity.Category
	if err := s.DB.First(&category, req.CategoryID).Error; err != nil {
		logrus.Errorf("category not found: %v", err)
		return nil, errors.New("category not found")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		logrus.Errorf("invalid date format: %v", err)
		return nil, errors.New("invalid date format")
	}

	transaction := entity.Transaction{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Type:        req.Type,
		Description: req.Description,
		Date:        date,
	}

	if err := s.DB.Create(&transaction).Error; err != nil {
		logrus.Errorf("Error creating transaction: %v", err)
		return nil, errors.New("failed to create transaction")
	}

	return &response.TransactionResponse{
		ID:          transaction.ID,
		CategoryID:  transaction.CategoryID,
		Category:    category.Name,
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		Description: transaction.Description,
		Date:        transaction.Date,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}, nil
}

func (s *TransactionService) UpdateTransaction(userID uint, transactionID uint, req request.UpdateTransactionRequest) (*response.TransactionResponse, error) {
	var transaction entity.Transaction
	if err := s.DB.Where("id = ? AND user_id = ?", transactionID, userID).
		First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		logrus.Errorf("Error getting transaction: %v", err)
		return nil, errors.New("failed to get transaction")
	}

	var category entity.Category
	if err := s.DB.First(&category, req.CategoryID).Error; err != nil {
		logrus.Errorf("Error category not found: %v", err)
		return nil, errors.New("category not found")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		logrus.Errorf("Error invalid date format: %v", err)
		return nil, errors.New("invalid date format")
	}

	transaction.CategoryID = req.CategoryID
	transaction.Amount = req.Amount
	transaction.Type = req.Type
	transaction.Description = req.Description
	transaction.Date = date

	if err := s.DB.Save(&transaction).Error; err != nil {
		logrus.Errorf("Error update transaction: %v", err)
		return nil, errors.New("failed to update transaction")
	}

	return &response.TransactionResponse{
		ID:          transaction.ID,
		CategoryID:  transaction.CategoryID,
		Category:    category.Name,
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		Description: transaction.Description,
		Date:        transaction.Date,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}, nil
}

func (s *TransactionService) DeleteTransaction(userID uint, transactionID uint) error {
	result := s.DB.Where("id = ? AND user_id = ?", transactionID, userID).Delete(&entity.Transaction{})
	if result.Error != nil {
		logrus.Errorf("Error to delete transaction: %v", result.Error)
		return errors.New("failed to delete transaction")
	}

	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}
