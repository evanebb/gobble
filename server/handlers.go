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
		if err != nil {
			var httpErr HTTPError
			var statusCode int
			if errors.As(err, &httpErr) {
				statusCode = httpErr.StatusCode
			} else {
				statusCode = http.StatusInternalServerError
			}

			w.WriteHeader(statusCode)
			if statusCode == http.StatusInternalServerError {
				log.Println(err)
			}
			if err := SendJSONResponse(w, err.Error()); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
