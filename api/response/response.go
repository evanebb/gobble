package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Status  string `json:"status"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

// Success sends a JSend-compliant response indicating success, with the passed data nested in it.
func Success(w http.ResponseWriter, code int, v any) error {
	r := response{
		Status: "success",
		Data:   v,
	}
	return JSON(w, code, r)
}

// Error sends a JSend-compliant response indicating an error with the passed error message.
func Error(w http.ResponseWriter, code int, message string) error {
	r := response{
		Status:  "error",
		Message: message,
	}
	return JSON(w, code, r)
}

// Fail sends a JSend-compliant response indicating failure with the passed error message.
func Fail(w http.ResponseWriter, code int, message string) error {
	r := response{
		Status:  "fail",
		Message: message,
	}
	return JSON(w, code, r)
}

func JSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func PlainText(w http.ResponseWriter, code int, v string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, v)
	return err
}

func HTML(w http.ResponseWriter, code int, v string) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, v)
	return err
}
