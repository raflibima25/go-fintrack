package service

import (
	"errors"
	"go-fintrack/internal/payload/entity"
	"go-fintrack/internal/utility"

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

func (s *UserService) Login(emailOrUsername, password string) (string, *entity.User, error) {
	var user entity.User

	if err := s.DB.Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).First(&user).Error; err != nil {
		return "", nil, errors.New("wrong email, username or password")
	}

	if err := utility.CompareHashAndPassword(user.Password, password); err != nil {
		return "", nil, errors.New("wrong email, username or password")
	}

	// generate jwt
	token, err := utility.GenerateJWT(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}
