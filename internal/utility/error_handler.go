package utility

import (
	"go-fintrack/internal/payload/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Error response 4xx
func ErrorResponse(ctx *gin.Context, statusCode int, message string, details interface{}) {
	ctx.JSON(statusCode, response.ApiResponse{
		ResponseStatus:  false,
		ResponseMessage: message,
		Data:            details,
	})
}

// Error response 5xx
func ServerErrorResponse(ctx *gin.Context, err error) {
	requestID, _ := ctx.Get("RequestID")

	logrus.WithFields(logrus.Fields{
		"request_id": requestID,
		"error":      err,
	}).Error("Internal server error")

	ctx.JSON(http.StatusInternalServerError, response.ApiResponse{
		ResponseStatus:  false,
		ResponseMessage: "Internal server error",
		Data: map[string]interface{}{
			"request_id": requestID,
			"message":    "Please contact support with this Request ID",
		},
	})
}

// untuk error validasi
func ValidationErrorResponse(ctx *gin.Context, err error) {
	ErrorResponse(ctx, http.StatusBadRequest, "Validation error", err.Error())
}

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
