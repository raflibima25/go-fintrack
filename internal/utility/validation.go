package utility

import "github.com/gin-gonic/gin"

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
