package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

// GetUUIDFromRequest gets and parses the UUID from the request. If it's not a valid UUID, an error is returned.
func GetUUIDFromRequest(r *http.Request) (uuid.UUID, error) {
	uuidString := chi.URLParam(r, "uuid")
	UUID, err := uuid.Parse(uuidString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("[%s] is not a valid UUID: [%w]", uuidString, err)
	}
	return UUID, nil
}
