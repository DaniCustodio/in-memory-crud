package api

import (
	"main/database"
	"net/http"
	"testing"
)

var users []database.User = []database.User{
	{

		FirstName: "John",
		LastName:  "Doe",
		Biography: "A simple guy who loves to write code and play games. He is a fan of technology and loves to read about new things.",
	},
	{
		FirstName: "Jane",
		LastName:  "Doe",
		Biography: "A nice lady who loves to write code and play games. She is a fan of technology and loves to read about new things.",
	},
}

func TestGetUsers(t *testing.T) {
	const URL = "/api/users"

	t.Run("get list of users", func(t *testing.T) {
		db := setupDB()

		request, err := createRequest(
			http.MethodGet,
			URL,
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, request)

		response, err := parseResponse[[]database.DBUser](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusOK, rec.Code)

		if len(response.Data) != len(users) {
			t.Fatalf("expected %d users, got %d", len(users), len(response.Data))
		}
	})

	t.Run("get user by ID", func(t *testing.T) {
		db := setupDB()

		dbUsers := db.FindAll()

		request, err := createRequest(
			http.MethodGet,
			URL+"/"+dbUsers[0].ID.String(),
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, request)

		response, err := parseResponse[database.DBUser](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusOK, rec.Code)

		assertUser(t, dbUsers[0], response.Data)
	})

	t.Run("get user by ID that does not exist", func(t *testing.T) {
		db := setupDB()

		request, err := createRequest(
			http.MethodGet,
			URL+"/"+database.ID{}.NewID().String(),
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, request)

		assertStatusCode(t, http.StatusNotFound, rec.Code)

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertErrorMessage(t, ErrUserNotFound.Error(), response.Message)
	})
}
