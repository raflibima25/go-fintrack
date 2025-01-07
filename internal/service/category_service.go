package service

import (
	"errors"
	"go-manajemen-keuangan/internal/payload/entity"
	"go-manajemen-keuangan/internal/payload/response"
	"gorm.io/gorm"
	"strings"
	"time"
)

type CategoryService struct {
	DB *gorm.DB
}

func (s *CategoryService) GetCategories() ([]response.CategoryResponse, error) {
	var categories []entity.Category
	if err := s.DB.Find(&categories).Error; err != nil {
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
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
			DeletedAt: deletedAtPtr,
		}
	}

	return categoryResponse, nil
}

func (s *CategoryService) GetCategoryByID(categoryID uint) (*response.CategoryResponse, error) {
	var category entity.Category
	if err := s.DB.First(&category, categoryID).Error; err != nil {
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
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
		DeletedAt: deletedAtPtr,
	}, nil
}

func (s *CategoryService) CreateCategory(name string) (*response.CategoryResponse, error) {
	nameToLower := strings.ToLower(strings.TrimSpace(name))

	// check existing
	var existingCategory entity.Category
	if err := s.DB.Where("LOWER(name) = ?", nameToLower).First(&existingCategory).Error; err == nil {
		return nil, errors.New("category name already exists")
	}

	// create category
	newCategory := entity.Category{
		Name: nameToLower,
	}

	if err := s.DB.Create(&newCategory).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return &response.CategoryResponse{
		ID:        newCategory.ID,
		Name:      newCategory.Name,
		CreatedAt: newCategory.CreatedAt,
		UpdatedAt: newCategory.UpdatedAt,
	}, nil
}

func (s *CategoryService) UpdateCategory(categoryID uint, name string) (*response.CategoryResponse, error) {
	nameToLower := strings.ToLower(strings.TrimSpace(name))

	// check category
	var category entity.Category
	if err := s.DB.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to get category")
	}

	// check nama baru setelah update already exists
	var existingCategory entity.Category
	if err := s.DB.Where("LOWER(name) = ? AND id != ?", nameToLower, categoryID).First(&existingCategory).Error; err == nil {
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
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (s *CategoryService) DeleteCategory(categoryID uint) error {
	result := s.DB.Delete(&entity.Category{}, categoryID)
	if result.Error != nil {
		return errors.New("failed to delete category")
	}

	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}
