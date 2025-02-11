package middleware

import (
	"go-fintrack/internal/payload/response"
	"go-fintrack/internal/utility"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		logrus.Infof("Auth header: %s", authHeader) // debug
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, response.SuccessResponse{
				ResponseStatus:  false,
				ResponseMessage: "Authorization header is missing",
				Data:            nil,
			})
			ctx.Abort()
			return
		}

		// ambil token setelah bearer
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		token, err := utility.ParseJWT(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, response.SuccessResponse{
				ResponseStatus:  false,
				ResponseMessage: "Invalid token",
				Data:            nil,
			})
			ctx.Abort()
			return
		}

		// menyimpan info user dari token ke dalam context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("userID", claims["sub"])
			ctx.Set("username", claims["username"])
		}

		ctx.Next()
	}
}
