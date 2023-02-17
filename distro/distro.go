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

	if err := validateKernel(kernel); err != nil {
		return d, err
	}

	if err := validateInitrd(initrd); err != nil {
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

func validateKernel(kernel string) error {
	if kernel == "" {
		return errors.New("kernel cannot be empty")
	}

	p := "\\s"
	matched, err := regexp.MatchString(p, kernel)
	if err != nil {
		return err
	}
	if matched {
		return errors.New("kernel contains illegal characters")
	}

	return nil
}

func validateInitrd(initrd string) error {
	if initrd == "" {
		return errors.New("initrd cannot be empty")
	}

	p := "\\s"
	matched, err := regexp.MatchString(p, initrd)
	if err != nil {
		return err
	}
	if matched {
		return errors.New("initrd contains illegal characters")
	}

	return nil
}
