package server

import (
	"encoding/json"
	"errors"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/repository"
	"github.com/evanebb/gobble/system"
	"github.com/google/uuid"
	"net"
	"net/http"
)

/*
 * Request and response structures, and their supporting functions
 */

// systemRequest is the JSON representation of a system.System that is accepted by the API.
type systemRequest struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Profile          uuid.UUID `json:"profile"`
	Mac              string    `json:"mac"`
	KernelParameters []string  `json:"kernelParameters"`
}

// systemResponse is the JSON representation of a system.System that is returned by the API.
type systemResponse struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Profile          uuid.UUID `json:"profile"`
	Mac              string    `json:"mac"`
	KernelParameters []string  `json:"kernelParameters"`
}

// newSystemResponse accepts a system.System, and casts it into a systemResponse.
func newSystemResponse(sys system.System) systemResponse {
	return systemResponse{
		Id:               sys.Id,
		Name:             sys.Name,
		Description:      sys.Description,
		Profile:          sys.Profile,
		Mac:              sys.Mac.String(),
		KernelParameters: kernelparameters.FormatKernelParameters(sys.KernelParameters),
	}
}

/*
 * The actual HTTP handlers
 */

func (s *Server) getSystems(w http.ResponseWriter, r *http.Request) error {
	systems, err := s.systemRepo.GetSystems()
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	resp := make([]systemResponse, 0)
	for _, sys := range systems {
		resp = append(resp, newSystemResponse(sys))
	}

	return SendSuccessResponse(w, http.StatusOK, resp)
}

func (s *Server) getSystem(w http.ResponseWriter, r *http.Request) error {
	systemId, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	sys, err := s.systemRepo.GetSystemById(systemId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NewHTTPError(err, http.StatusNotFound)
		}
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusOK, newSystemResponse(sys))
}

func (s *Server) createSystem(w http.ResponseWriter, r *http.Request) error {
	var req systemRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	systemId := uuid.New()

	kp, err := kernelparameters.New(req.KernelParameters)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	macAddress, err := net.ParseMAC(req.Mac)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	sys := system.New(systemId, req.Name, req.Description, req.Profile, macAddress, kp)
	err = s.systemRepo.SetSystem(sys)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusCreated, newSystemResponse(sys))
}

func (s *Server) putSystem(w http.ResponseWriter, r *http.Request) error {
	var req systemRequest

	systemId, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	kp, err := kernelparameters.New(req.KernelParameters)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	macAddress, err := net.ParseMAC(req.Mac)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	sys := system.New(systemId, req.Name, req.Description, req.Profile, macAddress, kp)
	err = s.systemRepo.SetSystem(sys)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusOK, newSystemResponse(sys))
}

func (s *Server) patchSystem(w http.ResponseWriter, r *http.Request) error {
	systemId, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	// Get and map the current system to the API DTO
	sys, err := s.systemRepo.GetSystemById(systemId)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	req := systemRequest{
		Name:             sys.Name,
		Description:      sys.Description,
		Profile:          sys.Profile,
		Mac:              sys.Mac.String(),
		KernelParameters: kernelparameters.FormatKernelParameters(sys.KernelParameters),
	}

	// Decode the request body into the current system;
	// Values supplied in the body will overwrite the current values,
	// and anything that isn't supplied will be left alone
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	kp, err := kernelparameters.New(req.KernelParameters)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	macAddress, err := net.ParseMAC(req.Mac)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	// Map the DTO back to the model, this time with the newly supplied values from the request body
	sys = system.New(systemId, req.Name, req.Description, req.Profile, macAddress, kp)
	err = s.systemRepo.SetSystem(sys)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusOK, newSystemResponse(sys))
}

func (s *Server) deleteSystem(w http.ResponseWriter, r *http.Request) error {
	systemId, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = s.systemRepo.DeleteSystemById(systemId)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusNoContent, nil)
}

func (s *Server) getPxeConfig(w http.ResponseWriter, r *http.Request) error {
	mac, err := net.ParseMAC(r.URL.Query().Get("mac"))
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	sys, err := s.systemRepo.GetSystemByMacAddress(mac)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NewHTTPError(err, http.StatusNotFound)
		}
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	p, err := s.profileRepo.GetProfileById(sys.Profile)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	d, err := s.distroRepo.GetDistroById(p.Distro)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	kp := kernelparameters.MergeKernelParameters(d.KernelParameters, p.KernelParameters, sys.KernelParameters)
	pxeConfig := system.NewPxeConfig(d.Kernel, d.Initrd, kp)

	return SendPlainTextResponse(w, http.StatusOK, pxeConfig.Render())
}
