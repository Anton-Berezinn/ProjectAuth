package jwt

import (
	"context"
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	secretKey = "PitPivoHalth"
)

func NewJwt(ctx context.Context, userId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
