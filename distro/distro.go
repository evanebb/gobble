package distro

import (
	"errors"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"regexp"
)

type Distro struct {
	Id               uuid.UUID
	Name             string
	Description      string
	Kernel           string
	Initrd           string
	KernelParameters kernelparameters.KernelParameters
}

func New(id uuid.UUID, name string, description string, kernel string, initrd string, kernelParameters kernelparameters.KernelParameters) (Distro, error) {
	var d Distro

	if err := validateName(name); err != nil {
		return d, err
	}

	return Distro{
		Id:               id,
		Name:             name,
		Description:      description,
		Kernel:           kernel,
		Initrd:           initrd,
		KernelParameters: kernelParameters,
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
