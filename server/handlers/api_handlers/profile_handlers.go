package api_handlers

import (
	"encoding/json"
	"errors"
	"github.com/evanebb/gobble/api/response"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/profile"
	"github.com/evanebb/gobble/repository"
	"github.com/evanebb/gobble/server/handlers"
	"github.com/google/uuid"
	"net/http"
)

/*
 * Request and response structures, and their supporting functions
 */

// profileRequest is the JSON representation of a profile.Profile that is accepted by the API.
type profileRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Kernel           string   `json:"kernel"`
	Initrd           string   `json:"initrd"`
	KernelParameters []string `json:"kernelParameters"`
}

// profileResponse is the JSON representation of a profile.Profile that is returned by the API.
type profileResponse struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Kernel           string    `json:"kernel"`
	Initrd           string    `json:"initrd"`
	KernelParameters []string  `json:"kernelParameters"`
}

// newProfileResponse accepts a profile.Profile, and casts it to a profileResponse.
func newProfileResponse(p profile.Profile) profileResponse {
	return profileResponse{
		Id:               p.Id,
		Name:             p.Name,
		Description:      p.Description,
		Kernel:           p.Kernel,
		Initrd:           p.Initrd,
		KernelParameters: p.KernelParameters.StringSlice(),
	}
}

/*
 * HTTP handlers
 */

// ProfileHandlerGroup is a group of http.HandlerFunc functions related to profiles
type ProfileHandlerGroup struct {
	profileRepo profile.Repository
}

func NewProfileHandlerGroup(pr profile.Repository) ProfileHandlerGroup {
	return ProfileHandlerGroup{pr}
}

func (h ProfileHandlerGroup) GetProfiles(w http.ResponseWriter, r *http.Request) error {
	profiles, err := h.profileRepo.GetProfiles()
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	resp := make([]profileResponse, 0)
	for _, p := range profiles {
		resp = append(resp, newProfileResponse(p))
	}

	return response.Success(w, http.StatusOK, resp)
}

func (h ProfileHandlerGroup) GetProfile(w http.ResponseWriter, r *http.Request) error {
	profileId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	p, err := h.profileRepo.GetProfileById(profileId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NewHTTPError(err, http.StatusNotFound)
		}
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusOK, newProfileResponse(p))
}

func (h ProfileHandlerGroup) CreateProfile(w http.ResponseWriter, r *http.Request) error {
	var req profileRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	profileId := uuid.New()

	kp, err := kernelparameters.ParseStringSlice(req.KernelParameters)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	p, err := profile.New(profileId, req.Name, req.Description, req.Kernel, req.Initrd, kp)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = h.profileRepo.SetProfile(p)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusCreated, newProfileResponse(p))
}

func (h ProfileHandlerGroup) PutProfile(w http.ResponseWriter, r *http.Request) error {
	var req profileRequest

	profileId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	kp, err := kernelparameters.ParseStringSlice(req.KernelParameters)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	p, err := profile.New(profileId, req.Name, req.Description, req.Kernel, req.Initrd, kp)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = h.profileRepo.SetProfile(p)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusOK, newProfileResponse(p))
}

func (h ProfileHandlerGroup) PatchProfile(w http.ResponseWriter, r *http.Request) error {
	profileId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	// Get and map the current profile to the API DTO
	p, err := h.profileRepo.GetProfileById(profileId)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	req := profileRequest{
		Name:             p.Name,
		Description:      p.Description,
		Kernel:           p.Kernel,
		Initrd:           p.Initrd,
		KernelParameters: p.KernelParameters.StringSlice(),
	}

	// Decode the request body into the current profile;
	// Values supplied in the body will overwrite the current values,
	// and anything that isn't supplied will be left alone
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	kp, err := kernelparameters.ParseStringSlice(req.KernelParameters)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	// Map the DTO back to the model, this time with the newly supplied values from the request body
	p, err = profile.New(profileId, req.Name, req.Description, req.Kernel, req.Initrd, kp)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = h.profileRepo.SetProfile(p)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusOK, newProfileResponse(p))
}

func (h ProfileHandlerGroup) DeleteProfile(w http.ResponseWriter, r *http.Request) error {
	profileId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = h.profileRepo.DeleteProfileById(profileId)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	// No data to return, just pass nil
	return response.Success(w, http.StatusNoContent, nil)
}
