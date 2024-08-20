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
			Biography: "A regular guy who loves to code in Go and JavaScript",
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

		if got.ID == ID(uuid.Nil) {
			t.Error("expected a non-empty ID")
		}
	})

	t.Run("first name length should be >= 2", func(t *testing.T) {
		user := CreateUserBody{
			FirstName: "",
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

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400; got %d", rec.Code)
		}

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if response.Message == "" {
			t.Error("expected an error message")
		}
	})

	t.Run("first name length should be <= 20", func(t *testing.T) {
		user := CreateUserBody{
			FirstName: "JohnJohnJohnJohnJohnJohn",
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

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400; got %d", rec.Code)
		}

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if response.Message == "" {
			t.Error("expected an error message")
		}
	})

	t.Run("last name length should be >= 2", func(t *testing.T) {
		user := CreateUserBody{
			FirstName: "John",
			LastName:  "",
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

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400; got %d", rec.Code)
		}

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if response.Message == "" {
			t.Error("expected an error message")
		}
	})

	t.Run("last name length should be <= 20", func(t *testing.T) {
		user := CreateUserBody{
			FirstName: "John",
			LastName:  "DoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoeDoe",
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

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400; got %d", rec.Code)
		}

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if response.Message == "" {
			t.Error("expected an error message")
		}
	})

	t.Run("biography length should be >= 20", func(t *testing.T) {
		user := CreateUserBody{
			FirstName: "John",
			LastName:  "Doe",
			Biography: "",
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

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400; got %d", rec.Code)
		}

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if response.Message == "" {
			t.Error("expected an error message")
		}
	})

	t.Run("biography length should be <= 450", func(t *testing.T) {
		user := CreateUserBody{
			FirstName: "John",
			LastName:  "Doe",
			Biography: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac eleifend felis, a dictum lacus. Vivamus nibh tellus, lobortis ac luctus vel, hendrerit in sapien. Pellentesque fringilla blandit interdum. Nullam at placerat dolor. Vivamus at hendrerit urna, eget interdum lorem. Curabitur a libero eget erat bibendum imperdiet. Morbi aliquet tellus id egestas vehicula.Curabitur eget elit pellentesque, ullamcorper est ut, vehicula nisi. Duis rhoncus cursus mi a convallis. Vestibulum sit amet vestibulum magna. Suspendisse posuere convallis nisi sed viverra. Sed molestie enim eget dignissim tincidunt. Curabitur eget sollicitudin dolor. In maximus dictum massa, sit amet commodo tellus. Nunc tempor sit amet libero vel tempor. Vestibulum sollicitudin risus sed augue pulvinar malesuada. Proin in tempus dolor, vel varius orci. Aliquam id nibh eu purus viverra vehicula ut.",
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

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400; got %d", rec.Code)
		}

		var response Response
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode the response: %v", err)
		}

		if response.Message == "" {
			t.Error("expected an error message")
		}
	})
}
