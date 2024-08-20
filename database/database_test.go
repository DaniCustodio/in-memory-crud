package database

import (
	"fmt"
	"main/api"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func TestInMemoryDB(t *testing.T) {
	users := []api.User{
		{
			ID:        api.ID(uuid.New()),
			FirstName: "John",
			LastName:  "Doe",
			Biography: "A simple guy who loves to write code and play games. He is a fan of technology and loves to read about new things.",
		},
		{
			ID:        api.ID(uuid.New()),
			FirstName: "Jane",
			LastName:  "Doe",
			Biography: "A nice lady who loves to write code and play games. She is a fan of technology and loves to read about new things.",
		},
	}

	t.Run("save a user and find by id", func(t *testing.T) {
		db := NewInMemoryDB()

		user := users[0]

		db.Insert(user.ID, user)

		got, exists := db.FindByID(user.ID)

		if !exists {
			t.Fatalf("expected the user to exist in the database")
		}

		if got != user {
			t.Fatalf("expected the user to be %v, got %v", user, got)
		}
	})

	t.Run("find all users", func(t *testing.T) {
		db := NewInMemoryDB()

		for _, user := range users {
			db.Insert(user.ID, user)
		}

		got := db.FindAll()

		if len(got) != len(users) {
			t.Fatalf("expected the number of users to be %d, got %d", len(users), len(got))
		}
	})

	t.Run("delete a user", func(t *testing.T) {
		db := NewInMemoryDB()

		user := users[0]

		db.Insert(user.ID, user)

		db.Delete(user.ID)

		_, exists := db.FindByID(user.ID)

		if exists {
			t.Fatalf("expected the user to be deleted from the database")
		}
	})

	t.Run("update a user", func(t *testing.T) {
		db := NewInMemoryDB()

		user := users[0]

		db.Insert(user.ID, user)

		updatedUser := users[1]

		db.Update(user.ID, updatedUser)

		got, exists := db.FindByID(user.ID)

		if !exists {
			t.Fatalf("expected the user to exist in the database")
		}

		if got != updatedUser {
			t.Fatalf("expected the user to be %v, got %v", updatedUser, got)
		}
	})

	t.Run("update a user that does not exist", func(t *testing.T) {
		db := NewInMemoryDB()

		user := users[0]

		updatedUser := users[1]

		_, err := db.Update(user.ID, updatedUser)

		if err != ErrUserDoesNotExist {
			t.Fatalf("expected the error to be %v, got %v", ErrUserDoesNotExist, err)
		}
	})
}

func TestConcurrent(t *testing.T) {
	db := NewInMemoryDB()

	var wg sync.WaitGroup
	numRoutines := 100

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			key := api.ID(uuid.New())
			db.Insert(key, api.User{
				ID:        key,
				FirstName: fmt.Sprintf("John%d", i),
				LastName:  "Doe",
				Biography: "A simple guy who loves to write code and play games. He is a fan of technology and loves to read about new things.",
			})
		}(i)
	}

	wg.Wait()

	users := db.FindAll()

	if len(users) != numRoutines {
		t.Fatalf("expected the number of users to be %d, got %d", numRoutines, len(users))
	}
}
