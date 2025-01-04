package controller

import (
	"github.com/gin-gonic/gin"
	"go-manajemen-keuangan/internal/model"
	"go-manajemen-keuangan/internal/service"
	"net/http"
)

type UserController struct {
	UserService *service.UserService
}

func (c UserController) RegisterHandler(ctx *gin.Context) {
	var req model.RegisterModel
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid input",
			Data:            err.Error(),
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, model.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Password dan Confirm password do not match",
			Data:            nil,
		})
		return
	}

	err := c.UserService.RegisterUser(req.Name, req.Email, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, model.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "User registered",
		Data: model.RegisterModel{
			Name:     req.Name,
			Username: req.Username,
			Email:    req.Email,
		},
	})
}
