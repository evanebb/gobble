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
	"os"
	"strconv"
)

type Server struct {
	apiUserRepo auth.ApiUserRepository
	profileRepo profile.Repository
	systemRepo  system.Repository
	router      chi.Router
}

func NewServer() (Server, error) {
	var s Server

	// FIXME: proper configuration handling
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")
	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(err)
	}

	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, pass, host, port, database)
	db, err := pgxpool.New(context.Background(), cs)
	if err != nil {
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
	log.Fatal(http.ListenAndServe(":80", s.router))
}
