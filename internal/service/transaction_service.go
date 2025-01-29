package service

import (
	"bytes"
	"errors"
	"fmt"
	"go-fintrack/internal/payload/entity"
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/payload/response"
	"go-fintrack/internal/utility"
	"math"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type TransactionService struct {
	DB              *gorm.DB
	transactionUtil *utility.TransactionUtil
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		DB:              db,
		transactionUtil: &utility.TransactionUtil{DB: db},
	}
}

func (s *TransactionService) GetTransactionByUser(userID uint, filter request.TransactionFilter) (*response.TransactionListResponse, error) {
	logrus.Infof("Applying filter: %+v", filter) // debug

	var transactions []entity.Transaction

	baseQuery := s.DB.Where("user_id = ?", userID)

	filteredQuery := s.transactionUtil.BuildFilterQuery(baseQuery, filter)

	// hitung total untuk pagination
	var total int64
	if err := filteredQuery.Model(&entity.Transaction{}).Count(&total).Error; err != nil {
		logrus.Errorf("Failed to count transaction: %v", err)
		return nil, errors.New("failed to count transaction")
	}

	summary, err := s.transactionUtil.CalculateTransactionSummary(filteredQuery, filter)
	if err != nil {
		logrus.Errorf("Failed to calculate transaction summary: %v", err)
		return nil, err
	}

	// terapkan pagination
	offset := (filter.Page - 1) * filter.Limit
	if err := filteredQuery.Preload("Category").
		Order("date DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&transactions).Error; err != nil {
		logrus.Errorf("Failed to get transactions: %v", err)
		return nil, errors.New("failed to get transactions")
	}

	// transform ke response format
	transactionResponses := make([]response.TransactionResponse, len(transactions))
	for i, tx := range transactions {
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
		Transactions: transactionResponses,
		Summary:      *summary,
		Pagination: response.Pagination{
			CurrentPage: filter.Page,
			TotalPage:   int(math.Ceil(float64(total) / float64(filter.Limit))),
			TotalItems:  total,
			ItemPerPage: filter.Limit,
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

func (s *TransactionService) ExportTransactionsExcel(userID uint, filter request.TransactionFilter) (*bytes.Buffer, error) {
	transactions, err := s.GetTransactionByUser(userID, filter)
	if err != nil {
		logrus.Errorf("Error getting transactions: %v", err)
		return nil, err
	}

	f := excelize.NewFile()

	// Buat sheet baru
	sheet := "Transactions"
	index, err := f.NewSheet(sheet)
	if err != nil {
		logrus.Errorf("Error creating sheet: %v", err)
		return nil, err
	}
	f.SetActiveSheet(index)

	// Set header
	headers := []string{"Date", "Type", "Category", "Amount", "Description"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		if err := f.SetCellValue(sheet, cell, header); err != nil {
			logrus.Errorf("Error setting header: %v", err)
			return nil, err
		}
	}

	// Isi data
	for i, tx := range transactions.Transactions {
		row := i + 2
		if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", row), tx.Date.Format("2006-01-02")); err != nil {
			return nil, err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("B%d", row), tx.Type); err != nil {
			return nil, err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("C%d", row), tx.Category); err != nil {
			return nil, err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("D%d", row), tx.Amount); err != nil {
			return nil, err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("E%d", row), tx.Description); err != nil {
			return nil, err
		}
	}

	// Tambah summary
	summaryRow := len(transactions.Transactions) + 4
	f.SetCellValue(sheet, fmt.Sprintf("A%d", summaryRow), "Summary")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", summaryRow), "Total Pemasukan")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", summaryRow), transactions.Summary.TotalIncome)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", summaryRow+1), "Total Pengeluaran")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", summaryRow+1), transactions.Summary.TotalExpense)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", summaryRow+2), "Saldo")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", summaryRow+2), transactions.Summary.Balance)

	// Styling
	if style, err := f.NewStyle(&excelize.Style{
		NumFmt: 44, // Format currency
	}); err == nil {
		// Set style untuk kolom amount dan summary
		for i := 2; i <= len(transactions.Transactions)+1; i++ {
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", i), fmt.Sprintf("D%d", i), style)
		}
		f.SetCellStyle(sheet, fmt.Sprintf("C%d", summaryRow), fmt.Sprintf("C%d", summaryRow+2), style)
	}

	// Save ke buffer
	buffer := new(bytes.Buffer)
	_, err = f.WriteTo(buffer)
	if err != nil {
		logrus.Errorf("Error writing to buffer: %v", err)
		return nil, err
	}

	return buffer, nil
}
