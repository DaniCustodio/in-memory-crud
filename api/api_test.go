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
		user := CreateUserBody{
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

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		dataBytes, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("could not marshal the data: %v", err)
		}

		var got User
		if err := json.Unmarshal(dataBytes, &got); err != nil {
			t.Fatalf("could not unmarshal the user: %v", err)
		}

		if got.FirstName != user.FirstName {
			t.Errorf("expected first name %q; got %q", user.FirstName, got.FirstName)
		}

		if got.LastName != user.LastName {
			t.Errorf("expected last name %q; got %q", user.LastName, got.LastName)
		}

		if got.Biography != user.Biography {
			t.Errorf("expected biography %q; got %q", user.Biography, got.Biography)
		}

		if got.ID == ID(uuid.Nil) {
			t.Error("expected a non-empty ID")
		}

	})
}
