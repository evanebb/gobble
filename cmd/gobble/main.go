package main

import (
	"context"
	"fmt"
	"github.com/evanebb/gobble/server"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"strconv"
)

func main() {
	// FIXME: proper configuration handling
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(err)
	}
	database := os.Getenv("database")

	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, pass, host, port, database)
	db, err := pgxpool.New(context.Background(), cs)
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	s, err := server.NewServer(db, router)
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}
