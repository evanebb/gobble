package memory

import (
	"github.com/evanebb/gobble/distro"
	"github.com/google/uuid"
)

type DistroRepository struct {
}

func NewDistroRepository() (DistroRepository, error) {
	return DistroRepository{}, nil
}

func (d DistroRepository) GetDistros() ([]distro.Distro, error) {
	//TODO implement me
	panic("implement me")
}

func (d DistroRepository) GetDistroById(id uuid.UUID) (distro.Distro, error) {
	//TODO implement me
	panic("implement me")
}

func (d DistroRepository) SetDistro(distro distro.Distro) error {
	//TODO implement me
	panic("implement me")
}

func (d DistroRepository) DeleteDistroById(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
