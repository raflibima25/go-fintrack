package controller

import (
	"go-manajemen-keuangan/internal/payload/request"
	"go-manajemen-keuangan/internal/payload/response"
	"go-manajemen-keuangan/internal/service"
	"go-manajemen-keuangan/internal/utility"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *service.UserService
}

func (c *UserController) RegisterHandler(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.ValidationErrorResponse(ctx, err)
		return
	}

	if req.Password != req.ConfirmPassword {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Password dan Confirm password do not match", nil)
		return
	}

	err := c.UserService.RegisterUser(req.Name, req.Email, req.Username, req.Password)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
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

func (c *UserController) LoginHandler(ctx *gin.Context) {
	var loginPayload request.LoginRequest

	if err := ctx.ShouldBindJSON(&loginPayload); err != nil {
		utility.ValidationErrorResponse(ctx, err)
		return
	}

	// validasi input
	if loginPayload.EmailOrUsername == "" || loginPayload.Password == "" {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Email, username atau password is required", nil)
		return
	}

	if len(loginPayload.Password) < 8 {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Password must be at least 8 characters", nil)
		return
	}

	// proses login
	token, user, err := c.UserService.Login(loginPayload.EmailOrUsername, loginPayload.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			utility.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid email/username or password", nil)
			return
		}

		utility.ServerErrorResponse(ctx, err)
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
