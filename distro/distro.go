package distro

import (
	"github.com/evanebb/gobble/kernelparameters"
)

type Distro struct {
	id               uint
	name             string
	description      string
	kernel           string
	initrd           string
	kernelParameters kernelparameters.KernelParameters
}

func New(id uint, name string, description string, kernel string, initrd string, kernelParameters kernelparameters.KernelParameters) Distro {
	return Distro{
		id:               id,
		name:             name,
		description:      description,
		kernel:           kernel,
		initrd:           initrd,
		kernelParameters: kernelParameters,
	}
}

func (d *Distro) Id() uint {
	return d.id
}

func (d *Distro) Name() string {
	return d.name
}

func (d *Distro) Description() string {
	return d.description
}

func (d *Distro) Kernel() string {
	return d.kernel
}

func (d *Distro) Initrd() string {
	return d.initrd
}

func (d *Distro) KernelParameters() kernelparameters.KernelParameters {
	return d.kernelParameters
}
