package api

import (
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

		response, err := parseResponse[[]database.DBUser](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusOK, rec.Code)

		if len(response.Data) != len(users) {
			t.Fatalf("expected %d users, got %d", len(users), len(response.Data))
		}
	})

	t.Run("get user by ID", func(t *testing.T) {
		db := setupDB()

		dbUsers := db.FindAll()

		rec := makeGetRequest(db, dbUsers[0].ID.String())

		response, err := parseResponse[database.DBUser](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusOK, rec.Code)

		if response.Data != dbUsers[0] {
			t.Fatalf("expected user to be %v, got %v", dbUsers[0], response.Data)
		}

	})

	t.Run("get user by ID that does not exist", func(t *testing.T) {
		db := setupDB()

		rec := makeGetRequest(db, database.ID{}.NewID().String())

		assertStatusCode(t, http.StatusNotFound, rec.Code)

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertErrorMessage(t, "The user with the specified ID does not exist", response.Message)
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
