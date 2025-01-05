package utility

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID uint, username string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
