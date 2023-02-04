package system

import "net"

type Repository interface {
	GetSystems() ([]System, error)
	GetSystemByMacAddress(macAddress net.HardwareAddr) (System, error)
	GetSystemById(id uint) (System, error)
	SetSystem(s System) error
	DeleteSystemById(id uint) error
}
