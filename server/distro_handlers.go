package server

import (
	"encoding/json"
	"github.com/evanebb/gobble/distro"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type distroRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Kernel           string   `json:"kernel"`
	Initrd           string   `json:"initrd"`
	KernelParameters []string `json:"kernelParameters"`
}

type distroResponse struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Kernel           string    `json:"kernel"`
	Initrd           string    `json:"initrd"`
	KernelParameters []string  `json:"kernelParameters"`
}

func (s *Server) getDistros(w http.ResponseWriter, r *http.Request) error {
	distros, err := s.distroRepo.GetDistros()
	if err != nil {
		return err
	}

	resp := make([]distroResponse, 0)
	for _, d := range distros {
		resp = append(resp, distroResponse{
			Id:               d.Id(),
			Name:             d.Name(),
			Description:      d.Description(),
			Kernel:           d.Kernel(),
			Initrd:           d.Initrd(),
			KernelParameters: kernelparameters.FormatKernelParameters(d.KernelParameters()),
		})
	}

	return SendJSONResponse(w, resp)
}

func (s *Server) getDistro(w http.ResponseWriter, r *http.Request) error {
	distroId, err := uuid.Parse(chi.URLParam(r, "distroID"))
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	d, err := s.distroRepo.GetDistroById(distroId)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, distroResponse{
		Id:               d.Id(),
		Name:             d.Name(),
		Description:      d.Description(),
		Kernel:           d.Kernel(),
		Initrd:           d.Initrd(),
		KernelParameters: kernelparameters.FormatKernelParameters(d.KernelParameters()),
	})
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
		return err
	}

	err = s.distroRepo.SetDistro(d)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, distroResponse{
		Id:               d.Id(),
		Name:             d.Name(),
		Description:      d.Description(),
		Kernel:           d.Kernel(),
		Initrd:           d.Initrd(),
		KernelParameters: kernelparameters.FormatKernelParameters(d.KernelParameters()),
	})
}

func (s *Server) putDistro(w http.ResponseWriter, r *http.Request) error {
	var req distroRequest

	distroId, err := uuid.Parse(chi.URLParam(r, "distroID"))
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
		return err
	}

	err = s.distroRepo.SetDistro(d)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, distroResponse{
		Id:               d.Id(),
		Name:             d.Name(),
		Description:      d.Description(),
		Kernel:           d.Kernel(),
		Initrd:           d.Initrd(),
		KernelParameters: kernelparameters.FormatKernelParameters(d.KernelParameters()),
	})
}

func (s *Server) patchDistro(w http.ResponseWriter, r *http.Request) error {
	distroID, err := uuid.Parse(chi.URLParam(r, "distroID"))
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	// Get and map the current distro to the API DTO
	d, err := s.distroRepo.GetDistroById(distroID)
	if err != nil {
		return err
	}

	req := distroRequest{
		Name:             d.Name(),
		Description:      d.Description(),
		Kernel:           d.Kernel(),
		Initrd:           d.Initrd(),
		KernelParameters: kernelparameters.FormatKernelParameters(d.KernelParameters()),
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
		return err
	}

	err = s.distroRepo.SetDistro(d)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, distroResponse{
		Id:               d.Id(),
		Name:             d.Name(),
		Description:      d.Description(),
		Kernel:           d.Kernel(),
		Initrd:           d.Initrd(),
		KernelParameters: kernelparameters.FormatKernelParameters(d.KernelParameters()),
	})
}

func (s *Server) deleteDistro(w http.ResponseWriter, r *http.Request) error {
	distroID, err := uuid.Parse(chi.URLParam(r, "distroID"))
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = s.distroRepo.DeleteDistroById(distroID)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, "successfully deleted distro")
}
