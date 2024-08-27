package api

import (
	"main/database"
	"net/http"
	"testing"
)

func TestUpdateUser(t *testing.T) {
	const URL = "/api/users/"

	t.Run("update a user successfully", func(t *testing.T) {
		db := setupDB()

		users := db.FindAll()

		updatedUser := database.User{
			FirstName: "updated first name",
			LastName:  "updated last name",
			Biography: "updated biography updated biography updated biography updated biography updated biography updated biography updated biography updated biography updated biography",
		}

		req, err := createRequest(http.MethodPut, URL+users[0].ID.String(), updatedUser)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		response, err := parseResponse[database.DBUser](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusOK, rec.Code)

		assertUser(t, database.DBUser{ID: users[0].ID, User: updatedUser}, response.Data)
	})

	t.Run("update a user with invalid data", func(t *testing.T) {
		db := setupDB()

		users := db.FindAll()

		updatedUser := database.User{
			FirstName: "updated first name",
			LastName:  "updated last name",
		}

		req, err := createRequest(http.MethodPut, URL+users[0].ID.String(), updatedUser)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusBadRequest, rec.Code)

		assertErrorMessage(t, "Please provide name and bio for the user", response.Message)
	})

	t.Run("update a user with invalid id", func(t *testing.T) {
		db := setupDB()

		updatedUser := database.User{
			FirstName: "updated first name",
			LastName:  "updated last name",
			Biography: "updated biography updated biography updated biography updated biography updated biography updated biography updated biography updated biography updated biography",
		}

		req, err := createRequest(http.MethodPut, URL+database.ID{}.NewID().String(), updatedUser)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusNotFound, rec.Code)

		assertErrorMessage(t, "The user with the specified ID does not exist", response.Message)
	})
}
