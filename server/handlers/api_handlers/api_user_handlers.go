package api_handlers

import (
	"encoding/json"
	"errors"
	"github.com/evanebb/gobble/api/auth"
	"github.com/evanebb/gobble/api/response"
	"github.com/evanebb/gobble/repository"
	"github.com/evanebb/gobble/server/handlers"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

/*
 * Request and response structures, and their supporting functions
 */

// userRequest is the JSON representation of an auth.ApiUser that is accepted by the API.
type userRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// userResponse is the JSON representation of an auth.ApiUser that is returned by the API.
type userResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// newUserResponse accepts an auth.ApiUser, and casts it into a userResponse.
func newUserResponse(a auth.ApiUser) userResponse {
	return userResponse{
		Id:   a.Id,
		Name: a.Name,
	}
}

/*
 * HTTP handlers
 */

// ApiUserHandlerGroup is a group of http.HandlerFunc functions related to API users
type ApiUserHandlerGroup struct {
	apiUserRepo auth.ApiUserRepository
}

func NewApiUserHandlerGroup(ar auth.ApiUserRepository) ApiUserHandlerGroup {
	return ApiUserHandlerGroup{ar}
}

func (h ApiUserHandlerGroup) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.apiUserRepo.GetApiUsers()
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	resp := make([]userResponse, 0)
	for _, a := range users {
		resp = append(resp, newUserResponse(a))
	}

	return response.Success(w, http.StatusOK, resp)
}

func (h ApiUserHandlerGroup) GetUser(w http.ResponseWriter, r *http.Request) error {
	userID, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	a, err := h.apiUserRepo.GetApiUserById(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NewHTTPError(err, http.StatusNotFound)
		}
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusOK, newUserResponse(a))
}

func (h ApiUserHandlerGroup) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var req userRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	userID := uuid.New()

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	a := auth.NewApiUser(userID, req.Name, pass)
	err = h.apiUserRepo.SetApiUser(a)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusCreated, newUserResponse(a))
}

func (h ApiUserHandlerGroup) PutUser(w http.ResponseWriter, r *http.Request) error {
	var req userRequest

	userID, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	a := auth.NewApiUser(userID, req.Name, pass)
	err = h.apiUserRepo.SetApiUser(a)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusOK, newUserResponse(a))
}

func (h ApiUserHandlerGroup) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	userID, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		return NewHTTPError(err, http.StatusBadRequest)
	}

	err = h.apiUserRepo.DeleteApiUserById(userID)
	if err != nil {
		return NewHTTPError(err, http.StatusInternalServerError)
	}

	return response.Success(w, http.StatusNoContent, nil)
}
