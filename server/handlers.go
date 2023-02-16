package server

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
)

var fatalErrorMsg = "something really bad happened, please check the logs!"

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func errorHandler(h ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		var httpErr HTTPError
		var statusCode int

		if errors.As(err, &httpErr) {
			statusCode = httpErr.StatusCode
		} else {
			statusCode = http.StatusInternalServerError
		}

		// FIXME: probably want to log more than just 500's, this avoids log spam for bad request errors for now
		if statusCode == http.StatusInternalServerError {
			log.Println(err)
		}

		if err := SendErrorResponse(w, statusCode, err.Error()); err != nil {
			// This shouldn't ever happen, if it does just return a bogus response?
			// We don't actually know whether a response has been written at this point; let's hope net/http handles that ;)
			log.Println(err)
			http.Error(w, fatalErrorMsg, http.StatusInternalServerError)
		}
	}
}

func unknownEndpointHandler(w http.ResponseWriter, r *http.Request) {
	if err := SendErrorResponse(w, http.StatusNotFound, "unknown endpoint, please refer to the documentation for available endpoints"); err != nil {
		log.Println(err)
		http.Error(w, fatalErrorMsg, http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	html := "<h1>Welcome to gobble!</h1><p>Refer to the documentation for the available API endpoints.</p>"
	if err := SendHTMLResponse(w, http.StatusOK, html); err != nil {
		log.Println(err)
		http.Error(w, fatalErrorMsg, http.StatusInternalServerError)
	}
}

// getUUIDFromRequest gets and parses the UUID from the request. If it's not a valid UUID, an error is returned.
func getUUIDFromRequest(r *http.Request) (uuid.UUID, error) {
	uuidString := chi.URLParam(r, "uuid")
	UUID, err := uuid.Parse(uuidString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("[%s] is not a valid UUID: [%w]", uuidString, err)
	}
	return UUID, nil
}
