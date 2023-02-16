package system

import (
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"net"
)

type System struct {
	Id               uuid.UUID
	Name             string
	Description      string
	Profile          uuid.UUID
	Mac              net.HardwareAddr
	KernelParameters kernelparameters.KernelParameters
}

func New(id uuid.UUID, name string, description string, profile uuid.UUID, mac net.HardwareAddr, kernelParameters kernelparameters.KernelParameters) System {
	return System{
		Id:               id,
		Name:             name,
		Description:      description,
		Profile:          profile,
		Mac:              mac,
		KernelParameters: kernelParameters,
	}
}
