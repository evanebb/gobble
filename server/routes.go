package server

import (
	"github.com/evanebb/gobble/api/auth"
	"github.com/evanebb/gobble/resources"
	"github.com/evanebb/gobble/server/handlers"
	"github.com/evanebb/gobble/server/handlers/api_handlers"
	"github.com/evanebb/gobble/server/handlers/ui_handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (s *Server) routes() {
	s.router.Use(handlers.MethodOverride, middleware.Logger)

	// API route group
	s.router.Route("/api", func(r chi.Router) {
		r.Use(auth.ApiBasicAuth(s.apiUserRepo))
		r.NotFound(api_handlers.ErrorHandler(api_handlers.UnknownEndpointHandler))

		r.Route("/profiles", func(r chi.Router) {
			h := api_handlers.NewProfileHandlerGroup(s.profileRepo)

			r.Get("/", api_handlers.ErrorHandler(h.GetProfiles))
			r.Post("/", api_handlers.ErrorHandler(h.CreateProfile))
			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", api_handlers.ErrorHandler(h.GetProfile))
				r.Put("/", api_handlers.ErrorHandler(h.PutProfile))
				r.Patch("/", api_handlers.ErrorHandler(h.PatchProfile))
				r.Delete("/", api_handlers.ErrorHandler(h.DeleteProfile))
			})
		})

		r.Route("/systems", func(r chi.Router) {
			h := api_handlers.NewSystemHandlerGroup(s.systemRepo)

			r.Get("/", api_handlers.ErrorHandler(h.GetSystems))
			r.Post("/", api_handlers.ErrorHandler(h.CreateSystem))
			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", api_handlers.ErrorHandler(h.GetSystem))
				r.Put("/", api_handlers.ErrorHandler(h.PutSystem))
				r.Patch("/", api_handlers.ErrorHandler(h.PatchSystem))
				r.Delete("/", api_handlers.ErrorHandler(h.DeleteSystem))
			})
		})

		r.Route("/users", func(r chi.Router) {
			h := api_handlers.NewApiUserHandlerGroup(s.apiUserRepo)

			r.Get("/", api_handlers.ErrorHandler(h.GetUsers))
			r.Post("/", api_handlers.ErrorHandler(h.CreateUser))
			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", api_handlers.ErrorHandler(h.GetUser))
				r.Put("/", api_handlers.ErrorHandler(h.PutUser))
				r.Delete("/", api_handlers.ErrorHandler(h.DeleteUser))
			})
		})
	})

	// This endpoint should not have authentication, so it lives outside the /api group above
	h := api_handlers.NewPxeConfigHandlerGroup(s.systemRepo, s.profileRepo)
	s.router.Get("/api/pxe-config", api_handlers.ErrorHandler(h.GetPxeConfig))

	// Redirect the index to the UI by default
	s.router.Handle("/", http.RedirectHandler("ui/", http.StatusMovedPermanently))

	// Front-end (UI) routes
	s.router.Route("/ui/", func(r chi.Router) {
		r.Use(auth.BrowserBasicAuth(s.apiUserRepo))

		r.NotFound(ui_handlers.PageNotFound)
		r.Get("/", ui_handlers.HomePage)
		r.Handle("/static/*", http.StripPrefix("/ui/", http.FileServer(http.FS(resources.Static))))

		r.Route("/profiles", func(r chi.Router) {
			h := ui_handlers.NewUiProfileHandlerGroup(s.profileRepo)

			r.Get("/", h.Overview)
			r.Get("/create", h.Create)
			r.Post("/", h.Store)

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", h.Show)
				r.Get("/edit", h.Edit)
				r.Put("/", h.Update)
				r.Delete("/", h.Delete)
			})
		})

		r.Route("/systems", func(r chi.Router) {
			h := ui_handlers.NewUiSystemHandlerGroup(s.systemRepo, s.profileRepo)

			r.Get("/", h.Overview)
			r.Get("/create", h.Create)
			r.Post("/", h.Store)

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", h.Show)
				r.Get("/edit", h.Edit)
				r.Put("/", h.Update)
				r.Delete("/", h.Delete)
			})
		})
	})
}
