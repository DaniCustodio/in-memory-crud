package api

import (
	"encoding/json"
	"log/slog"
	"main/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// WithRequiredStructEnabled: opt-in to new behavior that will become the default behavior in v11+
var validate = validator.New(validator.WithRequiredStructEnabled())

func NewHandler(db *database.InMemoryDB) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/api/users", handleCreateUser(db))
	router.Get("/api/users", handleGetUsers(db))
	router.Get("/api/users/{id}", handleGetUser(db))

	return router
}

func handleGetUser(db *database.InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		user, exists := db.FindByID(id)
		if !exists {
			sendJSON(
				w,
				Response{Message: "The user with the specified ID does not exist"},
				http.StatusNotFound,
			)
			return
		}

		sendJSON(
			w,
			Response{Data: user},
			http.StatusOK,
		)
	}
}

func handleGetUsers(db *database.InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := db.FindAll()

		sendJSON(
			w,
			Response{Data: users},
			http.StatusOK,
		)
	}
}

func handleCreateUser(db *database.InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body database.User
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(
				w,
				Response{Message: "could not decode the request"},
				http.StatusBadRequest,
			)
			return
		}

		if err := validate.Struct(&body); err != nil {
			sendJSON(
				w,
				Response{Message: "Please provide a valid FirstName, LastName and Bio for the user"},
				http.StatusBadRequest,
			)
			return
		}

		user := database.User{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Biography: body.Biography,
		}

		dbUser := db.Insert(user)

		sendJSON(
			w,
			Response{Data: dbUser},
			http.StatusCreated,
		)
	}
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("could not marshal the response", "error", err)
		sendJSON(
			w,
			Response{Message: "internal server error"},
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
