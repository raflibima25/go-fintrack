package service

import (
	"errors"
	"fmt"
	"go-fintrack/internal/payload/entity"
	"go-fintrack/internal/utility"
	"regexp"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

var (
	ErrUserExists         = errors.New("user with this email or username already exists")
	ErrInvalidCredentials = errors.New("invalid email/username or password")
	ErrWeakPassword       = errors.New("password must contain at least one uppercase letter, one lowercase letter, one number")
)

func (s *UserService) validatePassword(password string) error {
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return ErrWeakPassword
	}
	return nil
}

func (s *UserService) RegisterUser(name, email, username, password string) error {
	// validate password
	if err := s.validatePassword(password); err != nil {
		return err
	}

	// cek email or username exists
	var existingUser entity.User
	if err := s.DB.Where("email = ? OR username = ?", email, username).First(&existingUser).Error; err == nil {
		return ErrUserExists
	}

	// hash pass
	hashedPassword, err := utility.HashPassword(password)
	if err != nil {
		return errors.New("internal server error during registration")
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
		return fmt.Errorf("internal server error during registration")
	}

	return nil
}

func (s *UserService) Login(emailOrUsername, password string) (string, *entity.User, error) {
	var user entity.User

	if err := s.DB.Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("internal server error during login")
	}

	if err := utility.CompareHashAndPassword(user.Password, password); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// generate jwt
	token, err := utility.GenerateJWT(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return "", nil, fmt.Errorf("internal server error during login")
	}

	return token, &user, nil
}
