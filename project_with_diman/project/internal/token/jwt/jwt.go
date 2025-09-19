package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	secretKey = "PitPivoHalth"
)

type MyCustomClaims struct {
	UserId  int `json:"user_id"`
	Version int `json:"version"`
	jwt.StandardClaims
}

func NewJwt(ctx context.Context, userId int, version *int) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа
	}

	if version != nil {
		claims["version"] = *version
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func DecodeJwt(ctx context.Context, tokenString string) (*MyCustomClaims, error) {
	claims := &MyCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
