package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ApiUser struct {
	Id       uuid.UUID
	Name     string
	Password []byte
}

func NewApiUser(id uuid.UUID, name string, password []byte) ApiUser {
	return ApiUser{
		Id:       id,
		Name:     name,
		Password: password,
	}
}

func (u ApiUser) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password))
}
