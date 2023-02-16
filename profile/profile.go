package profile

import (
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
)

type Profile struct {
	Id               uuid.UUID
	Name             string
	Description      string
	Distro           uuid.UUID
	KernelParameters kernelparameters.KernelParameters
}

func New(id uuid.UUID, name string, description string, distro uuid.UUID, kernelParameters kernelparameters.KernelParameters) Profile {
	return Profile{
		Id:               id,
		Name:             name,
		Description:      description,
		Distro:           distro,
		KernelParameters: kernelParameters,
	}
}
