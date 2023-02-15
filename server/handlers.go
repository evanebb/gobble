package server

import (
	"errors"
	"log"
	"net/http"
)

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
			http.Error(w, "Something really bad happened, please check the logs!", http.StatusInternalServerError)
		}
	}
}
