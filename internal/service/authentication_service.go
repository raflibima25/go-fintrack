package service

import (
	"context"
	"errors"
	"fmt"
	"go-fintrack/internal/payload/entity"
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/utility"
	"regexp"
	"strings"
	"unicode"

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

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// cek email or username exists
		var existingUser entity.User
		if err := tx.Where("email = ? OR username = ?", email, username).First(&existingUser).Error; err == nil {
			return ErrUserExists
		} else if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("error checking user existence: %v", err)
		}

		// hash pass
		hashedPassword, err := utility.HashPassword(password)
		if err != nil {
			return fmt.Errorf("error hashing password: %v", err)
		}

		// create user
		newUser := entity.User{
			Name:     name,
			Email:    email,
			Username: username,
			Password: hashedPassword,
			Provider: "email",
		}

		// save database
		if err := tx.Create(&newUser).Error; err != nil {
			return fmt.Errorf("error creating user: %v", err)
		}

		return nil
	})

	return err
}

func (s *UserService) Login(emailOrUsername, password string) (string, *entity.User, error) {
	var user entity.User

	// find user
	if err := s.DB.Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("internal server error during login: %v", err)
	}

	// verify password
	if err := utility.CompareHashAndPassword(user.Password, password); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// generate jwt
	token, err := utility.GenerateJWT(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return "", nil, fmt.Errorf("internal server error during login: %v", err)
	}

	return token, &user, nil
}

func (s *UserService) UpsertGoogleUser(ctx context.Context, googleUser *request.GoogleUser) (*entity.User, error) {
	var user entity.User

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("email = ?", googleUser.Email).First(&user)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newUser := entity.User{
				Name:       googleUser.Name,
				Username:   strings.ToLower(strings.Join(strings.FieldsFunc(googleUser.Name, func(r rune) bool { return unicode.IsSpace(r) }), "")),
				Email:      googleUser.Email,
				IsAdmin:    false,
				Provider:   "google",
				ProfilePic: googleUser.Picture,
			}

			if err := tx.Create(&newUser).Error; err != nil {
				return fmt.Errorf("error creating user: %v", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("error checking user existence: %v", result.Error)
		} else {
			// update user jika ada
			user.Name = googleUser.Name
			user.Username = strings.ToLower(strings.Join(strings.FieldsFunc(googleUser.Name, func(r rune) bool { return unicode.IsSpace(r) }), ""))
			user.ProfilePic = googleUser.Picture
			user.Provider = "google"

			if err := tx.Save(&user).Error; err != nil {
				return fmt.Errorf("error updating user: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}
