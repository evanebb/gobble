package auth

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(struct {
		Status  string `json:"status"`
		Data    any    `json:"data"`
		Message string `json:"message"`
	}{
		Status:  "error",
		Message: "authentication failed",
	})
}
