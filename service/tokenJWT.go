package service

import (
	"server-article/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(claims jwt.MapClaims) (string, error) {
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.GetEnv("SECRET_TOKEN_JWT")))
}
