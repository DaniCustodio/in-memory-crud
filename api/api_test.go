package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	t.Run("create a user successfully", func(t *testing.T) {
		user := User{
			ID:        ID(uuid.New()),
			FirstName: "John",
			LastName:  "Doe",
			Biography: "A regular guy",
		}

		payload, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("could not marshal the user: %v", err)
		}

		router := NewHandler()

		req, err := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("could not create a request: %v", err)
		}

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201; got %d", rec.Code)
		}

		var got User
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if got != user {
			t.Errorf("expected the same user; got %v", got)
		}

	})
}
