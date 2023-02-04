package distro

import "github.com/google/uuid"

type Repository interface {
	GetDistros() ([]Distro, error)
	GetDistroById(id uuid.UUID) (Distro, error)
	SetDistro(d Distro) error
	DeleteDistroById(id uuid.UUID) error
}
