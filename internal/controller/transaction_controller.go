package controller

import (
	"fmt"
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/payload/response"
	"go-fintrack/internal/service"
	"go-fintrack/internal/utility"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	TransactionService *service.TransactionService
}

// GetAllTransactionsHandler godoc
// @Summary 	Get all transactions
// @Description Get all transactions with filter and pagination
// @Tags 		transactions
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Param 		start_date 	query 	string 	false 	"Start date (YYYY-MM-DD)"
// @Param 		end_date 	query 	string 	false 	"End date (YYYY-MM-DD)"
// @Param 		category_id	query 	int 	false 	"Category ID"
// @Param 		type 		query 	string 	false 	"Transaction type (income/expense)"
// @Param 		page 		query 	int 	false 	"Page number"
// @Param 		limit 		query 	int 	false 	"Limit per page"
// @Success 	200 {object} response.SuccessResponse{data=response.TransactionListResponse}
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/transaction [get]
func (c *TransactionController) GetTransactionHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	var filter request.TransactionFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logrus.Errorf("Error binding query params: %v", err)
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid filter parameters: "+err.Error(), nil)
		return
	}

	logrus.Infof("Received filter: %+v", filter) // debug

	transactions, err := c.TransactionService.GetTransactionByUser(userID, filter)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get transactions successful",
		Data:            transactions,
	})
}

// CreateTransactionHandler godoc
// @Summary 	Get new transaction
// @Description Get new transaction
// @Tags 		transactions
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Param 		request body request.CreateTransactionRequest true "Transaction data"
// @Success 	201 {object} response.SuccessResponse{data=response.TransactionResponse}
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/transaction [post]
func (c *TransactionController) CreateTransactionHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	var req request.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.ValidationErrorResponse(ctx, err)
		return
	}

	transaction, err := c.TransactionService.CreateTransaction(userID, req)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Transaction created",
		Data:            transaction,
	})
}

// UpdateTransactionHandler godoc
// @Summary 	Update transaction
// @Description Update transaction by ID
// @Tags 		transactions
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Param 		id path int true "Transaction ID"
// @Param 		request body request.UpdateTransactionRequest true "Transaction data"
// @Success 	201 {object} response.SuccessResponse{data=response.TransactionResponse}
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/transaction/{id} [put]
func (c *TransactionController) UpdateTransactionHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	transactionID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid transaction ID", nil)
		return
	}

	var req request.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	transaction, err := c.TransactionService.UpdateTransaction(userID, uint(transactionID), req)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Transaction updated",
		Data:            transaction,
	})
}

// DeleteTransactionHandler godoc
// @Summary 	Delete transaction
// @Description Delete transaction by ID
// @Tags 		transactions
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Param 		id path int true "Transaction ID"
// @Success 	200 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/transaction/{id} [delete]
func (c *TransactionController) DeleteTransactionHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	transactionID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid transaction ID", nil)
		return
	}

	if err := c.TransactionService.DeleteTransaction(userID, uint(transactionID)); err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Transaction deleted",
		Data:            nil,
	})
}

// ExportTransactionsExcelHandler godoc
// @Summary 	Export transactions to Excel
// @Description Export transactions to Excel file
// @Tags 		transactions
// @Accept 		json
// @Produce 	application/octet-stream
// @Security 	BearerAuth
// @Param 		start_date 	query 	string 	false 	"Start date (YYYY-MM-DD)"
// @Param 		end_date 	query 	string 	false 	"End date (YYYY-MM-DD)"
// @Param 		category_id	query 	int 	false 	"Category ID"
// @Param 		type 		query 	string 	false 	"Transaction type (income/expense)"
// @Success 	200 {file} file "Excel file download"
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/transaction/export [get]
func (c *TransactionController) ExportTransactionsExcelHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		logrus.Errorf("Error getting user ID: %v", err)
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	var filter request.TransactionFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logrus.Errorf("Error binding query params: %v", err)
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid filter parameters: "+err.Error(), nil)
		return
	}

	buffer, err := c.TransactionService.ExportTransactionsExcel(userID, filter)
	if err != nil {
		logrus.Errorf("Error exporting transactions: %v", err)
		utility.InternalServerErrorResponse(ctx, "Failed while export Excel", err)
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
