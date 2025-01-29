package unit

import (
	"database/sql"
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/service"
	"io"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	service *service.TransactionService
	sqlDB   *sql.DB
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	var err error
	suite.sqlDB, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)

	dialector := mysql.New(mysql.Config{
		Conn:                      suite.sqlDB,
		SkipInitializeWithVersion: true,
	})

	newLogger := logger.New(
		log.New(io.Discard, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	suite.DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})
	assert.NoError(suite.T(), err)

	suite.service = service.NewTransactionService(suite.DB)
}

func (suite *TransactionServiceTestSuite) TearDownTest() {
	suite.sqlDB.Close()
}

func (suite *TransactionServiceTestSuite) TestGetTransactionByUser() {
	userID := uint(1)
	now := time.Now()
	filter := request.TransactionFilter{
		Page:  1,
		Limit: 10,
	}

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `transactions` WHERE user_id = ? AND `transactions`.`deleted_at` IS NULL")).
		WithArgs(userID).
		WillReturnRows(countRows)

	// Mock income query
	incomeRows := sqlmock.NewRows([]string{"sum"}).AddRow(1000.0)
	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(amount), 0) FROM `transactions` WHERE user_id = ? AND `transactions`.`deleted_at` IS NULL AND type = ?")).
		WithArgs(userID, "income").
		WillReturnRows(incomeRows)

	// Mock expense query
	expenseRows := sqlmock.NewRows([]string{"sum"}).AddRow(500.0)
	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(amount), 0) FROM `transactions` WHERE user_id = ? AND `transactions`.`deleted_at` IS NULL AND type = ?")).
		WithArgs(userID, "expense").
		WillReturnRows(expenseRows)

	// Mock transaction list query
	transactionRows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"user_id", "category_id", "amount", "type",
		"description", "date",
	}).
		AddRow(1, now, now, nil, userID, 1, 1000.0, "income", "Salary", now).
		AddRow(2, now, now, nil, userID, 2, 500.0, "expense", "Shopping", now)

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? AND `transactions`.`deleted_at` IS NULL ORDER BY date DESC LIMIT ?")).
		WithArgs(userID, 10).
		WillReturnRows(transactionRows)

	// Mock category preload (using IN clause)
	categoryRows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"user_id", "name",
	}).
		AddRow(1, now, now, nil, userID, "Category 1").
		AddRow(2, now, now, nil, userID, "Category 2")

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` IN (?,?) AND `categories`.`deleted_at` IS NULL")).
		WithArgs(1, 2).
		WillReturnRows(categoryRows)

	result, err := suite.service.GetTransactionByUser(userID, filter)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Transactions, 2)
	assert.Equal(suite.T(), float64(1000), result.Summary.TotalIncome)
	assert.Equal(suite.T(), float64(500), result.Summary.TotalExpense)
	assert.Equal(suite.T(), float64(500), result.Summary.Balance)
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction() {
	userID := uint(1)
	now := time.Now()
	date := "2025-01-29"
	req := request.CreateTransactionRequest{
		CategoryID:  1,
		Amount:      1000.0,
		Type:        "income",
		Description: "Salary",
		Date:        date,
	}

	// Mock category check
	categoryRows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"user_id", "name",
	}).AddRow(1, now, now, nil, userID, "Salary")

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? AND `categories`.`deleted_at` IS NULL ORDER BY `categories`.`id` LIMIT ?")).
		WithArgs(req.CategoryID, 1).
		WillReturnRows(categoryRows)

	// Mock create transaction
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `transactions` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`category_id`,`amount`,`type`,`description`,`date`) VALUES (?,?,?,?,?,?,?,?,?)")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, userID, req.CategoryID, req.Amount, req.Type, req.Description, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	result, err := suite.service.CreateTransaction(userID, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), req.Amount, result.Amount)
	assert.Equal(suite.T(), req.Type, result.Type)
	assert.Equal(suite.T(), "Salary", result.Category)
}

func (suite *TransactionServiceTestSuite) TestUpdateTransaction() {
	userID := uint(1)
	transactionID := uint(1)
	now := time.Now()
	req := request.UpdateTransactionRequest{
		CategoryID:  1,
		Amount:      1500.0,
		Type:        "income",
		Description: "Updated Salary",
		Date:        "2025-01-29",
	}

	// Mock get existing transaction
	txRows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"user_id", "category_id", "amount", "type",
		"description", "date",
	}).AddRow(transactionID, now, now, nil, userID, 1, 1000.0, "income", "Salary", now)

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE (id = ? AND user_id = ?) AND `transactions`.`deleted_at` IS NULL ORDER BY `transactions`.`id` LIMIT ?")).
		WithArgs(transactionID, userID, 1).
		WillReturnRows(txRows)

	// Mock category check
	categoryRows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"user_id", "name",
	}).AddRow(1, now, now, nil, userID, "Salary")

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? AND `categories`.`deleted_at` IS NULL ORDER BY `categories`.`id` LIMIT ?")).
		WithArgs(req.CategoryID, 1).
		WillReturnRows(categoryRows)

	// Mock update
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `transactions` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`user_id`=?,`category_id`=?,`amount`=?,`type`=?,`description`=?,`date`=? WHERE `transactions`.`deleted_at` IS NULL AND `id` = ?")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, userID, req.CategoryID, req.Amount, req.Type, req.Description, sqlmock.AnyArg(), transactionID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	result, err := suite.service.UpdateTransaction(userID, transactionID, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), req.Amount, result.Amount)
	assert.Equal(suite.T(), req.Description, result.Description)
}

func (suite *TransactionServiceTestSuite) TestDeleteTransaction() {
	userID := uint(1)
	transactionID := uint(1)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `transactions` SET `deleted_at`=? WHERE (id = ? AND user_id = ?) AND `transactions`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), transactionID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.service.DeleteTransaction(userID, transactionID)
	assert.NoError(suite.T(), err)
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_CategoryNotFound() {
	userID := uint(1)
	req := request.CreateTransactionRequest{
		CategoryID:  999,
		Amount:      1000.0,
		Type:        "income",
		Description: "Salary",
		Date:        "2025-01-29",
	}

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? AND `categories`.`deleted_at` IS NULL ORDER BY `categories`.`id` LIMIT ?")).
		WithArgs(req.CategoryID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := suite.service.CreateTransaction(userID, req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "category not found")
}

func TestTransactionServiceSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}
