package database

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

var ErrUserDoesNotExist = errors.New("user does not exist")

type ID uuid.UUID

func (i ID) NewID() ID {
	return ID(uuid.New())
}

func (i ID) String() string {
	return uuid.UUID(i).String()
}

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

func (db *InMemoryDB) Update(id string, updatedUser User) (DBUser, error) {
	parsedID, err := parseID(id)
	if err != nil {
		return DBUser{}, err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.data[parsedID]; !exists {
		return DBUser{}, ErrUserDoesNotExist
	}

	db.data[parsedID] = updatedUser

	return DBUser{
		ID:   parsedID,
		User: updatedUser,
	}, nil
}

func (db *InMemoryDB) Delete(id string) (DBUser, error) {
	parsedID, err := parseID(id)
	if err != nil {
		return DBUser{}, err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	user, exists := db.data[parsedID]
	if !exists {
		return DBUser{}, ErrUserDoesNotExist
	}

	delete(db.data, parsedID)

	return DBUser{ID: parsedID, User: user}, nil
}

func (db *InMemoryDB) FindAll() []DBUser {
	db.mu.RLock()
	defer db.mu.RUnlock()

	users := make([]DBUser, 0, len(db.data))
	for id, user := range db.data {
		users = append(users, DBUser{
			ID:   id,
			User: user,
		})
	}

	return users
}

func (db *InMemoryDB) FindByID(id string) (DBUser, bool) {
	parsedID, err := parseID(id)
	if err != nil {
		return DBUser{}, false
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	user, exists := db.data[parsedID]
	return DBUser{ID: parsedID, User: user}, exists
}

func parseID(id string) (ID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("could not parse the id", "error", err)
		return ID{}, err
	}

	return ID(parsedID), nil
}
