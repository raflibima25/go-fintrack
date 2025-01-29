package service

import (
	"errors"
	"go-fintrack/internal/payload/entity"
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

	if err := s.DB.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, errors.New("failed to get all category")
	}

	// transform ke response format
	categoryResponse := make([]response.CategoryResponse, len(categories))

	for i, category := range categories {
		deletedAt := category.DeletedAt.Time
		var deletedAtPtr time.Time
		if category.DeletedAt.Valid {
			deletedAtPtr = deletedAt
		}

		categoryResponse[i] = response.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			UserID:    category.UserID,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
			DeletedAt: deletedAtPtr,
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

func (s *CategoryService) CreateCategory(name string, userID uint) (*response.CategoryResponse, error) {
	nameToLower := strings.ToLower(strings.TrimSpace(name))

	// check existing
	var existingCategory entity.Category
	if err := s.DB.Where("LOWER(name) = ? AND user_id = ?", nameToLower, userID).First(&existingCategory).Error; err == nil {
		return nil, errors.New("category name already exists")
	}

	// create category
	newCategory := entity.Category{
		UserID: userID,
		Name:   nameToLower,
	}

	if err := s.DB.Create(&newCategory).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return &response.CategoryResponse{
		ID:        newCategory.ID,
		Name:      newCategory.Name,
		UserID:    newCategory.UserID,
		CreatedAt: newCategory.CreatedAt,
		UpdatedAt: newCategory.UpdatedAt,
	}, nil
}

func (s *CategoryService) UpdateCategory(categoryID uint, userID uint, name string) (*response.CategoryResponse, error) {
	nameToLower := strings.ToLower(strings.TrimSpace(name))

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
	if err := s.DB.Save(&category).Error; err != nil {
		return nil, errors.New("failed to update category")
	}

	return &response.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		UserID:    category.UserID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
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
