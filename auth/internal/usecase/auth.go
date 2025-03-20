package usecase

import (
	"fmt"
	"time"

	"auth/internal/entity"
	"auth/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCase struct {
	repo      repository.UserRepository
	secretKey string
	ttl       time.Duration
}

func NewAuthUseCase(repo repository.UserRepository, key string, ttl time.Duration) *AuthUseCase {
	return &AuthUseCase{repo: repo, secretKey: key, ttl: ttl}
}

func (a *AuthUseCase) CreateUser(username, passwordHash string) error {
	user, _ := a.repo.GetUser(username)
	if user != nil {
		return fmt.Errorf("user already exists")
	}

	user = &entity.User{
		Username:  username,
		Password:  passwordHash,
		CreatedAt: time.Now(),
	}
	return a.repo.CreateUser(user)
}

func (a *AuthUseCase) Authenticate(username, password string) (string, error) {
	user, err := a.repo.GetUser(username)
	if err != nil {
		return "", fmt.Errorf("error fetching user: %w", err)
	}
	if user == nil || user.Password != password {
		return "", fmt.Errorf("invalid credentials")
	}

	return a.CreateToken(username)
}

func (a *AuthUseCase) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * a.ttl).Unix(),
	})
	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return tokenString, nil
}

func (a *AuthUseCase) GetUserByToken(tokenString string) (*entity.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if username, ok := claims["username"]; !ok {
		return nil, fmt.Errorf("error while getting claims")
	} else {
		user, _ := a.repo.GetUser(username.(string))
		if user == nil {
			return nil, fmt.Errorf("user not found")
		}
		return user, nil
	}
}
