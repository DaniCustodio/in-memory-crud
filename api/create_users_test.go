package api

import (
	"main/database"
	"net/http"
	"testing"
)

func TestCreateUser(t *testing.T) {
	const URL = "/api/users"

	requestBody := database.User{
		FirstName: "John",
		LastName:  "Doe",
		Biography: "A regular guy who loves to code in Go and JavaScript",
	}

	t.Run("create a user successfully", func(t *testing.T) {
		db := database.NewInMemoryDB()

		req, err := createRequest(http.MethodPost, URL, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		rec := makeRequest(db, req)

		response, err := parseResponse[database.DBUser](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusCreated, rec.Code)

		assertUser(
			t,
			database.DBUser{
				ID:   database.ID{}.NewID(),
				User: requestBody,
			},
			response.Data,
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

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusBadRequest, rec.Code)

		assertErrorMessage(t, "Please provide a valid FirstName, LastName and Bio for the user", response.Message)
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

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusBadRequest, rec.Code)

		assertErrorMessage(t, "Please provide a valid FirstName, LastName and Bio for the user", response.Message)
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

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusBadRequest, rec.Code)

		assertErrorMessage(t, "Please provide a valid FirstName, LastName and Bio for the user", response.Message)
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

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusBadRequest, rec.Code)

		assertErrorMessage(t, "Please provide a valid FirstName, LastName and Bio for the user", response.Message)
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

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusBadRequest, rec.Code)

		assertErrorMessage(t, "Please provide a valid FirstName, LastName and Bio for the user", response.Message)
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

		response, err := parseResponse[any](rec)
		if err != nil {
			t.Fatalf("could not parse the response: %v", err)
		}

		assertStatusCode(t, http.StatusBadRequest, rec.Code)

		assertErrorMessage(t, "Please provide a valid FirstName, LastName and Bio for the user", response.Message)
	})
}
