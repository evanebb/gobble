package server

import (
	"encoding/json"
	"github.com/evanebb/gobble/distro"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"net/http"
)

/*
 * Request and response structures, and their supporting functions
 */

// distroRequest is the JSON representation of a distro.Distro that is accepted by the API.
type distroRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Kernel           string   `json:"kernel"`
	Initrd           string   `json:"initrd"`
	KernelParameters []string `json:"kernelParameters"`
}

// distroResponse is the JSON representation of a distro.Distro that is returned by the API.
type distroResponse struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Kernel           string    `json:"kernel"`
	Initrd           string    `json:"initrd"`
	KernelParameters []string  `json:"kernelParameters"`
}

// newDistroResponse accepts a distro.Distro, and casts it into a distroResponse.
func newDistroResponse(d distro.Distro) distroResponse {
	return distroResponse{
		Id:               d.Id,
		Name:             d.Name,
		Description:      d.Description,
		Kernel:           d.Kernel,
		Initrd:           d.Initrd,
		KernelParameters: kernelparameters.FormatKernelParameters(d.KernelParameters),
	}
}

/*
 * The actual HTTP handlers
 */

func (s *Server) getDistros(w http.ResponseWriter, r *http.Request) error {
	distros, err := s.distroRepo.GetDistros()
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	resp := make([]distroResponse, 0)
	for _, d := range distros {
		resp = append(resp, newDistroResponse(d))
	}

	return SendSuccessResponse(w, http.StatusOK, resp)
}

func (s *Server) getDistro(w http.ResponseWriter, r *http.Request) error {
	distroId, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	d, err := s.distroRepo.GetDistroById(distroId)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusOK, newDistroResponse(d))
}

func (s *Server) createDistro(w http.ResponseWriter, r *http.Request) error {
	var req distroRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	distroId := uuid.New()

	kp, err := kernelparameters.New(req.KernelParameters)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	d, err := distro.New(distroId, req.Name, req.Description, req.Kernel, req.Initrd, kp)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = s.distroRepo.SetDistro(d)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusCreated, newDistroResponse(d))
}

func (s *Server) putDistro(w http.ResponseWriter, r *http.Request) error {
	var req distroRequest

	distroId, err := getUUIDFromRequest(r)
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

	d, err := distro.New(distroId, req.Name, req.Description, req.Kernel, req.Initrd, kp)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = s.distroRepo.SetDistro(d)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusOK, newDistroResponse(d))
}

func (s *Server) patchDistro(w http.ResponseWriter, r *http.Request) error {
	distroID, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	// Get and map the current distro to the API DTO
	d, err := s.distroRepo.GetDistroById(distroID)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	req := distroRequest{
		Name:             d.Name,
		Description:      d.Description,
		Kernel:           d.Kernel,
		Initrd:           d.Initrd,
		KernelParameters: kernelparameters.FormatKernelParameters(d.KernelParameters),
	}

	// Decode the request body into the current distro;
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

	// Map the DTO back to the model, this time with the newly supplied values from the request body
	d, err = distro.New(distroID, req.Name, req.Description, req.Kernel, req.Initrd, kp)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = s.distroRepo.SetDistro(d)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusOK, newDistroResponse(d))
}

func (s *Server) deleteDistro(w http.ResponseWriter, r *http.Request) error {
	distroID, err := getUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = s.distroRepo.DeleteDistroById(distroID)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return SendSuccessResponse(w, http.StatusNoContent, nil)
}
