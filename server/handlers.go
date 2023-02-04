package server

import (
	"encoding/json"
	"fmt"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/profile"
	"github.com/evanebb/gobble/system"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
	"strconv"
)

func (s *Server) getDistros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) getDistro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) createDistro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) putDistro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) patchDistro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) deleteDistro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

type profileRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Distro           uint     `json:"distro"`
	KernelParameters []string `json:"kernelParameters"`
}

type profileResponse struct {
	Id               uint     `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Distro           uint     `json:"distro"`
	KernelParameters []string `json:"kernelParameters"`
}

func (s *Server) getProfiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profiles, err := s.profileRepo.GetProfiles()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var profilesResp []profileResponse
		for _, p := range profiles {
			profilesResp = append(profilesResp, profileResponse{
				Id:               p.Id(),
				Name:             p.Name(),
				Description:      p.Description(),
				Distro:           p.Distro(),
				KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
			})
		}

		jsonResp, err := json.Marshal(profilesResp)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func (s *Server) getProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileId, err := strconv.Atoi(chi.URLParam(r, "profileID"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		p, err := s.profileRepo.GetProfileById(uint(profileId))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		profileJson, err := json.Marshal(profileResponse{
			Id:               p.Id(),
			Name:             p.Name(),
			Description:      p.Description(),
			Distro:           p.Distro(),
			KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(profileJson)
	}
}

func (s *Server) createProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req profileRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// FIXME: how the hell do I instantiate a new profile without an ID; fix: just don't leak the database id into your model :)
	}
}

func (s *Server) putProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req profileRequest

		id, err := strconv.Atoi(chi.URLParam(r, "profileID"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		profileId := uint(id)

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		kp, err := kernelparameters.New(req.KernelParameters)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		p := profile.New(profileId, req.Name, req.Description, req.Distro, kp)
		err = s.profileRepo.SetProfile(p)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		profileJson, err := json.Marshal(profileResponse{
			Id:               p.Id(),
			Name:             p.Name(),
			Description:      p.Description(),
			Distro:           p.Distro(),
			KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(profileJson)
	}
}

func (s *Server) patchProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "profileID"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		profileId := uint(id)

		// Get and map the current profile to the API DTO
		p, err := s.profileRepo.GetProfileById(profileId)
		if err != nil {
			http.Error(w, err.Error(), 400)
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
			http.Error(w, err.Error(), 400)
			return
		}

		kp, err := kernelparameters.New(req.KernelParameters)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		// Map the DTO back to the model, this time with the newly supplied values from the request body
		p = profile.New(profileId, req.Name, req.Description, req.Distro, kp)
		err = s.profileRepo.SetProfile(p)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		profileJson, err := json.Marshal(profileResponse{
			Id:               p.Id(),
			Name:             p.Name(),
			Description:      p.Description(),
			Distro:           p.Distro(),
			KernelParameters: kernelparameters.FormatKernelParameters(p.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(profileJson)
	}
}

func (s *Server) deleteProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) getSystems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) getSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) createSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) putSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) patchSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) deleteSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) getPxeConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mac, err := net.ParseMAC(r.URL.Query().Get("mac"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		systemService := system.NewService(s.distroRepo, s.profileRepo, s.systemRepo)
		pxeConfig, err := systemService.RenderPxeConfig(mac)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		_, _ = fmt.Fprint(w, pxeConfig)
	}
}