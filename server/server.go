package server

import (
	"github.com/evanebb/gobble/distro"
	"github.com/evanebb/gobble/profile"
	"github.com/evanebb/gobble/repository/postgres"
	"github.com/evanebb/gobble/system"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

type Server struct {
	distroRepo  distro.Repository
	profileRepo profile.Repository
	systemRepo  system.Repository
	router      chi.Router
}

func NewServer(db *pgxpool.Pool, router chi.Router) (Server, error) {
	var s Server

	// FIXME: don't instantiate the repositories here?
	dr, err := postgres.NewDistroRepository(db)
	if err != nil {
		return s, err
	}

	pr, err := postgres.NewProfileRepository(db)
	if err != nil {
		return s, err
	}

	sr, err := postgres.NewSystemRepository(db)
	if err != nil {
		return s, err
	}

	s.distroRepo = dr
	s.profileRepo = pr
	s.systemRepo = sr
	s.router = router
	return s, nil
}

func (s *Server) Run() {
	s.routes()
	log.Fatal(http.ListenAndServe(":8080", s.router))
}
