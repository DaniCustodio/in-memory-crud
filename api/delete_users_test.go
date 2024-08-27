package api

import (
	"main/database"
	"net/http"
	"testing"
)

func TestDeleteUser(t *testing.T) {
	const URL = "/api/users/"

	t.Run("delete a user successfully", func(t *testing.T) {
		db := setupDB()

		users := db.FindAll()

		request, err := createRequest(
			http.MethodDelete,
			URL+users[0].ID.String(),
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

		assertUser(t, users[0], response.Data)
	})

	t.Run("delete a user that doesn't exists", func(t *testing.T) {
		db := setupDB()

		request, err := createRequest(
			http.MethodDelete,
			URL+database.ID{}.NewID().String(),
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, request)

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusNotFound, rec.Code)

		assertErrorMessage(t, "The user with the specified ID does not exist", response.Message)
	})
}
