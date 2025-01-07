package utility

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(ctx *gin.Context) (uint, error) {
	// convert userID to uint
	userIDstr, exist := ctx.Get("userID")
	if !exist {
		return 0, errors.New("user ID not found in context")
	}

	userIDFloat, ok := userIDstr.(float64)
	if !ok {
		return 0, errors.New("invalid user ID")
	}

	return uint(userIDFloat), nil
}
