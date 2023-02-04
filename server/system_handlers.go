package server

import (
	"encoding/json"
	"fmt"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/system"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net"
	"net/http"
)

type systemRequest struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Profile          uuid.UUID `json:"profile"`
	Mac              string    `json:"mac"`
	KernelParameters []string  `json:"kernelParameters"`
}

type systemResponse struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Profile          uuid.UUID `json:"profile"`
	Mac              string    `json:"mac"`
	KernelParameters []string  `json:"kernelParameters"`
}

func (s *Server) getSystems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systems, err := s.systemRepo.GetSystems()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := make([]systemResponse, 0)
		for _, sys := range systems {
			resp = append(resp, systemResponse{
				Id:               sys.Id(),
				Name:             sys.Name(),
				Description:      sys.Description(),
				Profile:          sys.Profile(),
				Mac:              sys.Mac().String(),
				KernelParameters: kernelparameters.FormatKernelParameters(sys.KernelParameters()),
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

func (s *Server) getSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemId, err := uuid.Parse(chi.URLParam(r, "systemID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sys, err := s.systemRepo.GetSystemById(systemId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedSystem, err := json.Marshal(systemResponse{
			Id:               sys.Id(),
			Name:             sys.Name(),
			Description:      sys.Description(),
			Profile:          sys.Profile(),
			Mac:              sys.Mac().String(),
			KernelParameters: kernelparameters.FormatKernelParameters(sys.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedSystem)
	}
}

func (s *Server) createSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req systemRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		systemId := uuid.New()

		kp, err := kernelparameters.New(req.KernelParameters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		macAddress, err := net.ParseMAC(req.Mac)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sys := system.New(systemId, req.Name, req.Description, req.Profile, macAddress, kp)
		err = s.systemRepo.SetSystem(sys)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedSystem, err := json.Marshal(systemResponse{
			Id:               sys.Id(),
			Name:             sys.Name(),
			Description:      sys.Description(),
			Profile:          sys.Profile(),
			Mac:              sys.Mac().String(),
			KernelParameters: kernelparameters.FormatKernelParameters(sys.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedSystem)
	}
}

func (s *Server) putSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req systemRequest

		systemId, err := uuid.Parse(chi.URLParam(r, "systemID"))
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

		macAddress, err := net.ParseMAC(req.Mac)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sys := system.New(systemId, req.Name, req.Description, req.Profile, macAddress, kp)
		err = s.systemRepo.SetSystem(sys)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedSystem, err := json.Marshal(systemResponse{
			Id:               sys.Id(),
			Name:             sys.Name(),
			Description:      sys.Description(),
			Profile:          sys.Profile(),
			Mac:              sys.Mac().String(),
			KernelParameters: kernelparameters.FormatKernelParameters(sys.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedSystem)
	}
}

func (s *Server) patchSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemId, err := uuid.Parse(chi.URLParam(r, "systemID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get and map the current system to the API DTO
		sys, err := s.systemRepo.GetSystemById(systemId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := systemRequest{
			Name:             sys.Name(),
			Description:      sys.Description(),
			Profile:          sys.Profile(),
			Mac:              sys.Mac().String(),
			KernelParameters: kernelparameters.FormatKernelParameters(sys.KernelParameters()),
		}

		// Decode the request body into the current system;
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

		macAddress, err := net.ParseMAC(req.Mac)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Map the DTO back to the model, this time with the newly supplied values from the request body
		sys = system.New(systemId, req.Name, req.Description, req.Profile, macAddress, kp)
		err = s.systemRepo.SetSystem(sys)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedSystem, err := json.Marshal(systemResponse{
			Id:               sys.Id(),
			Name:             sys.Name(),
			Description:      sys.Description(),
			Profile:          sys.Profile(),
			Mac:              sys.Mac().String(),
			KernelParameters: kernelparameters.FormatKernelParameters(sys.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedSystem)
	}
}

func (s *Server) deleteSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemId, err := uuid.Parse(chi.URLParam(r, "systemID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.systemRepo.DeleteSystemById(systemId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("successfully deleted system"))
	}
}

func (s *Server) getPxeConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mac, err := net.ParseMAC(r.URL.Query().Get("mac"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		systemService := system.NewService(s.distroRepo, s.profileRepo, s.systemRepo)
		pxeConfig, err := systemService.RenderPxeConfig(mac)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, _ = fmt.Fprint(w, pxeConfig)
	}
}
