package jwtutil

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(os.Getenv("jwtSecret"))

type Claims struct {
	Account string `json:"account"`
	jwt.StandardClaims
}

func GenerateToken(account string) (string, error) {
	now := time.Now()
	claims := Claims{
		Account: account,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(24 * time.Hour).Unix(),
			IssuedAt:  now.Unix(),
			Subject:   account,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(token string) (*Claims, bool) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, true
	} else {
		return nil, false
	}
}
