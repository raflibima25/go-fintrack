package integration

import (
	"bytes"
	"encoding/json"
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/payload/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	// Setup
	ts := setupTestServer(t)
	defer cleanupDatabase(ts.DB)

	tests := []struct {
		name           string
		payload        request.RegisterRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful_registration",
			payload: request.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        "Password123",
				ConfirmPassword: "Password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "password_mismatch",
			payload: request.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        "Password123",
				ConfirmPassword: "password124",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password and confirm password must match",
		},
		{
			name: "duplicate_email",
			payload: request.RegisterRequest{
				Name:            "Test User 2",
				Email:           "test@example.com", // sama dengan test case pertama
				Username:        "testuser2",
				Password:        "Password123",
				ConfirmPassword: "Password123",
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "user with this email or username already exists",
		},
		{
			name: "empty_name",
			payload: request.RegisterRequest{
				Name:            "",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        "Password123",
				ConfirmPassword: "Password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name is required",
		},
		{
			name: "invalid_email_format",
			payload: request.RegisterRequest{
				Name:            "Test User",
				Email:           "invalid-email",
				Username:        "testuser",
				Password:        "Password123",
				ConfirmPassword: "Password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid email format",
		},
		{
			name: "short_password",
			payload: request.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        "123",
				ConfirmPassword: "123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password must be at least 8 characters",
		},
		{
			name: "username_with_special_chars",
			payload: request.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Username:        "test@user",
				Password:        "Password123",
				ConfirmPassword: "Password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "username can only contain letters and numbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert payload to JSON
			jsonData, _ := json.Marshal(tt.payload)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			ts.Router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response response.SuccessResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.False(t, response.ResponseStatus)
				assert.Equal(t, tt.expectedError, response.ResponseMessage)
			} else {
				assert.True(t, response.ResponseStatus)
				assert.Equal(t, "User registered", response.ResponseMessage)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	// Setup
	ts := setupTestServer(t)
	defer cleanupDatabase(ts.DB)

	// Register a user first
	user := request.RegisterRequest{
		Name:            "Test User",
		Email:           "test@example.com",
		Username:        "testuser",
		Password:        "Password123",
		ConfirmPassword: "Password123",
	}
	if err := ts.UserService.RegisterUser(user.Name, user.Email, user.Username, user.Password); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		payload        request.LoginRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful_login_with_email",
			payload: request.LoginRequest{
				EmailOrUsername: "test@example.com",
				Password:        "Password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "successful_login_with_username",
			payload: request.LoginRequest{
				EmailOrUsername: "testuser",
				Password:        "Password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_credentials",
			payload: request.LoginRequest{
				EmailOrUsername: "test@example.com",
				Password:        "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid email/username or password",
		},
		{
			name: "empty_credentials",
			payload: request.LoginRequest{
				EmailOrUsername: "",
				Password:        "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email, username or password is required",
		},
		{
			name: "short_password",
			payload: request.LoginRequest{
				EmailOrUsername: "test@example.com",
				Password:        "123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password must be at least 8 characters",
		},
		{
			name: "non_existent_user",
			payload: request.LoginRequest{
				EmailOrUsername: "nonexistent@example.com",
				Password:        "Password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid email/username or password",
		},
	}

	// Eksekusi setiap test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert payload ke JSON
			jsonData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			// Buat request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			ts.Router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response response.SuccessResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.False(t, response.ResponseStatus)
				assert.Equal(t, tt.expectedError, response.ResponseMessage)
			} else {
				assert.True(t, response.ResponseStatus)
				assert.Equal(t, "Login successful", response.ResponseMessage)

				// Untuk successful login, cek apakah token ada
				loginResponse, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, loginResponse["access_token"])
				assert.NotEmpty(t, loginResponse["name"])
				assert.NotEmpty(t, loginResponse["expiration"])
			}
		})
	}
}
