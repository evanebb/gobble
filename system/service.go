package system

import (
	"fmt"
	"github.com/evanebb/gobble/distro"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/profile"
	"net"
	"strings"
)

var pxeConfigTemplate = `#!ipxe

kernel %s %s
initrd %s

boot
`

type Service struct {
	distroRepo  distro.Repository
	profileRepo profile.Repository
	systemRepo  Repository
}

func NewService(dr distro.Repository, pr profile.Repository, sr Repository) Service {
	return Service{
		distroRepo:  dr,
		profileRepo: pr,
		systemRepo:  sr,
	}
}

func (s Service) RenderPxeConfig(mac net.HardwareAddr) (string, error) {
	sys, err := s.systemRepo.GetSystemByMacAddress(mac)
	if err != nil {
		return "", err
	}

	p, err := s.profileRepo.GetProfileById(sys.Profile())
	if err != nil {
		return "", err
	}

	d, err := s.distroRepo.GetDistroById(p.Distro())
	if err != nil {
		return "", err
	}

	kp := kernelparameters.MergeKernelParameters(d.KernelParameters(), p.KernelParameters(), sys.KernelParameters())
	kpSlice := kernelparameters.FormatKernelParameters(kp)
	kpString := strings.Join(kpSlice, " ")

	return fmt.Sprintf(pxeConfigTemplate, d.Kernel(), kpString, d.Initrd()), nil
}
