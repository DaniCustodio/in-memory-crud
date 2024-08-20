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
	requestBody := CreateUserBody{
		FirstName: "John",
		LastName:  "Doe",
		Biography: "A regular guy who loves to code in Go and JavaScript",
	}

	t.Run("create a user successfully", func(t *testing.T) {
		rec := makeRequest(t, requestBody)

		assertResponse(
			t,
			rec,
			http.StatusCreated,
			"",
			User{
				ID:        ID(uuid.Nil),
				FirstName: "John",
				LastName:  "Doe",
				Biography: "A regular guy who loves to code in Go and JavaScript",
			},
		)
	})

	t.Run("first name length should be >= 2", func(t *testing.T) {
		user := requestBody
		user.FirstName = ""

		rec := makeRequest(t, user)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			User{},
		)
	})

	t.Run("first name length should be <= 20", func(t *testing.T) {
		user := requestBody
		user.FirstName = "JohnDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoe"

		rec := makeRequest(t, user)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			User{},
		)
	})

	t.Run("last name length should be >= 2", func(t *testing.T) {
		user := requestBody
		user.LastName = ""

		rec := makeRequest(t, user)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			User{},
		)
	})

	t.Run("last name length should be <= 20", func(t *testing.T) {
		user := requestBody
		user.LastName = "DoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoe"

		rec := makeRequest(t, user)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			User{},
		)
	})

	t.Run("biography length should be >= 20", func(t *testing.T) {
		user := requestBody
		user.Biography = ""

		rec := makeRequest(t, user)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			User{},
		)
	})

	t.Run("biography length should be <= 450", func(t *testing.T) {
		user := requestBody
		user.Biography = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac eleifend felis, a dictum lacus. Vivamus nibh tellus, lobortis ac luctus vel, hendrerit in sapien. Pellentesque fringilla blandit interdum. Nullam at placerat dolor. Vivamus at hendrerit urna, eget interdum lorem. Curabitur a libero eget erat bibendum imperdiet. Morbi aliquet tellus id egestas vehicula.Curabitur eget elit pellentesque, ullamcorper est ut, vehicula nisi. Duis rhoncus cursus mi a convallis. Vestibulum sit amet vestibulum magna. Suspendisse posuere convallis nisi sed viverra. Sed molestie enim eget dignissim tincidunt. Curabitur eget sollicitudin dolor. In maximus dictum massa, sit amet commodo tellus. Nunc tempor sit amet libero vel tempor. Vestibulum sollicitudin risus sed augue pulvinar malesuada. Proin in tempus dolor, vel varius orci. Aliquam id nibh eu purus viverra vehicula ut."

		rec := makeRequest(t, user)

		assertResponse(
			t,
			rec,
			http.StatusBadRequest,
			"Please provide a valid FirstName, LastName and Bio for the user",
			User{},
		)
	})
}

func makeRequest(t testing.TB, user CreateUserBody) *httptest.ResponseRecorder {
	t.Helper()
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

	return rec
}

func assertResponse(
	t testing.TB,
	resp *httptest.ResponseRecorder,
	expectedStatus int,
	expectedMessage string,
	expectedData User,
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

	if !expectedData.isEmpty() {
		dataBytes, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("could not marshal the data: %v", err)
		}

		var got User
		if err := json.Unmarshal(dataBytes, &got); err != nil {
			t.Fatalf("could not unmarshal the user: %v", err)
		}

		if got.ID == ID(uuid.Nil) ||
			got.FirstName != expectedData.FirstName ||
			got.LastName != expectedData.LastName ||
			got.Biography != expectedData.Biography {
			t.Errorf("expected user %v; got %v", expectedData, got)
		}
	}
}
