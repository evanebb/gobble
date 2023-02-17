package system

import (
	"errors"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"net"
	"regexp"
)

type System struct {
	Id               uuid.UUID
	Name             string
	Description      string
	Profile          uuid.UUID
	Mac              net.HardwareAddr
	KernelParameters kernelparameters.KernelParameters
}

func New(id uuid.UUID, name string, description string, profile uuid.UUID, mac net.HardwareAddr, kernelParameters kernelparameters.KernelParameters) (System, error) {
	var s System

	if err := validateName(name); err != nil {
		return s, err
	}

	return System{
		Id:               id,
		Name:             name,
		Description:      description,
		Profile:          profile,
		Mac:              mac,
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
