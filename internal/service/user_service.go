package service

import (
	"errors"
	"go-manajemen-keuangan/internal/payload/entity"
	"go-manajemen-keuangan/internal/utility"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (s *UserService) RegisterUser(name, email, username, password string) error {
	// cek email or username exists
	var existingUser entity.User
	if err := s.DB.Where("email = ? OR username = ?", email, username).First(&existingUser).Error; err == nil {
		return errors.New("email or username already exists")
	}

	// hash pass
	hashedPassword, err := utility.HashPassword(password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// create user
	newUser := entity.User{
		Name:     name,
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	// save database
	if err := s.DB.Create(&newUser).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}
