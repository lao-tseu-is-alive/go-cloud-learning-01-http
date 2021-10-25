package todos

import (
	"errors"
	todos "github.com/lao-tseu-is-alive/go-cloud-learning-01-http/gen"
	"sync"
	"time"
)

const defaultMaxId = 2

type memoryStore struct {
	Todos map[int32]*todos.Todo
	maxId int32
	lock  sync.RWMutex
}

//Create will store the new task in the store
func (m memoryStore) Create(todo todos.NewTodo) (*todos.Todo, error) {
	if len(todo.Task) < 1 {
		return nil, errors.New("todo task cannot be empty")
	}
	if len(todo.Task) < 6 {
		return nil, errors.New("CreateTodo task minLength is 5")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	now := time.Now()
	m.maxId++
	t := &todos.Todo{
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   &now,
		Id:          m.maxId,
		Task:        todo.Task,
	}
	m.Todos[t.Id] = t
	return t, nil
}

func (m memoryStore) List(offset, limit int) ([]todos.Todo, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	res := make([]todos.Todo, len(m.Todos))
	//TODO convert map to res slice
	return res, nil
}

func (m memoryStore) Get(id int32) (todos.Todo, error) {
	panic("implement me")
}

// Exist returns true only if a todos with the specified id exists in store.
func (m memoryStore) Exist(id int32) bool {
	if m.Todos[id] == nil {
		return false
	}
	return true
}

func (m memoryStore) Count() (int, error) {
	return len(m.Todos), nil
}

func (m memoryStore) Update(id int32, todo todos.Todo) (*todos.Todo, error) {
	if m.Exist(id) {
		m.lock.Lock()
		defer m.lock.Unlock()
		delete(m.Todos, id)
		existingTodo := m.Todos[id]
		now := time.Now()
		// cannot override CreatedAt fields
		todo.CreatedAt = existingTodo.CreatedAt
		switch todo.Completed {
		case true:
			if existingTodo.Completed == false {
				todo.CompletedAt = &now
			}
		case false:
			if existingTodo.Completed == true {
				// task was completed, but user changed it to not completed
				todo.CompletedAt = nil
			}
		// in all other cases the value of CompletedAt should not be changed
		default:
			todo.CompletedAt = existingTodo.CompletedAt
		}

		m.Todos[id] = &todo
		return &todo, nil
	}
	return nil, errors.New("todo with this id does not exist")
}

func (m memoryStore) Delete(id int32) error {
	if m.Exist(id) {
		m.lock.Lock()
		defer m.lock.Unlock()
		delete(m.Todos, id)
		return nil
	}
	return errors.New("todo with this id does not exist")
}

// initializeStorage initialize some dummy data to get some results back
func initializeStorage() memoryStore {

	someTimeCreated, _ := time.Parse(time.RFC3339, "2020-02-21T08:00:23.877Z")
	someTimeCompleted, _ := time.Parse(time.RFC3339, "2021-10-07T15:02:23.877Z")
	defaultInitialData := map[int32]*todos.Todo{
		1: {
			Completed:   true,
			CompletedAt: &someTimeCompleted,
			CreatedAt:   &someTimeCreated,
			Id:          1,
			Task:        "Learn GO",
		},
		2: {
			Completed:   false,
			CompletedAt: nil,
			CreatedAt:   &someTimeCreated,
			Id:          2,
			Task:        "Learn OpenAPI",
		},
	}

	return memoryStore{
		Todos: defaultInitialData,
		maxId: defaultMaxId,
		lock:  sync.RWMutex{},
	}
}

func NewMemoryDB() (Storage, error) {
	DBinMemory := initializeStorage()
	return DBinMemory, nil
}
