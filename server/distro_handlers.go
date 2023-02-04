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

func (s *Server) getDistros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		distros, err := s.distroRepo.GetDistros()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

		encodedResponse, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedResponse)
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
		var req distroRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		distroId := uuid.New()

		kp, err := kernelparameters.New(req.KernelParameters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		d := distro.New(distroId, req.Name, req.Description, req.Kernel, req.Initrd, kp)
		err = s.distroRepo.SetDistro(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedDistro, err := json.Marshal(distroResponse{
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
		w.Write(encodedDistro)
	}
}

func (s *Server) putDistro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req distroRequest

		distroId, err := uuid.Parse(chi.URLParam(r, "distroID"))
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

		d := distro.New(distroId, req.Name, req.Description, req.Kernel, req.Initrd, kp)
		err = s.distroRepo.SetDistro(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedDistro, err := json.Marshal(distroResponse{
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
		w.Write(encodedDistro)
	}
}

func (s *Server) patchDistro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		distroID, err := uuid.Parse(chi.URLParam(r, "distroID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get and map the current distro to the API DTO
		d, err := s.distroRepo.GetDistroById(distroID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		kp, err := kernelparameters.New(req.KernelParameters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Map the DTO back to the model, this time with the newly supplied values from the request body
		d = distro.New(distroID, req.Name, req.Description, req.Kernel, req.Initrd, kp)
		err = s.distroRepo.SetDistro(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedDistro, err := json.Marshal(distroResponse{
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
		w.Write(encodedDistro)
	}
}

func (s *Server) deleteDistro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		distroID, err := uuid.Parse(chi.URLParam(r, "distroID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.distroRepo.DeleteDistroById(distroID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("successfully deleted distro"))
	}
}
