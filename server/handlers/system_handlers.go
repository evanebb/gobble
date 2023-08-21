package handlers

import (
	"encoding/json"
	"errors"
	"github.com/evanebb/gobble/api/response"
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
 * HTTP handlers
 */

// SystemHandlerGroup is a group of http.HandlerFunc functions related to systems
type SystemHandlerGroup struct {
	systemRepo system.Repository
}

func NewSystemHandlerGroup(sr system.Repository) SystemHandlerGroup {
	return SystemHandlerGroup{sr}
}

func (h SystemHandlerGroup) GetSystems(w http.ResponseWriter, r *http.Request) error {
	systems, err := h.systemRepo.GetSystems()
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	resp := make([]systemResponse, 0)
	for _, sys := range systems {
		resp = append(resp, newSystemResponse(sys))
	}

	return response.Success(w, http.StatusOK, resp)
}

func (h SystemHandlerGroup) GetSystem(w http.ResponseWriter, r *http.Request) error {
	systemId, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	sys, err := h.systemRepo.GetSystemById(systemId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NewHTTPError(err, http.StatusNotFound)
		}
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusOK, newSystemResponse(sys))
}

func (h SystemHandlerGroup) CreateSystem(w http.ResponseWriter, r *http.Request) error {
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

	sys, err := system.New(systemId, req.Name, req.Description, req.Profile, macAddress, kp)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = h.systemRepo.SetSystem(sys)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusCreated, newSystemResponse(sys))
}

func (h SystemHandlerGroup) PutSystem(w http.ResponseWriter, r *http.Request) error {
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

	sys, err := system.New(systemId, req.Name, req.Description, req.Profile, macAddress, kp)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = h.systemRepo.SetSystem(sys)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusOK, newSystemResponse(sys))
}

func (h SystemHandlerGroup) PatchSystem(w http.ResponseWriter, r *http.Request) error {
	systemId, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	// Get and map the current system to the API DTO
	sys, err := h.systemRepo.GetSystemById(systemId)
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
	sys, err = system.New(systemId, req.Name, req.Description, req.Profile, macAddress, kp)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = h.systemRepo.SetSystem(sys)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusOK, newSystemResponse(sys))
}

func (h SystemHandlerGroup) DeleteSystem(w http.ResponseWriter, r *http.Request) error {
	systemId, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = h.systemRepo.DeleteSystemById(systemId)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusNoContent, nil)
}
