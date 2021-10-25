package todos

import (
	"fmt"
	"github.com/labstack/echo/v4"
	todoGen "github.com/lao-tseu-is-alive/go-cloud-learning-01-http/gen"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	*echo.Echo
	log   *log.Logger
	store Storage
}

// GetMaxId returns the greatest todos id used by now
// curl -H "Content-Type: application/json" 'http://localhost:8080/todos/maxid'
func (s Server) GetMaxId(ctx echo.Context) error {
	s.log.Println("# Entering GetMaxId()")
	var maxTodoId int32 = 0
	maxTodoId, _ = s.store.Count()
	/*defer s.data.lock.RUnlock()
	for myTodoId, _ := range s.data.Todos {
		if myTodoId > maxTodoId {
			maxTodoId = myTodoId
		}
	}*/
	s.log.Printf("# Exit GetMaxId() maxTodoId: %d", maxTodoId)
	return ctx.JSON(http.StatusOK, maxTodoId)
}

//GetTodos will retrieve all Todos in the store and return then
//to test it with curl you can try :
//curl -H "Content-Type: application/json" 'http://localhost:8080/todos' |json_pp

func (s Server) GetTodos(ctx echo.Context, params GetTodosParams) error {
	s.log.Printf("# Entering GetTodos() %v", params)
	list, err := s.store.List(0, 100)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, list)
}

//CreateTodo will store the NewTodo task in the store
//to test it with curl you can try :
//curl -XPOST -H "Content-Type: application/json" -d '{"task":"learn Linux"}'  'http://localhost:8080/todos'
//curl -XPOST -H "Content-Type: application/json" -d '{"task":""}'  'http://localhost:8080/todos'
func (s Server) CreateTodo(ctx echo.Context) error {
	s.log.Println("# Entering CreateTodo()")
	newTodo := &todoGen.NewTodo{}
	if err := ctx.Bind(newTodo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("CreateTodo has invalid format [%v]", err))
	}
	if len(newTodo.Task) < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateTodo task cannot be empty"))
	}
	if len(newTodo.Task) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateTodo task minLength is 5"))
	}
	s.log.Printf("# CreateTodo() newTodo : %#v\n", newTodo)
	t.Task = newTodo.Task
	s.log.Printf("# CreateTodo() Todo %#v\n", t)
	s.data.Todos[t.Id] = t
	return ctx.JSON(http.StatusCreated, t)

}

// UpdateTodo will store the modified information in the store for the given todoId
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": true}'  'http://localhost:8080/todos/3'
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": false}'  'http://localhost:8080/todos/3'
func (s Server) UpdateTodo(ctx echo.Context, todoId int32) error {
	s.log.Printf("# Entering UpdateTodo(%d)", todoId)
	if s.data.Todos[todoId] == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	s.data.lock.Lock()
	defer s.data.lock.Unlock()
	t := new(Todo)
	if err := ctx.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("UpdateTodo has invalid format [%v]", err))
	}
	if len(t.Task) < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateTodo task cannot be empty"))
	}
	if len(t.Task) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateTodo task minLength is 5"))
	}
	existingTodo := s.data.Todos[todoId]
	now := time.Now()
	// cannot override CreatedAt fields
	t.CreatedAt = existingTodo.CreatedAt
	switch t.Completed {
	case true:
		if existingTodo.Completed == false {
			t.CompletedAt = &now
		}
	case false:
		if existingTodo.Completed == true {
			// task was completed, but user changed it to not completed
			t.CompletedAt = nil
		}
	// in all other cases the value of CompletedAt should not be changed
	default:
		t.CompletedAt = existingTodo.CompletedAt
	}

	s.data.Todos[todoId] = t
	return ctx.JSON(http.StatusOK, t)
}

// DeleteTodo will remove the given todoID entry from the store, and if not present will return 400 Bad Request
//curl -v -XDELETE -H "Content-Type: application/json" 'http://localhost:8080/todos/3' ->  204 No Content if present and delete it
//curl -v -XDELETE -H "Content-Type: application/json" 'http://localhost:8080/todos/93333' -> 400 Bad Request
func (s Server) DeleteTodo(ctx echo.Context, todoId int32) error {
	s.log.Printf("# Entering DeleteTodo(%d)", todoId)
	if s.data.Todos[todoId] == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		s.data.lock.Lock()
		defer s.data.lock.Unlock()
		delete(s.data.Todos, todoId)
		return ctx.NoContent(http.StatusNoContent)
	}
}

// GetNewServer initialize a new Echo server and returns it
func GetNewServer(discardLog bool) *echo.Echo {
	var l *log.Logger
	if discardLog == true {
		l = log.New(ioutil.Discard, "todo-api_", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		l = log.New(os.Stdout, "todo-api_", log.Ldate|log.Ltime|log.Lshortfile)
	}
	e := echo.New()
	s, _ := GetInstance("memory", "")
	myApi := Server{
		Echo:  e,
		log:   l,
		store: s,
	}

	todoGen.RegisterHandlers(e, &myApi)
	// add a route for maxId
	e.GET("/todos/maxid", myApi.GetMaxId)
	return e
}
