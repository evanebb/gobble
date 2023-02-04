package profile

import (
	"github.com/evanebb/gobble/kernelparameters"
)

type Profile struct {
	id               uint
	name             string
	description      string
	distro           uint
	kernelParameters kernelparameters.KernelParameters
}

func New(id uint, name string, description string, distro uint, kernelParameters kernelparameters.KernelParameters) Profile {
	return Profile{
		id:               id,
		name:             name,
		description:      description,
		distro:           distro,
		kernelParameters: kernelParameters,
	}
}

func (p Profile) Id() uint {
	return p.id
}

func (p Profile) Name() string {
	return p.name
}

func (p Profile) Description() string {
	return p.description
}

func (p Profile) Distro() uint {
	return p.distro
}

func (p Profile) KernelParameters() kernelparameters.KernelParameters {
	return p.kernelParameters
}
