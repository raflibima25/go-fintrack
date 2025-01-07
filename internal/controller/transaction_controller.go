package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-manajemen-keuangan/internal/payload/request"
	"go-manajemen-keuangan/internal/payload/response"
	"go-manajemen-keuangan/internal/service"
	"go-manajemen-keuangan/internal/utility"
	"net/http"
	"strconv"
	"time"
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

	var filter request.TransactionFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logrus.Errorf("Error binding query params: %v", err)
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid filter parameters: " + err.Error(),
			Data:            nil,
		})
		return
	}

	logrus.Infof("Received filter: %+v", filter) // debug

	transactions, err := c.TransactionService.GetTransactionByUser(userID, filter)
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

func (c *TransactionController) UpdateTransactionHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
		return
	}

	transactionID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid transaction ID",
			Data:            nil,
		})
		return
	}

	var req request.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid input",
			Data:            nil,
		})
		return
	}

	transaction, err := c.TransactionService.UpdateTransaction(userID, uint(transactionID), req)
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
		ResponseMessage: "Transaction updated",
		Data:            transaction,
	})
}

func (c *TransactionController) DeleteTransactionHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
		return
	}

	transactionID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid transaction ID",
			Data:            nil,
		})
		return
	}

	if err := c.TransactionService.DeleteTransaction(userID, uint(transactionID)); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "Transaction deleted",
		Data:            nil,
	})
}

func (c *TransactionController) ExportTransactionsExcelHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		logrus.Errorf("Error getting user ID: %v", err)
		ctx.JSON(http.StatusUnauthorized, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: err.Error(),
			Data:            nil,
		})
		return
	}

	var filter request.TransactionFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logrus.Errorf("Error binding query params: %v", err)
		ctx.JSON(http.StatusBadRequest, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid filter parameters: " + err.Error(),
			Data:            nil,
		})
		return
	}

	buffer, err := c.TransactionService.ExportTransactionsExcel(userID, filter)
	if err != nil {
		logrus.Errorf("Error exporting transactions: %v", err)
		ctx.JSON(http.StatusInternalServerError, response.ApiResponse{
			ResponseStatus:  false,
			ResponseMessage: "Failed to export transactions: " + err.Error(),
			Data:            nil,
		})
		return
	}

	// Set response headers
	filename := fmt.Sprintf("transactions_%s.xlsx", time.Now().Format("20060102"))
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")

	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}
