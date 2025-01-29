package unit

import (
	"errors"
	"go-fintrack/internal/payload/entity"
	"go-fintrack/internal/service"
	"go-fintrack/internal/utility"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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
			inputPassword: "password123",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Check existing user - perhatikan GORM menambahkan LIMIT 1
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WithArgs("test@example.com", "testUser", 1). // tambahkan argument LIMIT
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
			inputPassword: "password123",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WithArgs("existing@example.com", "existingUser", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "username"}).
						AddRow(1, "Test User", "existing@example.com", "existingUser"))
			},
			expectedError: errors.New("email or username already exists"),
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
	hashedPassword, _ := utility.HashPassword("password123")

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
			inputPassword:   "password123",
			setupMock: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "name", "email", "username", "password", "is_admin"}
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WithArgs("test@example.com", "test@example.com", 1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "Test User", "test@example.com", "testUser", hashedPassword, false))
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
			inputPassword:   "wrongpassword",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WithArgs("wrong@example.com", "wrong@example.com", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectedUser:  nil,
			expectedError: errors.New("wrong email, username or password"),
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
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Empty(t, token)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				if tt.expectedUser != nil {
					assert.Equal(t, tt.expectedUser.ID, user.ID)
					assert.Equal(t, tt.expectedUser.Username, user.Username)
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
