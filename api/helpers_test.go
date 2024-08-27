package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createRequest(method string, url string, body any) (*http.Request, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("could not marshal the body: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("could not create a request: %w", err)
	}

	return req, nil
}

func makeRequest(db *database.InMemoryDB, request *http.Request) *httptest.ResponseRecorder {
	router := NewHandler(db)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)

	return rec
}

func assertStatusCode(t testing.TB, want int, got int) {
	t.Helper()

	if got != want {
		t.Errorf("expected status code %d, got %d", want, got)
	}
}

func assertErrorMessage(t testing.TB, want string, got string) {
	t.Helper()

	if got != want {
		t.Errorf("expected error message %q, got %q", want, got)
	}
}

func assertUser(t testing.TB, want database.DBUser, got database.DBUser) {
	t.Helper()

	if got.ID.IsEmpty() ||
		got.User.FirstName != want.User.FirstName ||
		got.User.LastName != want.User.LastName ||
		got.User.Biography != want.User.Biography {
		t.Errorf("expected user %v, got %v", want, got)
	}
}

func parseResponse[T any](response *httptest.ResponseRecorder) (Response[T], error) {
	var parsedResponse Response[T]
	if err := json.NewDecoder(response.Body).Decode(&parsedResponse); err != nil {
		return Response[T]{}, fmt.Errorf("could not decode the response: %w", err)
	}

	var parsedData T

	bytes, err := json.Marshal(parsedResponse.Data)
	if err != nil {
		return Response[T]{}, fmt.Errorf("could not marshal the data: %w", err)
	}

	if err := json.Unmarshal(bytes, &parsedData); err != nil {
		return Response[T]{}, fmt.Errorf("could not unmarshal the data: %w", err)
	}

	parsedResponse.Data = parsedData

	return parsedResponse, nil
}

func setupDB() *database.InMemoryDB {
	db := database.NewInMemoryDB()
	db.Insert(users[0])
	db.Insert(users[1])

	return db
}
