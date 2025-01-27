package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-manajemen-keuangan/internal/payload/response"
	"net/http"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, response.ApiResponse{
				ResponseStatus:  false,
				ResponseMessage: "Unauthorized",
				Data:            nil,
			})
			ctx.Abort()
			return
		}

		if mapClaims, ok := claims.(jwt.MapClaims); ok {
			isAdmin, exists := mapClaims["is_admin"]
			if !exists || !isAdmin.(bool) {
				ctx.JSON(http.StatusForbidden, response.ApiResponse{
					ResponseStatus:  false,
					ResponseMessage: "Access denied: Admin only",
					Data:            nil,
				})
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
