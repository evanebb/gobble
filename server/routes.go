package server

import (
	"github.com/evanebb/gobble/api/auth"
	"github.com/evanebb/gobble/server/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) routes() {
	s.router.Use(middleware.Logger)

	s.router.Get("/", handlers.ErrorHandler(handlers.IndexHandler))
	s.router.Route("/api", func(r chi.Router) {
		r.Use(auth.BasicAuth(s.apiUserRepo))
		r.NotFound(handlers.ErrorHandler(handlers.UnknownEndpointHandler))

		r.Route("/profiles", func(r chi.Router) {
			h := handlers.NewProfileHandlerGroup(s.profileRepo)

			r.Get("/", handlers.ErrorHandler(h.GetProfiles))
			r.Post("/", handlers.ErrorHandler(h.CreateProfile))
			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", handlers.ErrorHandler(h.GetProfile))
				r.Put("/", handlers.ErrorHandler(h.PutProfile))
				r.Patch("/", handlers.ErrorHandler(h.PatchProfile))
				r.Delete("/", handlers.ErrorHandler(h.DeleteProfile))
			})
		})

		r.Route("/systems", func(r chi.Router) {
			h := handlers.NewSystemHandlerGroup(s.systemRepo)

			r.Get("/", handlers.ErrorHandler(h.GetSystems))
			r.Post("/", handlers.ErrorHandler(h.CreateSystem))
			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", handlers.ErrorHandler(h.GetSystem))
				r.Put("/", handlers.ErrorHandler(h.PutSystem))
				r.Patch("/", handlers.ErrorHandler(h.PatchSystem))
				r.Delete("/", handlers.ErrorHandler(h.DeleteSystem))
			})
		})

		r.Route("/users", func(r chi.Router) {
			h := handlers.NewApiUserHandlerGroup(s.apiUserRepo)

			r.Get("/", handlers.ErrorHandler(h.GetUsers))
			r.Post("/", handlers.ErrorHandler(h.CreateUser))
			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", handlers.ErrorHandler(h.GetUser))
				r.Put("/", handlers.ErrorHandler(h.PutUser))
				r.Delete("/", handlers.ErrorHandler(h.DeleteUser))
			})
		})
	})
	// This endpoint should not have authentication, so it lives outside the /api group above
	h := handlers.NewPxeConfigHandlerGroup(s.systemRepo, s.profileRepo)
	s.router.Get("/api/pxe-config", handlers.ErrorHandler(h.GetPxeConfig))
}
