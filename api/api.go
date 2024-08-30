package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"main/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Response[T any] struct {
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

var ErrInvalidUserParams = errors.New("please provide a valid FirstName, LastName and Bio for the user")
var ErrUserNotFound = errors.New("the user with the specified ID does not exist")
var ErrInvalidUpdateUserParams = errors.New("please provide name and bio for the user")

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
	router.Delete("/api/users/{id}", handleDeleteUser(db))
	router.Put("/api/users/{id}", handleUpdateUser(db))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return router
}

// GetUser godoc
//
//	@Summary		Get a user by ID
//	@Description	Get a user by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	Response[database.DBUser]{data=database.DBUser}
//	@Failure		404	{object}	Response[any]{message=string}
//	@Router			/users/{id} [get]
func handleGetUser(db *database.InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		user, exists := db.FindByID(id)
		if !exists {
			sendJSON(
				w,
				Response[any]{Message: ErrUserNotFound.Error()},
				http.StatusNotFound,
			)
			return
		}

		sendJSON(
			w,
			Response[database.DBUser]{Data: user},
			http.StatusOK,
		)
	}
}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Get all users
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response[[]database.DBUser]{data=[]database.DBUser}
//	@Router			/users [get]
func handleGetUsers(db *database.InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := db.FindAll()

		sendJSON(
			w,
			Response[[]database.DBUser]{Data: users},
			http.StatusOK,
		)
	}
}

// CreateUser godoc
//
//	@Summary		Create a user
//	@Description	Create a user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		database.User	true	"User details"
//	@Success		201		{object}	Response[database.DBUser]{data=database.DBUser}
//	@Failure		400		{object}	Response[any]{message=string}
//	@Router			/users [post]
func handleCreateUser(db *database.InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body database.User
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(
				w,
				Response[any]{Message: "could not decode the request"},
				http.StatusBadRequest,
			)
			return
		}

		if err := validate.Struct(&body); err != nil {
			sendJSON(
				w,
				Response[any]{Message: ErrInvalidUserParams.Error()},
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
			Response[database.DBUser]{Data: dbUser},
			http.StatusCreated,
		)
	}
}

// DeleteUser godoc
//
//	@Summary		Delete a user by ID
//	@Description	Delete a user by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	Response[database.DBUser]{data=database.DBUser}
//	@Failure		404	{object}	Response[any]{message=string}
//	@Router			/users/{id} [delete]
func handleDeleteUser(db *database.InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		user, err := db.Delete(id)
		if err != nil {
			sendJSON(
				w,
				Response[any]{Message: ErrUserNotFound.Error()},
				http.StatusNotFound,
			)
		}

		sendJSON(
			w,
			Response[database.DBUser]{Data: user},
			http.StatusOK,
		)
	}
}

// UpdateUser godoc
//
//	@Summary		Update a user by ID
//	@Description	Update a user by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"User ID"
//	@Param			body	body		database.User	true	"User details"
//	@Success		200		{object}	Response[database.DBUser]{data=database.DBUser}
//	@Failure		400		{object}	Response[any]{message=string}
//	@Failure		404		{object}	Response[any]{message=string}
//	@Router			/users/{id} [put]
func handleUpdateUser(db *database.InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var body database.User
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(
				w,
				Response[any]{Message: "could not decode the request"},
				http.StatusBadRequest,
			)
			return
		}

		if err := validate.Struct(&body); err != nil {
			sendJSON(
				w,
				Response[any]{Message: ErrInvalidUpdateUserParams.Error()},
				http.StatusBadRequest,
			)
			return
		}

		user, err := db.Update(id, body)
		if err != nil {
			sendJSON(
				w,
				Response[database.DBUser]{Message: ErrUserNotFound.Error()},
				http.StatusNotFound,
			)
			return
		}

		sendJSON(
			w,
			Response[database.DBUser]{Data: user},
			http.StatusOK,
		)
	}
}

func sendJSON[T any](w http.ResponseWriter, resp Response[T], status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("could not marshal the response", "error", err)
		sendJSON(
			w,
			Response[any]{Message: "internal server error"},
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
