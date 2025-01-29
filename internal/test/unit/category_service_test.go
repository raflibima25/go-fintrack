package unit

import (
	"database/sql"
	"go-fintrack/internal/service"
	"io"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CategoryServiceTestSuite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	service *service.CategoryService
	sqlDB   *sql.DB
}

func (suite *CategoryServiceTestSuite) SetupTest() {
	var err error
	suite.sqlDB, suite.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
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

	suite.service = &service.CategoryService{
		DB: suite.DB,
	}
}

func (suite *CategoryServiceTestSuite) TearDownTest() {
	suite.sqlDB.Close()
}

func (suite *CategoryServiceTestSuite) TestGetCategories() {
	userID := uint(1)
	now := time.Now()

	query := "SELECT * FROM `categories` WHERE user_id = ? AND `categories`.`deleted_at` IS NULL"
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "user_id"}).
		AddRow(1, now, now, nil, "food", userID).
		AddRow(2, now, now, nil, "transport", userID)

	suite.mock.ExpectQuery(query).
		WithArgs(userID).
		WillReturnRows(rows)

	categories, err := suite.service.GetCategories(userID)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), categories, 2)
	if len(categories) > 0 {
		assert.Equal(suite.T(), "food", categories[0].Name)
		assert.Equal(suite.T(), "transport", categories[1].Name)
	}
}

func (suite *CategoryServiceTestSuite) TestCreateCategory() {
	userID := uint(1)
	name := "groceries"

	// Check existing
	checkQuery := "SELECT * FROM `categories` WHERE (LOWER(name) = ? AND user_id = ?) AND `categories`.`deleted_at` IS NULL ORDER BY `categories`.`id` LIMIT ?"
	suite.mock.ExpectQuery(checkQuery).
		WithArgs("groceries", userID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Create
	suite.mock.ExpectBegin()
	createQuery := "INSERT INTO `categories` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`name`) VALUES (?,?,?,?,?)"
	suite.mock.ExpectExec(createQuery).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, userID, name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	category, err := suite.service.CreateCategory(name, userID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), category)
	if category != nil {
		assert.Equal(suite.T(), name, category.Name)
	}
}

func (suite *CategoryServiceTestSuite) TestUpdateCategory() {
	categoryID := uint(1)
	userID := uint(1)
	oldName := "food"
	newName := "updated food"
	now := time.Now()

	// Get existing category
	getQuery := "SELECT * FROM `categories` WHERE (id = ? AND user_id = ?) AND `categories`.`deleted_at` IS NULL ORDER BY `categories`.`id` LIMIT ?"
	getRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "user_id"}).
		AddRow(categoryID, now, now, nil, oldName, userID)
	suite.mock.ExpectQuery(getQuery).
		WithArgs(categoryID, userID, 1).
		WillReturnRows(getRows)

	// Check duplicate name
	checkQuery := "SELECT * FROM `categories` WHERE (LOWER(name) = ? AND user_id = ? AND id != ?) AND `categories`.`deleted_at` IS NULL ORDER BY `categories`.`id` LIMIT ?"
	suite.mock.ExpectQuery(checkQuery).
		WithArgs(newName, userID, categoryID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Update
	suite.mock.ExpectBegin()
	updateQuery := "UPDATE `categories` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`user_id`=?,`name`=? WHERE `categories`.`deleted_at` IS NULL AND `id` = ?"
	suite.mock.ExpectExec(updateQuery).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, userID, newName, categoryID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	result, err := suite.service.UpdateCategory(categoryID, userID, newName)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	if result != nil {
		assert.Equal(suite.T(), newName, result.Name)
	}
}

func (suite *CategoryServiceTestSuite) TestDeleteCategory() {
	categoryID := uint(1)
	userID := uint(1)

	// Mock soft delete
	suite.mock.ExpectBegin()
	deleteQuery := "UPDATE `categories` SET `deleted_at`=? WHERE (id = ? AND user_id = ?) AND `categories`.`deleted_at` IS NULL"
	suite.mock.ExpectExec(deleteQuery).
		WithArgs(sqlmock.AnyArg(), categoryID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.service.DeleteCategory(categoryID, userID)
	assert.NoError(suite.T(), err)
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceTestSuite))
}
