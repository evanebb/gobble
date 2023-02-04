package system

import (
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"net"
)

type System struct {
	id               uuid.UUID
	name             string
	description      string
	profile          uuid.UUID
	mac              net.HardwareAddr
	kernelParameters kernelparameters.KernelParameters
}

func New(id uuid.UUID, name string, description string, profile uuid.UUID, mac net.HardwareAddr, kernelParameters kernelparameters.KernelParameters) System {
	return System{
		id:               id,
		name:             name,
		description:      description,
		profile:          profile,
		mac:              mac,
		kernelParameters: kernelParameters,
	}
}

func (s *System) Id() uuid.UUID {
	return s.id
}

func (s *System) Name() string {
	return s.name
}

func (s *System) Description() string {
	return s.description
}

func (s *System) Profile() uuid.UUID {
	return s.profile
}

func (s *System) Mac() net.HardwareAddr {
	return s.mac
}

func (s *System) KernelParameters() kernelparameters.KernelParameters {
	return s.kernelParameters
}
