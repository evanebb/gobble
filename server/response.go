package server

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

// SendSuccessResponse sends a JSend-compliant response indicating success, with the passed data nested in it.
func SendSuccessResponse(w http.ResponseWriter, code int, v any) error {
	r := response{
		Status: "success",
		Data:   v,
	}
	return SendJSONResponse(w, code, r)
}

// SendErrorResponse sends a JSend-compliant response indicating an error with the passed error message.
func SendErrorResponse(w http.ResponseWriter, code int, message string) error {
	r := response{
		Status:  "error",
		Message: message,
	}
	return SendJSONResponse(w, code, r)
}

// SendFailResponse sends a JSend-compliant response indicating failure with the passed error message.
func SendFailResponse(w http.ResponseWriter, code int, message string) error {
	r := response{
		Status:  "fail",
		Message: message,
	}
	return SendJSONResponse(w, code, r)
}

func SendJSONResponse(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func SendPlainTextResponse(w http.ResponseWriter, code int, v string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, v)
	return err
}
