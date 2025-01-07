package controller

import (
	"github.com/gin-gonic/gin"
	"go-manajemen-keuangan/internal/payload/request"
	"go-manajemen-keuangan/internal/payload/response"
	"go-manajemen-keuangan/internal/service"
	"go-manajemen-keuangan/internal/utility"
	"net/http"
)

type TransactionController struct {
	TransactionService *service.TransactionService
}

func (c *TransactionController) GetTransactionHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
		return
	}

	transactions, err := c.TransactionService.GetTransactionByUser(userID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get transactions successful",
		Data:            transactions,
	})
}

func (c *TransactionController) CreateTransactionHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
		return
	}

	var req request.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid input",
			Data:            nil,
		})
		return
	}

	transaction, err := c.TransactionService.CreateTransaction(userID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "Transaction created",
		Data:            transaction,
	})
}
