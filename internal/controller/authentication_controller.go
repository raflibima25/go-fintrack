package controller

import (
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/payload/response"
	"go-fintrack/internal/service"
	"go-fintrack/internal/utility"
	"net/http"
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
// @Success 	200 {object} response.ApiResponse
// @Failure 	400 {object} response.ApiResponse
// @Router 		/user/register [post]
func (c *UserController) RegisterHandler(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utility.FormatValidationError(validationErrors)
			messages := make([]string, len(formattedErrors))
			for i, err := range formattedErrors {
				messages[i] = utility.GetReadableErrorMessage(err)
			}
			utility.ErrorResponse(ctx, http.StatusBadRequest, messages[0], messages)
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
	if err != nil {
		switch err {
		case service.ErrWeakPassword:
			utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		case service.ErrUserExists:
			utility.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
		default:
			utility.ServerErrorResponse(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{
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
// @Success 	200 {object} response.ApiResponse
// @Failure 	400 {object} response.ApiResponse
// @Failure 	401 {object} response.ApiResponse
// @Router 		/user/login [post]
func (c *UserController) LoginHandler(ctx *gin.Context) {
	var loginPayload request.LoginRequest

	if err := ctx.ShouldBindJSON(&loginPayload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utility.FormatValidationError(validationErrors)
			messages := make([]string, len(formattedErrors))
			for i, err := range formattedErrors {
				messages[i] = utility.GetReadableErrorMessage(err)
			}
			utility.ErrorResponse(ctx, http.StatusBadRequest, messages[0], messages)
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
			utility.ServerErrorResponse(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{
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
