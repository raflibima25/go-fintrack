package service

import (
	"errors"
	"go-manajemen-keuangan/internal/payload/entity"
	"go-manajemen-keuangan/internal/payload/request"
	"go-manajemen-keuangan/internal/payload/response"
	"gorm.io/gorm"
	"time"
)

type TransactionService struct {
	DB *gorm.DB
}

func (s *TransactionService) CreateTransaction(userID uint, req request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	// validasi category
	var category entity.Category
	if err := s.DB.First(&category, req.CategoryID).Error; err != nil {
		return nil, errors.New("category not found")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
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
