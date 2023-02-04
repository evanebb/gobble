package system

import (
	"github.com/google/uuid"
	"net"
)

type Repository interface {
	GetSystems() ([]System, error)
	GetSystemByMacAddress(macAddress net.HardwareAddr) (System, error)
	GetSystemById(id uuid.UUID) (System, error)
	SetSystem(s System) error
	DeleteSystemById(id uuid.UUID) error
}
