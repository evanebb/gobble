package server

import (
	"context"
	"fmt"
	"github.com/evanebb/gobble/api/auth"
	"github.com/evanebb/gobble/profile"
	"github.com/evanebb/gobble/repository/postgres"
	"github.com/evanebb/gobble/system"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

type Server struct {
	apiUserRepo auth.ApiUserRepository
	profileRepo profile.Repository
	systemRepo  system.Repository
	router      chi.Router
	config      AppConfig
}

func NewServer() (Server, error) {
	var s Server

	var err error
	s.config, err = NewAppConfig()
	if err != nil {
		return s, err
	}

	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", s.config.dbUser, s.config.dbPass, s.config.dbHost, s.config.dbPort, s.config.dbName)
	db, err := pgxpool.New(context.Background(), cs)
	if err != nil {
		return s, err
	}

	if err = db.Ping(context.Background()); err != nil {
		return s, err
	}

	ar, err := postgres.NewApiUserRepository(db)
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

	router := chi.NewRouter()

	s.apiUserRepo = ar
	s.profileRepo = pr
	s.systemRepo = sr
	s.router = router
	return s, nil
}

func (s *Server) Run() {
	log.Println("starting API...")
	s.routes()
	log.Fatal(http.ListenAndServe(s.config.listenAddress, s.router))
}
