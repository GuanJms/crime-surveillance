package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenExpiry = time.Minute * 20
)

var secret = []byte(os.Getenv("SECRET"))

type CustomClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateToken(id string, roles []string) (string, error) {
	claims := CustomClaims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiry)),
			Issuer:    "auth-service",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
