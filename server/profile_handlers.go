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

func (s *Server) getProfiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profiles, err := s.profileRepo.GetProfiles()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

		encodedResponse, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedResponse)
	}
}

func (s *Server) getProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileId, err := uuid.Parse(chi.URLParam(r, "profileID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := s.profileRepo.GetProfileById(profileId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedProfile, err := json.Marshal(profileResponse{
			Id:               p.Id(),
			Name:             p.Name(),
			Description:      p.Description(),
			Distro:           p.Distro(),
			KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedProfile)
	}
}

func (s *Server) createProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req profileRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		profileId := uuid.New()

		kp, err := kernelparameters.New(req.KernelParameters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := profile.New(profileId, req.Name, req.Description, req.Distro, kp)
		err = s.profileRepo.SetProfile(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedProfile, err := json.Marshal(profileResponse{
			Id:               p.Id(),
			Name:             p.Name(),
			Description:      p.Description(),
			Distro:           p.Distro(),
			KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedProfile)
	}
}

func (s *Server) putProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req profileRequest

		profileId, err := uuid.Parse(chi.URLParam(r, "profileID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		kp, err := kernelparameters.New(req.KernelParameters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := profile.New(profileId, req.Name, req.Description, req.Distro, kp)
		err = s.profileRepo.SetProfile(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedProfile, err := json.Marshal(profileResponse{
			Id:               p.Id(),
			Name:             p.Name(),
			Description:      p.Description(),
			Distro:           p.Distro(),
			KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedProfile)
	}
}

func (s *Server) patchProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileId, err := uuid.Parse(chi.URLParam(r, "profileID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get and map the current profile to the API DTO
		p, err := s.profileRepo.GetProfileById(profileId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		kp, err := kernelparameters.New(req.KernelParameters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Map the DTO back to the model, this time with the newly supplied values from the request body
		p = profile.New(profileId, req.Name, req.Description, req.Distro, kp)
		err = s.profileRepo.SetProfile(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedProfile, err := json.Marshal(profileResponse{
			Id:               p.Id(),
			Name:             p.Name(),
			Description:      p.Description(),
			Distro:           p.Distro(),
			KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedProfile)
	}
}

func (s *Server) deleteProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileId, err := uuid.Parse(chi.URLParam(r, "profileID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.profileRepo.DeleteProfileById(profileId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("successfully deleted profile"))
	}
}
