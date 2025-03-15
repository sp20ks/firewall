package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	secretKey string
	ttl       time.Duration
}

func NewAuth(key string, ttl time.Duration) *Auth {
	return &Auth{secretKey: key, ttl: ttl}
}

func (a *Auth) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * a.ttl).Unix(),
		})
	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("error while sign token: %v", err)
	}
	return tokenString, nil
}

func (a *Auth) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
