package auth

import (
	"fmt"
	"github.com/evanebb/gobble/api/response"
	"net/http"
)

// ApiBasicAuth will check the basic auth credentials sent in the request against the known users,
// and return a JSON response if authentication has failed.
func ApiBasicAuth(db ApiUserRepository) func(next http.Handler) http.Handler {
	return basicAuth(db, sendBasicAuthFailedResponse)
}

// BrowserBasicAuth will check the basic auth credentials sent in the request against the known users,
// and (re)-request basic auth credentials if authentication has failed.
func BrowserBasicAuth(db ApiUserRepository) func(next http.Handler) http.Handler {
	return basicAuth(db, requestBasicAuth)
}

// basicAuth will check the basic auth credentials sent in the request against the known users,
// and execute the passed callback if authentication has failed.
func basicAuth(db ApiUserRepository, authFailureCallback func(w http.ResponseWriter)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				authFailureCallback(w)
				return
			}

			apiUser, err := db.GetApiUserByName(user)
			if err != nil {
				authFailureCallback(w)
				return
			}

			err = apiUser.CheckPassword(pass)
			if err != nil {
				authFailureCallback(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// sendBasicAuthFailedResponse will write a JSON response indicating authentication failure to the passed http.ResponseWriter variable.
func sendBasicAuthFailedResponse(w http.ResponseWriter) {
	err := response.Error(w, http.StatusUnauthorized, "authentication failed")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// sendBasicAuthFailedResponse will request basic authentication by sending the 'WWW-Authenticate' header with a 401 status code.
func requestBasicAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="gobble"`)
	w.WriteHeader(401)
	_, _ = fmt.Fprint(w, "Unauthorised")
}
