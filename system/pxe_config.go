package system

import (
	"fmt"
	"github.com/evanebb/gobble/kernelparameters"
	"strings"
)

// TODO: create a PxeConfig interface with a Render() method, have the System struct implement it

var template = `#!ipxe

kernel %s %s
initrd %s

boot
`

type PxeConfig struct {
	Kernel           string
	Initrd           string
	KernelParameters kernelparameters.KernelParameters
}

func NewPxeConfig(kernel string, initrd string, kernelParameters kernelparameters.KernelParameters) PxeConfig {
	return PxeConfig{
		Kernel:           kernel,
		Initrd:           initrd,
		KernelParameters: kernelParameters,
	}
}

func (p PxeConfig) Render() string {
	kp := kernelparameters.FormatKernelParameters(p.KernelParameters)
	kpStr := strings.Join(kp, " ")
	return fmt.Sprintf(template, p.Kernel, kpStr, p.Initrd)
}
