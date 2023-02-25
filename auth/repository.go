package auth

import "github.com/google/uuid"

type ApiUserRepository interface {
	GetApiUsers() ([]ApiUser, error)
	GetApiUserById(id uuid.UUID) (ApiUser, error)
	GetApiUserByName(name string) (ApiUser, error)
	SetApiUser(a ApiUser) error
	DeleteApiUserById(id uuid.UUID) error
}
