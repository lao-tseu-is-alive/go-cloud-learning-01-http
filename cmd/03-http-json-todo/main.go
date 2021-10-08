package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type memoryStore struct {
	Todos map[int32]*Todo
	IdSeq int32
	Lock  sync.Mutex
}

var DBinMemory memoryStore

type goTodoServer struct {
	data   *memoryStore
	logger *log.Logger
}

func (s goTodoServer) GetTodos(ctx echo.Context, params GetTodosParams) error {
	s.logger.Println("# Entering GetTodos()")
	return ctx.JSON(http.StatusOK, s.data.Todos)
}

//CreateTodo will store the NewTodo task in the store
//to test it with curl you can try :
//curl -XPOST -H "Content-Type: application/json" -d '{"task":"learn Linux2"}'  'http://localhost:8080/todos'
func (s goTodoServer) CreateTodo(ctx echo.Context) error {
	s.logger.Println("# Entering CreateTodo()")
	now := time.Now()
	t := &Todo{
		Id:        s.data.IdSeq,
		CreatedAt: &now,
	}
	newTodo := &NewTodo{Task: ""}
	if err := ctx.Bind(newTodo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("CreateTodo has invalid format [%v]", err))
	}
	if len(newTodo.Task) < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateTodo task cannot be empty"))
	}
	if len(newTodo.Task) < 6 {
		return echo.NewHTTPError(http.StatusOK, fmt.Sprint("CreateTodo task minLength is 5"))
	}
	s.logger.Printf("# CreateTodo() newTodo : %#v\n", newTodo)
	t.Task = newTodo.Task
	s.logger.Printf("# CreateTodo() Todo %#v\n", t)
	s.data.Todos[t.Id] = t
	s.data.IdSeq++
	return ctx.JSON(http.StatusCreated, t)

}

// DeleteTodo will remove the given todoID entry from the store, and if not present will return 400 Bad Request
//curl -v -XDELETE -H "Content-Type: application/json" 'http://localhost:8080/todos/3' ->  204 No Content if present and delete it
//curl -v -XDELETE -H "Content-Type: application/json" 'http://localhost:8080/todos/93333' -> 400 Bad Request
func (s goTodoServer) DeleteTodo(ctx echo.Context, todoId int32) error {
	s.logger.Println("# Entering DeleteTodo()")
	if s.data.Todos[todoId] == nil {
		return ctx.NoContent(http.StatusBadRequest)
	} else {
		delete(s.data.Todos, todoId)
		return ctx.NoContent(http.StatusNoContent)
	}
}

// UpdateTodo will store the modified information in the store for the given todoId
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": true}'  'http://localhost:8080/todos/3'
func (s goTodoServer) UpdateTodo(ctx echo.Context, todoId int32) error {
	s.logger.Println("# Entering UpdateTodo()")
	if s.data.Todos[todoId] == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	t := new(Todo)
	if err := ctx.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("UpdateTodo has invalid format [%v]", err))
	}
	s.data.Todos[todoId] = t
	return ctx.JSON(http.StatusOK, t)
}

// initializeStorage initialize some dummy data to get some results back
func initializeStorage() memoryStore {

	someTimeCreated, _ := time.Parse(time.RFC3339, "2020-02-21T08:00:23.877Z")
	someTimeCompleted, _ := time.Parse(time.RFC3339, "2021-10-07T15:02:23.877Z")
	return memoryStore{
		Todos: map[int32]*Todo{
			1: {
				Completed:   true,
				CompletedAt: &someTimeCompleted,
				CreatedAt:   &someTimeCreated,
				Id:          1,
				Task:        "Learn GO ",
			},
			2: {
				Completed:   false,
				CompletedAt: nil,
				CreatedAt:   &someTimeCreated,
				Id:          2,
				Task:        "Learn OpenAPI ",
			},
		},
		IdSeq: 3,
		Lock:  sync.Mutex{},
	}
}

func main() {
	l := log.New(os.Stdout, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	DBinMemory = initializeStorage()
	myApi := goTodoServer{logger: l, data: &DBinMemory}
	e := echo.New()
	RegisterHandlers(e, &myApi)

	e.Logger.Fatal(e.Start(":8080"))
}
