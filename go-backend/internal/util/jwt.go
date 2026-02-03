package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	Secret string
	Issuer string
	TTL    int64
}

func (j JWTManager) Generate(claims jwt.MapClaims) (string, error) {
	now := time.Now().Unix()
	claims["iss"] = j.Issuer
	claims["iat"] = now
	claims["exp"] = now + j.TTL

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
}

func (j JWTManager) Parse(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}
