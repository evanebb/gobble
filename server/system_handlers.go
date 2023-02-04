package server

import (
	"encoding/json"
	"fmt"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/system"
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

		systemJson, err := json.Marshal(systemResponse{
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
		w.Write(systemJson)
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
