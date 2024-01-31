package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	expirationTime = time.Hour
	signingMethod  = jwt.SigningMethodHS512
)

type Claims struct {
	jwt.RegisteredClaims

	UserID    uint
	UserAgent string
}

func GenerateToken(signingKey []byte, userID uint, userAgent string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expirationTime)),
		},
		UserID:    userID,
		UserAgent: userAgent,
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ParseWithClaims(signingKey, tokenString string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
}
