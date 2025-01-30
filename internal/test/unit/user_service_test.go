package unit

import (
	"errors"
	"go-fintrack/internal/payload/entity"
	"go-fintrack/internal/service"
	"go-fintrack/internal/utility"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	// Set JWT secret untuk testing
	os.Setenv("JWT_SECRET", "test-secret-key")

	// Create SQL mock
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	// Create GORM DB with the SQL mock
	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	return db, mock
}

func TestRegisterUser(t *testing.T) {
	db, mock := setupTestDB(t)

	tests := []struct {
		name          string
		inputName     string
		inputEmail    string
		inputUsername string
		inputPassword string
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name:          "Success Registration",
			inputName:     "Test User",
			inputEmail:    "test@example.com",
			inputUsername: "testUser",
			inputPassword: "Password123",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Check existing user - perhatikan GORM menambahkan LIMIT 1
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WithArgs("test@example.com", "testUser", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				// Expect insert with all fields
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs(
						sqlmock.AnyArg(),   // created_at
						sqlmock.AnyArg(),   // updated_at
						nil,                // deleted_at
						"Test User",        // name
						"test@example.com", // email
						"testUser",         // username
						sqlmock.AnyArg(),   // password (hashed)
						false,              // is_admin
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:          "Existing user",
			inputName:     "Test User",
			inputEmail:    "existing@example.com",
			inputUsername: "existingUser",
			inputPassword: "Password123",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WithArgs("existing@example.com", "existingUser", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "username"}).
						AddRow(1, "Test User", "existing@example.com", "existingUser"))
			},
			expectedError: errors.New("user with this email or username already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock(mock)

			userService := &service.UserService{
				DB: db,
			}

			// Execute
			err := userService.RegisterUser(
				tt.inputName,
				tt.inputEmail,
				tt.inputUsername,
				tt.inputPassword,
			)

			// Assert
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	db, mock := setupTestDB(t)

	defer func() {
		os.Unsetenv("JWT_SECRET")
	}()

	hashedPassword, _ := utility.HashPassword("Password123")

	tests := []struct {
		name            string
		inputCredential string
		inputPassword   string
		setupMock       func(sqlmock.Sqlmock)
		expectedUser    *entity.User
		expectedError   error
	}{
		{
			name:            "Success Login with Email",
			inputCredential: "test@example.com",
			inputPassword:   "Password123",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Tambahkan regex untuk query yang lebih spesifik
				mock.ExpectQuery(`^SELECT \* FROM `+"`users`"+` WHERE .+$`).
					WithArgs("test@example.com", "test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "created_at", "updated_at", "deleted_at",
						"name", "email", "username", "password", "is_admin",
					}).AddRow(
						1, time.Now(), time.Now(), nil,
						"Test User", "test@example.com", "testUser",
						hashedPassword, false,
					))
			},
			expectedUser: &entity.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:     "Test User",
				Email:    "test@example.com",
				Username: "testUser",
				Password: hashedPassword,
				IsAdmin:  false,
			},
			expectedError: nil,
		},
		{
			name:            "Invalid Credentials",
			inputCredential: "wrong@example.com",
			inputPassword:   "Password123",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`^SELECT \* FROM `+"`users`"+` WHERE .+$`).
					WithArgs("wrong@example.com", "wrong@example.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: service.ErrInvalidCredentials,
		},
		{
			name:            "Invalid Password",
			inputCredential: "test@example.com",
			inputPassword:   "WrongPass123",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`^SELECT \* FROM `+"`users`"+` WHERE .+$`).
					WithArgs("test@example.com", "test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "created_at", "updated_at", "deleted_at",
						"name", "email", "username", "password", "is_admin",
					}).AddRow(
						1, time.Now(), time.Now(), nil,
						"Test User", "test@example.com", "testUser",
						hashedPassword, false,
					))
			},
			expectedUser:  nil,
			expectedError: service.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock(mock)

			userService := &service.UserService{
				DB: db,
			}

			// Execute
			token, user, err := userService.Login(tt.inputCredential, tt.inputPassword)

			// Assert
			if tt.expectedError != nil {
				if tt.expectedError == service.ErrInvalidCredentials {
					assert.ErrorIs(t, err, service.ErrInvalidCredentials)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
				assert.Empty(t, token)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err, "Should not have error")
				assert.NotEmpty(t, token, "Token should not be empty")
				if tt.expectedUser != nil {
					assert.NotNil(t, user, "User should not be nil")
					assert.Equal(t, tt.expectedUser.ID, user.ID)
					assert.Equal(t, tt.expectedUser.Email, user.Email)
					assert.Equal(t, tt.expectedUser.Username, user.Username)
					assert.Equal(t, tt.expectedUser.Password, user.Password)
					assert.Equal(t, tt.expectedUser.IsAdmin, user.IsAdmin)
				}
			}

			// Verify all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
