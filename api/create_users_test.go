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

func TestCreateUser(t *testing.T) {
	const URL = "/api/users"

	requestBody := database.User{
		FirstName: "John",
		LastName:  "Doe",
		Biography: "A regular guy who loves to code in Go and JavaScript",
	}

	emptyDBUser := database.DBUser{}

	t.Run("create a user successfully", func(t *testing.T) {
		db := database.NewInMemoryDB()

		req, err := createRequest(http.MethodPost, URL, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		assertResponse(
			t,
			rec,
			http.StatusCreated,
			"",
			database.DBUser{
				ID:   database.ID{}.NewID(),
				User: requestBody,
			},
		)
	})

	t.Run("first name length should be >= 2", func(t *testing.T) {
		user := requestBody
		user.FirstName = ""

		db := database.NewInMemoryDB()

		req, err := createRequest(http.MethodPost, URL, user)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			emptyDBUser,
		)
	})

	t.Run("first name length should be <= 20", func(t *testing.T) {
		user := requestBody
		user.FirstName = "JohnDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoe"

		db := database.NewInMemoryDB()

		req, err := createRequest(http.MethodPost, URL, user)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			emptyDBUser,
		)
	})

	t.Run("last name length should be >= 2", func(t *testing.T) {
		user := requestBody
		user.LastName = ""

		db := database.NewInMemoryDB()

		req, err := createRequest(http.MethodPost, URL, user)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			emptyDBUser,
		)
	})

	t.Run("last name length should be <= 20", func(t *testing.T) {
		user := requestBody
		user.LastName = "DoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoe"

		db := database.NewInMemoryDB()

		req, err := createRequest(http.MethodPost, URL, user)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			emptyDBUser,
		)
	})

	t.Run("biography length should be >= 20", func(t *testing.T) {
		user := requestBody
		user.Biography = ""

		db := database.NewInMemoryDB()

		req, err := createRequest(http.MethodPost, URL, user)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			emptyDBUser,
		)
	})

	t.Run("biography length should be <= 450", func(t *testing.T) {
		user := requestBody
		user.Biography = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac eleifend felis, a dictum lacus. Vivamus nibh tellus, lobortis ac luctus vel, hendrerit in sapien. Pellentesque fringilla blandit interdum. Nullam at placerat dolor. Vivamus at hendrerit urna, eget interdum lorem. Curabitur a libero eget erat bibendum imperdiet. Morbi aliquet tellus id egestas vehicula.Curabitur eget elit pellentesque, ullamcorper est ut, vehicula nisi. Duis rhoncus cursus mi a convallis. Vestibulum sit amet vestibulum magna. Suspendisse posuere convallis nisi sed viverra. Sed molestie enim eget dignissim tincidunt. Curabitur eget sollicitudin dolor. In maximus dictum massa, sit amet commodo tellus. Nunc tempor sit amet libero vel tempor. Vestibulum sollicitudin risus sed augue pulvinar malesuada. Proin in tempus dolor, vel varius orci. Aliquam id nibh eu purus viverra vehicula ut."

		db := database.NewInMemoryDB()

		req, err := createRequest(http.MethodPost, URL, user)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			emptyDBUser,
		)
	})
}

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

func assertResponse(
	t testing.TB,
	resp *httptest.ResponseRecorder,
	expectedStatus int,
	expectedMessage string,
	expectedData database.DBUser,
) {
	t.Helper()
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("could not decode the response: %v", err)
	}

	if resp.Code != expectedStatus {
		t.Errorf("expected status %d; got %d", expectedStatus, resp.Code)
	}

	if response.Message != expectedMessage {
		t.Errorf("expected message %q; got %q", expectedMessage, response.Message)
	}

	if !expectedData.IsEmpty() {
		dataBytes, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("could not marshal the data: %v", err)
		}

		var got database.DBUser
		if err := json.Unmarshal(dataBytes, &got); err != nil {
			t.Fatalf("could not unmarshal the user: %v", err)
		}

		if got.ID.IsEmpty() ||
			got.User.FirstName != expectedData.User.FirstName ||
			got.User.LastName != expectedData.User.LastName ||
			got.User.Biography != expectedData.User.Biography {
			t.Errorf("expected user %v; got %v", expectedData, got)
		}
	}
}
