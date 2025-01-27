package utility

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrUserIDNotFound = errors.New("user ID not found in context")
	ErrInvalidUserID  = errors.New("invalid user ID format")
)

func GetUserIDFromContext(ctx *gin.Context) (uint, error) {
	// trace ID logging
	traceID := ctx.GetString("X-Request-ID")
	logger := logrus.WithFields(logrus.Fields{
		"trace_id": traceID,
		"function": "GetUserIDFromContext",
	})

	// convert userID to uint
	userIDVal, exist := ctx.Get("userID")
	logrus.Debugf("Retrieved UserID from context: %v, exists: %v", userIDVal, exist) // debug

	if !exist {
		return 0, ErrUserIDNotFound
	}

	switch v := userIDVal.(type) {
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case uint:
		return v, nil
	case string:
		u64, _ := strconv.ParseUint(v, 10, 32)
		return uint(u64), nil
	default:
		logger.Errorf("Unexpected user ID type: %T", userIDVal)
		return 0, ErrInvalidUserID
	}
}
