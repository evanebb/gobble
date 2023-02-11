package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendJSONResponse(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func SendPlainTextResponse(w http.ResponseWriter, v string) error {
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprint(w, v)
	return err
}
