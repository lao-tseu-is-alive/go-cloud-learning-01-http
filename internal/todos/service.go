package todos

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Service struct {
	Log   *log.Logger
	Store Storage
}

type ErrorService struct {
	Err    error  `json:"err" `
	Status int    `json:"status" `
	Msg    string `json:"msg" `
}

func (e *ErrorService) Error() string {
	return fmt.Sprintf("Status[%d] %s. error: %v", e.Status, e.Msg, e.Err)
}

// GetMaxId returns the greatest todos id used by now
// curl -H "Content-Type: application/json" 'http://localhost:8080/todos/maxid'
func (s Service) GetMaxId(ctx echo.Context) error {
	s.Log.Println("# Entering GetMaxId()")
	var maxTodoId int32 = 0
	maxTodoId, _ = s.Store.GetMaxId()
	s.Log.Printf("# Exit GetMaxId() maxTodoId: %d", maxTodoId)
	return ctx.JSON(http.StatusOK, maxTodoId)
}

func (s Service) GetTodo(ctx echo.Context, todoId int32) error {
	s.Log.Printf("# Entering GetTodo(%d)", todoId)
	if s.Store.Exist(todoId) == false {
		return ctx.JSON(http.StatusNotFound, ErrorService{
			Err:    errors.New("not found"),
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("todo id : %d does not exist", todoId),
		})
	}
	todo, err := s.Store.Get(todoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem retrieving todo :%v", err))
	}
	return ctx.JSON(http.StatusOK, todo)
}

//GetTodos will retrieve all Todos in the store and return then
//to test it with curl you can try :
//curl -H "Content-Type: application/json" 'http://localhost:8080/todos' |json_pp
func (s Service) GetTodos(ctx echo.Context, params GetTodosParams) error {
	s.Log.Printf("# Entering GetTodos() %v", params)
	list, err := s.Store.List(0, 100)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("there was a problem when calling store.List :%v", err))
	}
	return ctx.JSON(http.StatusOK, list)
}

//CreateTodo will store the NewTodo task in the store
//to test it with curl you can try :
//curl -XPOST -H "Content-Type: application/json" -d '{"task":"learn Linux"}'  'http://localhost:8080/todos'
//curl -XPOST -H "Content-Type: application/json" -d '{"task":""}'  'http://localhost:8080/todos'
func (s Service) CreateTodo(ctx echo.Context) error {
	s.Log.Println("# Entering CreateTodo()")
	newTodo := &NewTodo{}
	if err := ctx.Bind(newTodo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("CreateTodo has invalid format [%v]", err))
	}
	if len(newTodo.Task) < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateTodo task cannot be empty"))
	}
	if len(newTodo.Task) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateTodo task minLength is 5"))
	}
	s.Log.Printf("# CreateTodo() newTodo : %#v\n", newTodo)
	todoCreated, err := s.Store.Create(*newTodo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem saving new todo :%v", err))
	}
	s.Log.Printf("# CreateTodo() Todo %#v\n", todoCreated)
	return ctx.JSON(http.StatusCreated, todoCreated)

}

// UpdateTodo will store the modified information in the store for the given todoId
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": true}'  'http://localhost:8080/todos/3'
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": false}'  'http://localhost:8080/todos/3'
func (s Service) UpdateTodo(ctx echo.Context, todoId int32) error {
	s.Log.Printf("# Entering UpdateTodo(%d)", todoId)
	if s.Store.Exist(todoId) == false {
		return ctx.JSON(http.StatusNotFound, ErrorService{
			Err:    errors.New("not found"),
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("todo id : %d does not exist", todoId),
		})
	}
	t := new(Todo)
	if err := ctx.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("UpdateTodo has invalid format [%v]", err))
	}
	if len(t.Task) < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateTodo task cannot be empty"))
	}
	//refuse an attempt to modify a todoId (in url) with a different id in the body !
	if t.Id != todoId {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("UpdateTodo id : [%d] and posted Id [%d] cannot differ ", todoId, t.Id))
	}

	updatedTodo, err := s.Store.Update(todoId, *t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem updating todo :%v", err))
	}
	return ctx.JSON(http.StatusOK, updatedTodo)
}

// DeleteTodo will remove the given todoID entry from the store, and if not present will return 400 Bad Request
//curl -v -XDELETE -H "Content-Type: application/json" 'http://localhost:8080/todos/3' ->  204 No Content if present and delete it
//curl -v -XDELETE -H "Content-Type: application/json" 'http://localhost:8080/todos/93333' -> 400 Bad Request
func (s Service) DeleteTodo(ctx echo.Context, todoId int32) error {
	s.Log.Printf("# Entering DeleteTodo(%d)", todoId)
	if s.Store.Exist(todoId) == false {
		return ctx.JSON(http.StatusNotFound, ErrorService{
			Err:    errors.New("not found"),
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("todo id : %d does not exist", todoId),
		})
	} else {
		err := s.Store.Delete(todoId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem deleting todo :%v", err))
		}
		return ctx.NoContent(http.StatusNoContent)
	}
}
