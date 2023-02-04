package server

import (
	"encoding/json"
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

func (s *Server) getDistros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", 501)
	}
}

func (s *Server) getDistro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		distroId, err := uuid.Parse(chi.URLParam(r, "distroID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		d, err := s.distroRepo.GetDistroById(distroId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		distroJson, err := json.Marshal(distroResponse{
			Id:               d.Id(),
			Name:             d.Name(),
			Description:      d.Description(),
			Kernel:           d.Kernel(),
			Initrd:           d.Initrd(),
			KernelParameters: kernelparameters.FormatKernelParameters(d.KernelParameters()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(distroJson)
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
