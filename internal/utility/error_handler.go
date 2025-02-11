package utility

import (
	"go-fintrack/internal/payload/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Error response 4xx
func ErrorResponse(ctx *gin.Context, statusCode int, message string, errors []response.ErrorDetail) {
	ctx.JSON(statusCode, response.ErrorResponse{
		ResponseStatus:  false,
		ResponseMessage: message,
		Errors:          errors,
	})
}

// Error response 5xx
func InternalServerErrorResponse(ctx *gin.Context, message string, err error) {
	requestID, _ := ctx.Get("RequestID")

	logrus.WithFields(logrus.Fields{
		"request_id": requestID,
		"error":      err,
		"message":    message,
	}).Error("Internal server error")

	ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
		ResponseStatus:  false,
		ResponseMessage: message,
		Errors: []response.ErrorDetail{
			{
				Field:   "server",
				Message: "Please contact support with this Request ID: " + requestID.(string),
			},
		},
	})
}

// untuk error validasi
func ValidationErrorResponse(ctx *gin.Context, err error) {
	ErrorResponse(ctx, http.StatusBadRequest, "Validation error", []response.ErrorDetail{
		{
			Field:   "validation",
			Message: err.Error(),
		},
	})
}

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				requestID, _ := ctx.Get("RequestID")
				logrus.WithFields(logrus.Fields{
					"request_id": requestID,
					"panic":      r,
				}).Error("Panic recovered")

				ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error", []response.ErrorDetail{
					{
						Field:   "server",
						Message: "An unexpected error occurred",
					},
				})

				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
