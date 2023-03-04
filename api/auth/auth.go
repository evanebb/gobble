package auth

import (
	"github.com/evanebb/gobble/api/response"
	"net/http"
)

func BasicAuth(db ApiUserRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				basicAuthFailed(w)
				return
			}

			apiUser, err := db.GetApiUserByName(user)
			if err != nil {
				basicAuthFailed(w)
				return
			}

			err = apiUser.CheckPassword(pass)
			if err != nil {
				basicAuthFailed(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func basicAuthFailed(w http.ResponseWriter) {
	err := response.Error(w, http.StatusUnauthorized, "authentication failed")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
