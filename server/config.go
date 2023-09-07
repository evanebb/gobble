package server

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

var ErrIncompleteDatabaseCredentials = errors.New("incomplete or no database credentials supplied")

type AppConfig struct {
	dbUser        string
	dbPass        string
	dbHost        string
	dbName        string
	dbPort        int
	listenAddress string
}

func NewAppConfig() (AppConfig, error) {
	var a AppConfig
	var err error

	// Default values if applicable
	a.dbPort = 5432
	a.listenAddress = ":80"

	// Parse environment variables
	a.dbUser = os.Getenv("GOBBLE_DB_USER")
	a.dbPass = os.Getenv("GOBBLE_DB_PASS")
	a.dbHost = os.Getenv("GOBBLE_DB_HOST")
	a.dbName = os.Getenv("GOBBLE_DB_NAME")
	portString := os.Getenv("GOBBLE_DB_PORT")
	if portString != "" {
		a.dbPort, err = strconv.Atoi(portString)
		if err != nil {
			return a, err
		}
	}

	listenAddress := os.Getenv("GOBBLE_LISTEN_ADDRESS")
	if listenAddress != "" {
		a.listenAddress = listenAddress
	}

	// Parse command line flags
	flag.StringVar(&a.dbUser, "db-user", a.dbUser, "the database user")
	flag.StringVar(&a.dbPass, "db-pass", a.dbPass, "the database password")
	flag.StringVar(&a.dbHost, "db-host", a.dbHost, "the database host")
	flag.StringVar(&a.dbName, "db-name", a.dbHost, "the database to use")
	flag.IntVar(&a.dbPort, "db-port", a.dbPort, "the database port to connect to")
	flag.StringVar(&a.listenAddress, "listen-address", a.listenAddress, "the address that the application should listen on")
	flag.Parse()

	if a.dbUser == "" || a.dbPass == "" || a.dbHost == "" || a.dbName == "" {
		return a, ErrIncompleteDatabaseCredentials
	}

	return a, nil
}
