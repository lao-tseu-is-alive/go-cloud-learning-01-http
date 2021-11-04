package todos

import (
	"errors"
	"sort"
	"sync"
	"time"
)

const DefaultMaxId = 2

type memoryStore struct {
	Todos map[int32]*Todo
	maxId int32
	lock  sync.RWMutex
}

//Create will store the new task in the store
func (m *memoryStore) Create(todo NewTodo) (*Todo, error) {
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
	t := &Todo{
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   &now,
		Id:          m.maxId,
		Task:        todo.Task,
	}
	m.Todos[t.Id] = t
	return t, nil
}

func (m *memoryStore) List(offset, limit int) ([]*Todo, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var res []*Todo
	if offset > 0 && offset < len(m.Todos) {
		// handle offset
	} else {
		if limit > len(m.Todos) {
			//return all
			keys := make([]int, 0, len(m.Todos))
			for k := range m.Todos {
				keys = append(keys, int(k))
			}
			sort.Ints(keys)

			for _, k := range keys {
				res = append(res, m.Todos[int32(k)])
			}
		}
	}
	return res, nil
}

func (m *memoryStore) Get(id int32) (*Todo, error) {
	if m.Exist(id) {
		m.lock.RLock()
		defer m.lock.RUnlock()
		existingTodo := m.Todos[id]
		return existingTodo, nil
	}
	return nil, errors.New("todo with this id does not exist")
}

// GetMaxId returns the maximum value of todos id existing in store.
func (m *memoryStore) GetMaxId() (int32, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	existingMaxId := int32(0)
	for _, t := range m.Todos {
		if t.Id > existingMaxId {
			existingMaxId = t.Id
		}
	}
	return existingMaxId, nil
}

// Exist returns true only if a todos with the specified id exists in store.
func (m *memoryStore) Exist(id int32) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.Todos[id] == nil {
		return false
	}
	return true
}

func (m *memoryStore) Count() (int32, error) {
	return int32(len(m.Todos)), nil
}

func (m *memoryStore) Update(id int32, todo Todo) (*Todo, error) {
	if m.Exist(id) {
		m.lock.Lock()
		defer m.lock.Unlock()
		existingTodo := m.Todos[id]
		now := time.Now()
		// cannot override id field
		todo.Id = existingTodo.Id
		// cannot override CreatedAt field
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

func (m *memoryStore) Delete(id int32) error {
	if m.Exist(id) {
		m.lock.Lock()
		defer m.lock.Unlock()
		delete(m.Todos, id)
		return nil
	}
	return errors.New("todo with this id does not exist")
}

// Close : will do cleanup for all todos stored in memory
func (m *memoryStore) Close() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for idx, _ := range m.Todos {
		delete(m.Todos, idx)
	}
	return
}

// initializeStorage initialize some dummy data to get some results back
func initializeStorage() *memoryStore {

	someTimeCreated, _ := time.Parse(time.RFC3339, "2020-02-21T08:00:23.877Z")
	someTimeCompleted, _ := time.Parse(time.RFC3339, "2021-10-07T15:02:23.877Z")
	defaultInitialData := map[int32]*Todo{
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

	return &memoryStore{
		Todos: defaultInitialData,
		maxId: DefaultMaxId,
		lock:  sync.RWMutex{},
	}
}

func NewMemoryDB() (Storage, error) {
	DBinMemory := initializeStorage()
	return DBinMemory, nil
}
