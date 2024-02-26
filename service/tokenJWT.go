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

func DecodeToken(t string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(t, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(utils.GetEnv("SECRET_TOKEN_JWT")), nil
	})

	if err != nil && !token.Valid {
		return nil, err
	}

	return claims, nil

}
