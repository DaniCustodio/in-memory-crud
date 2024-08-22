package api

import (
	"encoding/json"
	"main/database"
	"net/http"
	"net/http/httptest"
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

	t.Run("get list of users", func(t *testing.T) {
		db := setupDB()

		rec := makeGetRequest(db, "")

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code to be %d, got %d", http.StatusOK, rec.Code)
		}

		data, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("could not marshal the data: %v", err)
		}

		var got []database.User
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("could not unmarshal the user: %v", err)
		}

		if len(got) != len(users) {
			t.Fatalf("expected %d users, got %d", len(users), len(got))
		}
	})

	t.Run("get user by ID", func(t *testing.T) {
		db := setupDB()

		dbUsers := db.FindAll()

		rec := makeGetRequest(db, dbUsers[0].ID.String())

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code to be %d, got %d", http.StatusOK, rec.Code)
		}

		data, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("could not marshal the data: %v", err)
		}

		var got database.DBUser
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("could not unmarshal the user: %v", err)
		}

		if got != dbUsers[0] {
			t.Fatalf("expected user to be %v, got %v", dbUsers[0], got)
		}

	})

	t.Run("get user by ID that does not exist", func(t *testing.T) {
		db := setupDB()

		rec := makeGetRequest(db, database.ID{}.NewID().String())

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status code to be %d, got %d", http.StatusNotFound, rec.Code)
		}

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if response.Message != "The user with the specified ID does not exist" {
			t.Fatalf("expected message to be %s, got %s", "The user with the specified ID does not exist", response.Message)
		}
	})
}

func makeGetRequest(db *database.InMemoryDB, userID string) *httptest.ResponseRecorder {
	router := NewHandler(db)

	var url string
	if userID == "" {
		url = "/api/users"
	} else {
		url = "/api/users/" + userID
	}

	req := httptest.NewRequest(http.MethodGet, url, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func setupDB() *database.InMemoryDB {
	db := database.NewInMemoryDB()
	db.Insert(users[0])
	db.Insert(users[1])

	return db
}
