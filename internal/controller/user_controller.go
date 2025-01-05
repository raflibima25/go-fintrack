package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-manajemen-keuangan/internal/payload/request"
	"go-manajemen-keuangan/internal/payload/response"
	"go-manajemen-keuangan/internal/service"
	"net/http"
	"time"
)

type UserController struct {
	UserService *service.UserService
}

func (c *UserController) RegisterHandler(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid input",
			Data:            err.Error(),
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Password dan Confirm password do not match",
			Data:            nil,
		})
		return
	}

	err := c.UserService.RegisterUser(req.Name, req.Email, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
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
		logrus.Errorf("Invalid input: %v", err) // debug
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid input",
			Data:            err.Error(),
		})
		return
	}

	logrus.Infof("Login attempt for: %s:", loginPayload.EmailOrUsername) // debug

	// proses login
	token, user, err := c.UserService.Login(loginPayload.EmailOrUsername, loginPayload.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Email, username atau password salah",
			Data:            nil,
		})
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
