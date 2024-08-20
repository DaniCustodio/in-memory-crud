package database

import (
	"errors"
	"main/api"
	"sync"
)

var ErrUserDoesNotExist = errors.New("user does not exist")

type InMemoryDB struct {
	mu   sync.RWMutex
	data map[api.ID]api.User
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[api.ID]api.User),
	}
}

func (db *InMemoryDB) Insert(id api.ID, value api.User) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[id] = value
}

func (db *InMemoryDB) Update(id api.ID, updatedUser api.User) (api.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.data[id]; !exists {
		return api.User{}, ErrUserDoesNotExist
	}

	db.data[id] = updatedUser

	return updatedUser, nil
}

func (db *InMemoryDB) Delete(id api.ID) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, id)
}

func (db *InMemoryDB) FindAll() []api.User {
	db.mu.RLock()
	defer db.mu.RUnlock()
	users := make([]api.User, 0, len(db.data))
	for _, user := range db.data {
		users = append(users, user)
	}

	return users
}

func (db *InMemoryDB) FindByID(id api.ID) (api.User, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	user, exists := db.data[id]
	return user, exists
}
