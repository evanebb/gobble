package server

import (
	"encoding/json"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/profile"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type profileRequest struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Distro           uuid.UUID `json:"distro"`
	KernelParameters []string  `json:"kernelParameters"`
}

type profileResponse struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Distro           uuid.UUID `json:"distro"`
	KernelParameters []string  `json:"kernelParameters"`
}

func (s *Server) getProfiles(w http.ResponseWriter, r *http.Request) error {
	profiles, err := s.profileRepo.GetProfiles()
	if err != nil {
		return err
	}

	resp := make([]profileResponse, 0)
	for _, p := range profiles {
		resp = append(resp, profileResponse{
			Id:               p.Id(),
			Name:             p.Name(),
			Description:      p.Description(),
			Distro:           p.Distro(),
			KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
		})
	}

	return SendJSONResponse(w, resp)
}

func (s *Server) getProfile(w http.ResponseWriter, r *http.Request) error {
	profileId, err := uuid.Parse(chi.URLParam(r, "profileID"))
	if err != nil {
		return err
	}

	p, err := s.profileRepo.GetProfileById(profileId)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, profileResponse{
		Id:               p.Id(),
		Name:             p.Name(),
		Description:      p.Description(),
		Distro:           p.Distro(),
		KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
	})
}

func (s *Server) createProfile(w http.ResponseWriter, r *http.Request) error {
	var req profileRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return err
	}

	profileId := uuid.New()

	kp, err := kernelparameters.New(req.KernelParameters)
	if err != nil {
		return err
	}

	p := profile.New(profileId, req.Name, req.Description, req.Distro, kp)
	err = s.profileRepo.SetProfile(p)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, profileResponse{
		Id:               p.Id(),
		Name:             p.Name(),
		Description:      p.Description(),
		Distro:           p.Distro(),
		KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
	})
}

func (s *Server) putProfile(w http.ResponseWriter, r *http.Request) error {
	var req profileRequest

	profileId, err := uuid.Parse(chi.URLParam(r, "profileID"))
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		return err
	}

	kp, err := kernelparameters.New(req.KernelParameters)
	if err != nil {
		return err
	}

	p := profile.New(profileId, req.Name, req.Description, req.Distro, kp)
	err = s.profileRepo.SetProfile(p)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, profileResponse{
		Id:               p.Id(),
		Name:             p.Name(),
		Description:      p.Description(),
		Distro:           p.Distro(),
		KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
	})
}

func (s *Server) patchProfile(w http.ResponseWriter, r *http.Request) error {
	profileId, err := uuid.Parse(chi.URLParam(r, "profileID"))
	if err != nil {
		return err
	}

	// Get and map the current profile to the API DTO
	p, err := s.profileRepo.GetProfileById(profileId)
	if err != nil {
		return err
	}

	req := profileRequest{
		Name:             p.Name(),
		Description:      p.Description(),
		Distro:           p.Distro(),
		KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
	}

	// Decode the request body into the current profile;
	// Values supplied in the body will overwrite the current values,
	// and anything that isn't supplied will be left alone
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		return err
	}

	kp, err := kernelparameters.New(req.KernelParameters)
	if err != nil {
		return err
	}

	// Map the DTO back to the model, this time with the newly supplied values from the request body
	p = profile.New(profileId, req.Name, req.Description, req.Distro, kp)
	err = s.profileRepo.SetProfile(p)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, profileResponse{
		Id:               p.Id(),
		Name:             p.Name(),
		Description:      p.Description(),
		Distro:           p.Distro(),
		KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
	})
}

func (s *Server) deleteProfile(w http.ResponseWriter, r *http.Request) error {
	profileId, err := uuid.Parse(chi.URLParam(r, "profileID"))
	if err != nil {
		return err
	}

	err = s.profileRepo.DeleteProfileById(profileId)
	if err != nil {
		return err
	}

	return SendJSONResponse(w, "successfully deleted profile")
}