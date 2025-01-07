package utility

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "jwtsecret"
	}
	return []byte(secret)
}

func GenerateJWT(userID uint, username string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 24 jam
	}

	jwtSecret := getJWTSecret()

	fmt.Println("jwt secret:", string(jwtSecret)) // debug
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	jwtSecret := getJWTSecret()

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
