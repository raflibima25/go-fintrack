package service

import (
	"errors"
	"go-fintrack/internal/payload/entity"
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/payload/response"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CategoryService struct {
	DB *gorm.DB
}

func (s *CategoryService) GetCategories(userID uint) ([]response.CategoryResponse, error) {
	var categories []entity.Category

	var totalTransactions int64
	s.DB.Model(&entity.Transaction{}).Where("user_id = ?", userID).Count(&totalTransactions)

	if err := s.DB.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, errors.New("failed to get all category")
	}

	// transform ke response format
	categoryResponse := make([]response.CategoryResponse, len(categories))

	for i, category := range categories {
		// hitung usage count
		var usageCount int64
		s.DB.Model(&entity.Transaction{}).Where("category_id = ?", category.ID).Count(&usageCount)

		// hitung usage percentage
		var usagePercentage float64
		if totalTransactions > 0 {
			usagePercentage = float64(usageCount) / float64(totalTransactions) * 100
		}

		deletedAt := category.DeletedAt.Time
		var deletedAtPtr time.Time
		if category.DeletedAt.Valid {
			deletedAtPtr = deletedAt
		}

		categoryResponse[i] = response.CategoryResponse{
			ID:              category.ID,
			Name:            category.Name,
			Color:           category.Color,
			IconColor:       category.IconColor,
			UsageCount:      usageCount,
			UsagePercentage: usagePercentage,
			UserID:          category.UserID,
			CreatedAt:       category.CreatedAt,
			UpdatedAt:       category.UpdatedAt,
			DeletedAt:       deletedAtPtr,
		}
	}

	return categoryResponse, nil
}

func (s *CategoryService) GetCategoryByID(categoryID uint, userID uint) (*response.CategoryResponse, error) {
	var category entity.Category
	if err := s.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to get category")
	}

	deletedAt := category.DeletedAt.Time
	var deletedAtPtr time.Time
	if category.DeletedAt.Valid {
		deletedAtPtr = deletedAt
	}

	return &response.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		UserID:    category.UserID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
		DeletedAt: deletedAtPtr,
	}, nil
}

func (s *CategoryService) CreateCategory(req *request.CategoryRequest, userID uint) (*response.CategoryResponse, error) {
	nameToLower := strings.ToLower(strings.TrimSpace(req.Name))

	// check existing
	var existingCategory entity.Category
	if err := s.DB.Where("LOWER(name) = ? AND user_id = ?", nameToLower, userID).First(&existingCategory).Error; err == nil {
		return nil, errors.New("category name already exists")
	}

	// create category
	newCategory := entity.Category{
		UserID:    userID,
		Name:      nameToLower,
		Color:     req.Color,
		IconColor: req.IconColor,
	}

	if err := s.DB.Create(&newCategory).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return &response.CategoryResponse{
		ID:        newCategory.ID,
		Name:      newCategory.Name,
		Color:     newCategory.Color,
		IconColor: newCategory.IconColor,
		UserID:    newCategory.UserID,
		CreatedAt: newCategory.CreatedAt,
		UpdatedAt: newCategory.UpdatedAt,
	}, nil
}

func (s *CategoryService) UpdateCategory(categoryID uint, userID uint, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error) {
	nameToLower := strings.ToLower(strings.TrimSpace(req.Name))

	// check category
	var category entity.Category
	if err := s.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to get category")
	}

	// check nama baru setelah update already exists
	var existingCategory entity.Category
	if err := s.DB.Where("LOWER(name) = ? AND user_id = ? AND id != ?", nameToLower, userID, categoryID).First(&existingCategory).Error; err == nil {
		return nil, errors.New("category name already exists")
	}

	// update category
	category.Name = nameToLower
	if req.Color != "" {
		category.Color = req.Color
	}

	if req.IconColor != "" {
		category.IconColor = req.IconColor
	}

	if err := s.DB.Save(&category).Error; err != nil {
		return nil, errors.New("failed to update category")
	}

	// Hitung usage count dan percentage untuk response
	var totalTransactions int64
	s.DB.Model(&entity.Transaction{}).Where("user_id = ?", userID).Count(&totalTransactions)

	var usageCount int64
	s.DB.Model(&entity.Transaction{}).Where("category_id = ?", category.ID).Count(&usageCount)

	var usagePercentage float64
	if totalTransactions > 0 {
		usagePercentage = float64(usageCount) / float64(totalTransactions) * 100
	}

	return &response.CategoryResponse{
		ID:              category.ID,
		Name:            category.Name,
		Color:           category.Color,
		IconColor:       category.IconColor,
		UsageCount:      usageCount,
		UsagePercentage: usagePercentage,
		UserID:          category.UserID,
		CreatedAt:       category.CreatedAt,
		UpdatedAt:       category.UpdatedAt,
	}, nil
}

func (s *CategoryService) DeleteCategory(categoryID uint, userID uint) error {
	result := s.DB.Where("id = ? AND user_id = ?", categoryID, userID).Delete(&entity.Category{})
	if result.Error != nil {
		return errors.New("failed to delete category")
	}

	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}
