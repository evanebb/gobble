package server

import (
	"errors"
	"fmt"
	"github.com/evanebb/gobble/api/response"
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
		var resp string

		if errors.As(err, &httpErr) {
			statusCode = httpErr.StatusCode
		} else {
			statusCode = http.StatusInternalServerError
		}

		// For server-side errors, return a generic message and log the error; I don't want to expose potentially sensitive information from the error to the client.
		// I don't care about logging client errors (e.g. bad requests), the error message should be descriptive enough for them to figure it out themselves.
		if statusCode >= 500 && statusCode <= 599 {
			log.Println(err)
			resp = fatalErrorMsg
		} else {
			resp = err.Error()
		}

		if err := response.Error(w, statusCode, resp); err != nil {
			// This shouldn't ever happen, if it does just return a bogus response?
			// I don't actually know whether a response has been written at this point; let's hope net/http handles that ;)
			log.Println(err)
			http.Error(w, fatalErrorMsg, http.StatusInternalServerError)
		}
	}
}

func unknownEndpointHandler(w http.ResponseWriter, r *http.Request) error {
	return response.Error(w, http.StatusNotFound, "unknown endpoint, please refer to the documentation for available endpoints")
}

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	html := "<h1>Welcome to gobble!</h1><p>Refer to the documentation for the available API endpoints.</p>"
	return response.HTML(w, http.StatusOK, html)
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
