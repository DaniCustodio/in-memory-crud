package api

import (
	"encoding/json"
	"main/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteUser(t *testing.T) {
	t.Run("delete a user successfully", func(t *testing.T) {
		db := setupDB()

		users := db.FindAll()

		rec := makeDeleteRequest(db, users[0].ID.String())

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code to be %d, got %d", http.StatusOK, rec.Code)
		}

		got := parseData[database.DBUser](t, response.Data)

		if got.ID != users[0].ID {
			t.Fatalf("expected user id to be %s, got %s", users[0].ID, got.ID)
		}
	})

	t.Run("delete a user that doesn't exists", func(t *testing.T) {
		db := setupDB()

		rec := makeDeleteRequest(db, database.ID{}.NewID().String())

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status code to be %d, got %d", http.StatusNotFound, rec.Code)
		}

		if response.Message != "The user with the specified ID does not exist" {
			t.Fatalf("expected message to be %s, got %s", "The user with the specified ID does not exist", response.Message)
		}
	})
}

func makeDeleteRequest(db *database.InMemoryDB, userID string) *httptest.ResponseRecorder {
	router := NewHandler(db)

	req := httptest.NewRequest(http.MethodDelete, "/api/users/"+userID, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}
