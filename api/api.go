package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type ID uuid.UUID

type User struct {
	ID        ID
	FirstName string
	LastName  string
	Biography string
}

type CreateUserBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func NewHandler() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/api/users", handleCreateUser)

	return router
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: validate the request body
	var body CreateUserBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendJSON(
			w,
			Response{Error: "could not decode the request"},
			http.StatusBadRequest,
		)
		return
	}

	userId := ID(uuid.New())
	user := User{
		ID:        userId,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Biography: body.Biography,
	}

	// TODO: save the user to the database
	// TODO: return the user in the response

	sendJSON(
		w,
		Response{Data: user},
		http.StatusCreated,
	)
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("could not marshal the response", "error", err)
		sendJSON(
			w,
			Response{Error: "internal server error"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("could not write the response", "error", err)
		return
	}
}
