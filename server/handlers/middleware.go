package handlers

import "net/http"

// MethodOverride will check for a hidden "_method" input in the POST form,
// so that PUT, PATCH and DELETE requests can be supported
func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			method := r.PostFormValue("_method")

			if method == "PUT" || method == "PATCH" || method == "DELETE" {
				r.Method = method
			}
		}

		next.ServeHTTP(w, r)
	})
}
