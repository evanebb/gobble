package memory

import "github.com/evanebb/gobble/distro"

type DistroRepository struct {
}

func NewDistroRepository() (DistroRepository, error) {
	return DistroRepository{}, nil
}

func (d DistroRepository) GetDistros() ([]distro.Distro, error) {
	//TODO implement me
	panic("implement me")
}

func (d DistroRepository) GetDistroById(id uint) (distro.Distro, error) {
	//TODO implement me
	panic("implement me")
}

func (d DistroRepository) SetDistro(distro distro.Distro) error {
	//TODO implement me
	panic("implement me")
}

func (d DistroRepository) DeleteDistroById(id uint) error {
	//TODO implement me
	panic("implement me")
}
