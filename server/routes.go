package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) routes() {
	s.router.Use(middleware.Logger)

	s.router.Get("/", errorHandler(indexHandler))
	s.router.Route("/api", func(r chi.Router) {
		r.NotFound(errorHandler(unknownEndpointHandler))

		r.Route("/distros", func(r chi.Router) {
			r.Get("/", errorHandler(s.getDistros))
			r.Post("/", errorHandler(s.createDistro))

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", errorHandler(s.getDistro))
				r.Put("/", errorHandler(s.putDistro))
				r.Patch("/", errorHandler(s.patchDistro))
				r.Delete("/", errorHandler(s.deleteDistro))
			})
		})

		r.Route("/profiles", func(r chi.Router) {
			r.Get("/", errorHandler(s.getProfiles))
			r.Post("/", errorHandler(s.createProfile))

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", errorHandler(s.getProfile))
				r.Put("/", errorHandler(s.putProfile))
				r.Patch("/", errorHandler(s.patchProfile))
				r.Delete("/", errorHandler(s.deleteProfile))
			})
		})

		r.Route("/systems", func(r chi.Router) {
			r.Get("/", errorHandler(s.getSystems))
			r.Post("/", errorHandler(s.createSystem))

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", errorHandler(s.getSystem))
				r.Put("/", errorHandler(s.putSystem))
				r.Patch("/", errorHandler(s.patchSystem))
				r.Delete("/", errorHandler(s.deleteSystem))
			})
		})

		r.Get("/pxe-config", errorHandler(s.getPxeConfig))
	})
}
