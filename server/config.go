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
	httpsEnabled  bool
	httpsCertFile string
	httpsKeyFile  string
	listenAddress string
}

func NewAppConfig() (AppConfig, error) {
	var a AppConfig
	var err error

	// Default values if applicable
	a.dbPort = 5432

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

	a.httpsCertFile = os.Getenv("GOBBLE_HTTPS_CERT_FILE")
	a.httpsKeyFile = os.Getenv("GOBBLE_HTTPS_KEY_FILE")

	listenAddress := os.Getenv("GOBBLE_LISTEN_ADDRESS")
	if listenAddress != "" {
		a.listenAddress = listenAddress
	}

	// Parse command line flags
	flag.StringVar(&a.dbUser, "db-user", a.dbUser, "the database user")
	flag.StringVar(&a.dbPass, "db-pass", a.dbPass, "the database password")
	flag.StringVar(&a.dbHost, "db-host", a.dbHost, "the database host")
	flag.StringVar(&a.dbName, "db-name", a.dbName, "the database to use")
	flag.IntVar(&a.dbPort, "db-port", a.dbPort, "the database port to connect to")
	flag.StringVar(&a.httpsCertFile, "https-cert-file", a.httpsCertFile, "the TLS certificate file to use for HTTPS")
	flag.StringVar(&a.httpsKeyFile, "https-key-file", a.httpsKeyFile, "the TLS certificate key file to use for HTTPS")
	flag.StringVar(&a.listenAddress, "listen-address", a.listenAddress, "the address that the application should listen on")
	flag.Parse()

	if a.dbUser == "" || a.dbPass == "" || a.dbHost == "" || a.dbName == "" {
		return a, ErrIncompleteDatabaseCredentials
	}

	// If both a certificate and corresponding key file path have been passed, HTTPS will be enabled
	if a.httpsCertFile != "" && a.httpsKeyFile != "" {
		a.httpsEnabled = true
	} else {
		a.httpsEnabled = false
	}

	// If no listen address has been passed, set an appropriate default depending on whether HTTPS has been enabled
	if a.listenAddress == "" {
		if a.httpsEnabled {
			a.listenAddress = ":443"
		} else {
			a.listenAddress = ":80"
		}
	}

	return a, nil
}
