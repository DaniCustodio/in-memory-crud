package database

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrUserDoesNotExist = errors.New("user does not exist")

type ID uuid.UUID

type User struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=20"`
	LastName  string `json:"last_name" validate:"required,min=2,max=20"`
	Biography string `json:"biography" validate:"required,min=20,max=450"`
}

type DBUser struct {
	ID   ID   `json:"id"`
	User User `json:"user"`
}

func (d DBUser) IsEmpty() bool {
	return d.User == User{} && d.ID == ID(uuid.Nil)
}

type InMemoryDB struct {
	mu   sync.RWMutex
	data map[ID]User
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[ID]User),
	}
}

func (db *InMemoryDB) Insert(value User) DBUser {
	db.mu.Lock()
	defer db.mu.Unlock()

	id := ID(uuid.New())
	db.data[id] = value

	return DBUser{
		ID:   id,
		User: value,
	}
}

func (db *InMemoryDB) Update(id ID, updatedUser User) (DBUser, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.data[id]; !exists {
		return DBUser{}, ErrUserDoesNotExist
	}

	db.data[id] = updatedUser

	return DBUser{
		ID:   id,
		User: updatedUser,
	}, nil
}

func (db *InMemoryDB) Delete(id ID) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, id)
}

func (db *InMemoryDB) FindAll() []User {
	db.mu.RLock()
	defer db.mu.RUnlock()
	users := make([]User, 0, len(db.data))
	for _, user := range db.data {
		users = append(users, user)
	}

	return users
}

func (db *InMemoryDB) FindByID(id ID) (User, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	user, exists := db.data[id]
	return user, exists
}
