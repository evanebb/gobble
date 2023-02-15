package distro

import (
	"errors"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"regexp"
)

type Distro struct {
	id               uuid.UUID
	name             string
	description      string
	kernel           string
	initrd           string
	kernelParameters kernelparameters.KernelParameters
}

func New(id uuid.UUID, name string, description string, kernel string, initrd string, kernelParameters kernelparameters.KernelParameters) (Distro, error) {
	var d Distro

	if err := validateName(name); err != nil {
		return d, err
	}

	return Distro{
		id:               id,
		name:             name,
		description:      description,
		kernel:           kernel,
		initrd:           initrd,
		kernelParameters: kernelParameters,
	}, nil
}

func validateName(name string) error {
	p := "^[a-zA-Z0-9-_.()]{1,64}$"
	matched, err := regexp.MatchString(p, name)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("name contains illegal characters")
	}

	return nil
}

func (d Distro) Id() uuid.UUID {
	return d.id
}

func (d Distro) Name() string {
	return d.name
}

func (d Distro) Description() string {
	return d.description
}

func (d Distro) Kernel() string {
	return d.kernel
}

func (d Distro) Initrd() string {
	return d.initrd
}

func (d Distro) KernelParameters() kernelparameters.KernelParameters {
	return d.kernelParameters
}
