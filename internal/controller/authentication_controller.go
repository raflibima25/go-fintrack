package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-fintrack/config"
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/payload/response"
	"go-fintrack/internal/service"
	"go-fintrack/internal/utility"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	UserService *service.UserService
}

// RegisterHandler godoc
// @Summary 	Register new user
// @Description Register new user with name, username, email, password and confirm password
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param 		request body request.RegisterRequest true "Register credentials"
// @Success 	200 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Router 		/auth/register [post]
func (c *UserController) RegisterHandler(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utility.FormatValidationError(validationErrors)
			messages := make([]string, len(formattedErrors))
			for i, err := range formattedErrors {
				messages[i] = utility.GetReadableErrorMessage(err)
			}
			utility.ErrorResponse(ctx, http.StatusBadRequest, messages[0], []response.ErrorDetail{
				{
					Field:   "validation",
					Message: messages,
				},
			})
			return
		}

		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid input format", nil)
		return
	}

	if req.Password != req.ConfirmPassword {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Password dan Confirm password do not match", nil)
		return
	}

	err := c.UserService.RegisterUser(req.Name, req.Email, req.Username, req.Password)
	fmt.Println("err register", err)
	if err != nil {
		switch err {
		case service.ErrWeakPassword:
			utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		case service.ErrUserExists:
			utility.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
		default:
			utility.InternalServerErrorResponse(ctx, "Failed register user", err)
		}
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "User registered",
		Data: request.RegisterRequest{
			Name:     req.Name,
			Username: req.Username,
			Email:    req.Email,
		},
	})
}

// LoginHandler godoc
// @Summary 	Login user
// @Description Login user with email/username and password
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param 		request body request.LoginRequest true "Login credentials"
// @Success 	200 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/auth/login [post]
func (c *UserController) LoginHandler(ctx *gin.Context) {
	var loginPayload request.LoginRequest

	if err := ctx.ShouldBindJSON(&loginPayload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utility.FormatValidationError(validationErrors)
			messages := make([]string, len(formattedErrors))
			for i, err := range formattedErrors {
				messages[i] = utility.GetReadableErrorMessage(err)
			}
			utility.ErrorResponse(ctx, http.StatusBadRequest, messages[0], []response.ErrorDetail{
				{
					Field:   "validation",
					Message: messages,
				},
			})
			return
		}

		utility.ErrorResponse(ctx, http.StatusBadRequest, "invalid input format", nil)
		return
	}

	// proses login
	token, user, err := c.UserService.Login(loginPayload.EmailOrUsername, loginPayload.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			utility.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid email/username or password", nil)
		default:
			utility.InternalServerErrorResponse(ctx, "Failed to login", err)
		}
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Login successful",
		Data: response.LoginResponse{
			Name:        user.Name,
			AccessToken: token,
			Expiration:  time.Now().Add(24 * time.Hour),
			IsAdmin:     user.IsAdmin,
		},
	})
}

func (c *UserController) GoogleLogin(ctx *gin.Context) {
	if config.GoogleOauthConfig == nil {
		utility.InternalServerErrorResponse(ctx, "Google OAuth config is not initialized", errors.New("oauth config is nil"))
		return
	}

	state := utility.GenerateRandomString(32)
	url := config.GoogleOauthConfig.AuthCodeURL(state)

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Successfully generated Google login URL",
		Data: gin.H{
			"redirect_url": url,
			"state":        state,
		},
	})
}

func (c *UserController) GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" || state == "" {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request parameters", []response.ErrorDetail{
			{
				Field:   "code",
				Message: "Authorization code is required",
			},
			{
				Field:   "state",
				Message: "Authorization state is required",
			},
		})
		return
	}

	userChan := make(chan *request.GoogleUser)
	errorChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)

	// timeout context
	ctx.Request = ctx.Request.WithContext(
		utility.ContextWithTimeout(ctx, 10*time.Second),
	)

	// goroutine for fetching user data
	go func() {
		defer wg.Done()

		// exchange code with token
		token, err := config.GoogleOauthConfig.Exchange(ctx, code)
		if err != nil {
			errorChan <- err
			return
		}

		// create http client with token
		client := config.GoogleOauthConfig.Client(ctx, token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		// read response body
		userData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorChan <- err
			return
		}

		// unmarshal data user
		var googleUser request.GoogleUser
		if err := json.Unmarshal(userData, &googleUser); err != nil {
			errorChan <- err
			return
		}

		userChan <- &googleUser
	}()

	// goroutine for cleanup
	go func() {
		wg.Wait()
		close(userChan)
		close(errorChan)
	}()

	// handle response using select
	select {
	case user := <-userChan:
		// save or delete user data ke database
		dbUser, err := c.UserService.UpsertGoogleUser(ctx, user)
		if err != nil {
			utility.InternalServerErrorResponse(ctx, "Failed to process user data", err)
			return
		}

		token, err := utility.GenerateJWT(dbUser.ID, dbUser.Username, dbUser.IsAdmin)
		if err != nil {
			utility.InternalServerErrorResponse(ctx, "Failed to generate JWT", err)
		}

		ctx.JSON(http.StatusOK, response.SuccessResponse{
			ResponseStatus:  true,
			ResponseMessage: "Successfully authenticated with Google",
			Data: gin.H{
				"access_token": token,
				"is_admin":     false,
				"user":         user,
			},
		})

	case err := <-errorChan:
		utility.InternalServerErrorResponse(ctx, "Authentication failed", err)

	case <-ctx.Request.Context().Done():
		utility.ErrorResponse(ctx, http.StatusGatewayTimeout, "Request timeout", []response.ErrorDetail{
			{
				Field:   "timeout",
				Message: "The request took too long to process",
			},
		})
	}
}
