package memory

import (
	"github.com/evanebb/gobble/system"
	"github.com/google/uuid"
	"net"
)

type SystemRepository struct {
}

func NewSystemRepository() (SystemRepository, error) {
	return SystemRepository{}, nil
}

func (s SystemRepository) GetSystems() ([]system.System, error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemRepository) GetSystemByMacAddress(macAddress net.HardwareAddr) (system.System, error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemRepository) GetSystemById(id uuid.UUID) (system.System, error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemRepository) SetSystem(sys system.System) error {
	//TODO implement me
	panic("implement me")
}

func (s SystemRepository) DeleteSystemById(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
