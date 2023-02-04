package profile

import (
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
)

type Profile struct {
	id               uuid.UUID
	name             string
	description      string
	distro           uuid.UUID
	kernelParameters kernelparameters.KernelParameters
}

func New(id uuid.UUID, name string, description string, distro uuid.UUID, kernelParameters kernelparameters.KernelParameters) Profile {
	return Profile{
		id:               id,
		name:             name,
		description:      description,
		distro:           distro,
		kernelParameters: kernelParameters,
	}
}

func (p Profile) Id() uuid.UUID {
	return p.id
}

func (p Profile) Name() string {
	return p.name
}

func (p Profile) Description() string {
	return p.description
}

func (p Profile) Distro() uuid.UUID {
	return p.distro
}

func (p Profile) KernelParameters() kernelparameters.KernelParameters {
	return p.kernelParameters
}
