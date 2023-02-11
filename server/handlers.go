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
			if errors.As(err, &httpErr) {
				w.WriteHeader(httpErr.StatusCode)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			// FIXME: only log server errors, don't care about bad request errors
			log.Println(err)
			if err := SendJSONResponse(w, err.Error()); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
