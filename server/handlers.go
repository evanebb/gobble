package server

import "net/http"

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func errorHandler(h ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			// FIXME: add more specific HTTP error that allows setting status code
			w.WriteHeader(http.StatusInternalServerError)
			if err := SendJSONResponse(w, err.Error()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
