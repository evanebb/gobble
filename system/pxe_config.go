package system

import (
	"fmt"
	"github.com/evanebb/gobble/kernelparameters"
	"strings"
)

// TODO: create a PxeConfig interface with a Render() method, have the System struct implement it

// iPXE script template that is served to clients.
// The configured kernel, initrd and kernel parameters are substituted into it.
var template = `#!ipxe

kernel %s %s
initrd %s

boot
`

// iPXE script template that is served if the system is not registered and no profile can be found for it.
var notFoundTemplate = `#!ipxe

echo No matching profile found for system!
prompt --timeout 30000 Press any key or wait 30 seconds to continue ||

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

func RenderNotFound() string {
	return notFoundTemplate
}
