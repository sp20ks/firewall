package repository

import "auth/internal/entity"

type UserRepository interface {
	GetUser(username string) (*entity.User, error)
	CreateUser(user *entity.User) error
}
