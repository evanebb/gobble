package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) routes() {
	s.router.Use(middleware.Logger)
	s.router.Route("/api", func(r chi.Router) {
		// FIXME: set JSON header right here :)
		r.Route("/distros", func(r chi.Router) {
			r.Get("/", s.getDistros())
			r.Post("/", s.createDistro())

			r.Route("/{distroID}", func(r chi.Router) {
				r.Get("/", s.getDistro())
				r.Put("/", s.putDistro())
				r.Patch("/", s.patchDistro())
				r.Delete("/", s.deleteDistro())
			})
		})

		r.Route("/profiles", func(r chi.Router) {
			r.Get("/", s.getProfiles())
			r.Post("/", s.createProfile())

			r.Route("/{profileID}", func(r chi.Router) {
				r.Get("/", s.getProfile())
				r.Put("/", s.putProfile())
				r.Patch("/", s.patchProfile())
				r.Delete("/", s.deleteProfile())
			})
		})

		r.Route("/systems", func(r chi.Router) {
			r.Get("/", s.getSystems())
			r.Post("/", s.createSystem())

			r.Route("/{systemID}", func(r chi.Router) {
				r.Get("/", s.getSystem())
				r.Put("/", s.putSystem())
				r.Patch("/", s.patchSystem())
				r.Delete("/", s.deleteSystem())
			})
		})

		r.Get("/pxe-config", s.getPxeConfig())
	})
}
