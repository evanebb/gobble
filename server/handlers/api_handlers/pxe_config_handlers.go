package api_handlers

import (
	"errors"
	"github.com/evanebb/gobble/api/response"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/profile"
	"github.com/evanebb/gobble/repository"
	"github.com/evanebb/gobble/system"
	"net"
	"net/http"
)

/*
 * HTTP handlers
 */

// PxeConfigHandlerGroup is a group of http.HandlerFunc functions related to PXE configs
type PxeConfigHandlerGroup struct {
	systemRepo  system.Repository
	profileRepo profile.Repository
}

func NewPxeConfigHandlerGroup(sr system.Repository, pr profile.Repository) PxeConfigHandlerGroup {
	return PxeConfigHandlerGroup{
		sr,
		pr,
	}
}

func (h PxeConfigHandlerGroup) GetPxeConfig(w http.ResponseWriter, r *http.Request) error {
	mac, err := net.ParseMAC(r.URL.Query().Get("mac"))
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	sys, err := h.systemRepo.GetSystemByMacAddress(mac)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return response.PlainText(w, http.StatusNotFound, system.RenderNotFound())
		}
		// This should be a 404, but iPXE won't load the script if that is the response code
		return NewHTTPError(err, http.StatusOK)
	}

	p, err := h.profileRepo.GetProfileById(sys.Profile)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	kp := kernelparameters.MergeKernelParameters(p.KernelParameters, sys.KernelParameters)
	pxeConfig := system.NewPxeConfig(p.Kernel, p.Initrd, kp)

	return response.PlainText(w, http.StatusOK, pxeConfig.Render())
}
