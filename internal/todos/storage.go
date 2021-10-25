package todos

import (
	"errors"
	"fmt"
	todos "github.com/lao-tseu-is-alive/go-cloud-learning-01-http/gen"
)

type ErrorTodos struct {
	err error
	msg string
}

func (e *ErrorTodos) Error() string {
	return fmt.Sprintf("%s : %v", e.msg, e.err)
}

type Storage interface {
	// List returns the list of existing todos with the given offset and limit.
	List(offset, limit int) ([]todos.Todo, error)
	// Get returns the todos with the specified todos ID.
	Get(id int32) (*todos.Todo, error)
	// Exist returns true only if a todos with the specified id exists in store.
	Exist(id int32) bool
	// Count returns the total number of todos.
	Count() (int, error)
	// Create saves a new todos in the storage.
	Create(todo todos.NewTodo) (*todos.Todo, error)
	// Update updates the todos with given ID in the storage.
	Update(id int32, todo todos.Todo) (*todos.Todo, error)
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
