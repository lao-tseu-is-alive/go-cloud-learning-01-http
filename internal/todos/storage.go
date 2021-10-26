package todos

import (
	"errors"
	"fmt"
)

// Storage is an interface to different implementation of persistence for Todos
type Storage interface {
	// List returns the list of existing todos with the given offset and limit.
	List(offset, limit int) ([]Todo, error)
	// Get returns the todos with the specified todos ID.
	Get(id int32) (*Todo, error)
	// GetMaxId returns the maximum value of todos id existing in store.
	GetMaxId() (int32, error)
	// Exist returns true only if a todos with the specified id exists in store.
	Exist(id int32) bool
	// Count returns the total number of todos.
	Count() (int32, error)
	// Create saves a new todos in the storage.
	Create(todo NewTodo) (*Todo, error)
	// Update updates the todos with given ID in the storage.
	Update(id int32, todo Todo) (*Todo, error)
	// Delete removes the todos with given ID from the storage.
	Delete(id int32) error
}

func GetInstance(dbDriver, dbConnectionString string) (Storage, error) {
	var db Storage
	var err error
	switch dbDriver {
	/*case "pgx":
	db, err = NewPgxDB(dbConnectionString, runtime.NumCPU())
	if err != nil {
		return nil, fmt.Errorf("error opening postgresql database with pgx driver: %s", err)
	}*/
	case "memory":
		db, err = NewMemoryDB()
		if err != nil {
			return nil, fmt.Errorf("error opening memory store: %s", err)
		}
	default:
		return nil, errors.New("unsupported DB driver type")

	}
	return db, nil
}
