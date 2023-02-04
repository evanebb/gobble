package profile

import "github.com/google/uuid"

type Repository interface {
	GetProfiles() ([]Profile, error)
	GetProfileById(id uuid.UUID) (Profile, error)
	SetProfile(p Profile) error
	DeleteProfileById(id uuid.UUID) error
}
