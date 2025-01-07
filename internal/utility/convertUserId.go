package utility

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetUserIDFromContext(ctx *gin.Context) (uint, error) {
	// convert userID to uint
	userIDstr, exist := ctx.Get("userID")
	logrus.Infof("UserID from context: %v, exists: %v", userIDstr, exist) // debug

	if !exist {
		return 0, errors.New("user ID not found in context")
	}

	userIDFloat, ok := userIDstr.(float64)
	logrus.Infof("Converted user ID: %v, conversion ok: %v", userIDFloat, ok) // debug
	if !ok {
		return 0, errors.New("invalid user ID")
	}

	return uint(userIDFloat), nil
}
