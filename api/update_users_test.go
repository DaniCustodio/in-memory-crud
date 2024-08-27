package api

import (
	"encoding/json"
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

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		got := parseData[database.DBUser](t, response.Data)

		if got.ID != users[0].ID {
			t.Errorf("expected id %d, got %d", users[0].ID, got.ID)
		}
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

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		if response.Message != "Please provide name and bio for the user" {
			t.Fatalf("expected message to be %s, got %s", "Please provide name and bio for the user", response.Message)
		}
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

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		if response.Message != "The user with the specified ID does not exist" {
			t.Fatalf("expected message to be %s, got %s", "The user with the specified ID does not exist", response.Message)
		}
	})
}
