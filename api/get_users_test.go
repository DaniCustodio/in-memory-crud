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

	t.Run("it should return a list of users", func(t *testing.T) {
		rec := makeGetRequest()

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
}

func makeGetRequest() *httptest.ResponseRecorder {
	db := database.NewInMemoryDB()
	db.Insert(users[0])
	db.Insert(users[1])

	router := NewHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}
