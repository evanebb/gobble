package server

import (
	"github.com/evanebb/gobble/api/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) routes() {
	s.router.Use(middleware.Logger)

	s.router.Get("/", errorHandler(indexHandler))
	s.router.Route("/api", func(r chi.Router) {
		r.Use(auth.BasicAuth(s.apiUserRepo))
		r.NotFound(errorHandler(unknownEndpointHandler))

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

		r.Route("/users", func(r chi.Router) {
			r.Get("/", errorHandler(s.getUsers))
			r.Post("/", errorHandler(s.createUser))

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", errorHandler(s.getUser))
				r.Put("/", errorHandler(s.putUser))
				r.Delete("/", errorHandler(s.deleteUser))
			})
		})
	})
}
