package api

import (
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

		response, err := parseResponse[database.DBUser](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusOK, rec.Code)

		if response.Data.ID != users[0].ID {
			t.Fatalf("expected user id to be %s, got %s", users[0].ID, response.Data.ID)
		}
	})

	t.Run("delete a user that doesn't exists", func(t *testing.T) {
		db := setupDB()

		rec := makeDeleteRequest(db, database.ID{}.NewID().String())

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusNotFound, rec.Code)

		assertErrorMessage(t, "The user with the specified ID does not exist", response.Message)
	})
}

func makeDeleteRequest(db *database.InMemoryDB, userID string) *httptest.ResponseRecorder {
	router := NewHandler(db)

	req := httptest.NewRequest(http.MethodDelete, "/api/users/"+userID, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}
