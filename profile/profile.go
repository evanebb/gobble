package profile

import (
	"errors"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"regexp"
)

type Profile struct {
	Id               uuid.UUID
	Name             string
	Description      string
	Distro           uuid.UUID
	KernelParameters kernelparameters.KernelParameters
}

func New(id uuid.UUID, name string, description string, distro uuid.UUID, kernelParameters kernelparameters.KernelParameters) (Profile, error) {
	var p Profile

	if err := validateName(name); err != nil {
		return p, err
	}

	return Profile{
		Id:               id,
		Name:             name,
		Description:      description,
		Distro:           distro,
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
