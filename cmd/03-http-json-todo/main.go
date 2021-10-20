package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultServerPort = 8080
	defaultServerIp   = "127.0.0.1"
	defaultMaxId      = 2
)

type memoryStore struct {
	Todos map[int32]*Todo
	maxId int32
	lock  sync.RWMutex
}

type goTodoServer struct {
	data   *memoryStore
	logger *log.Logger
}

// GetMaxId returns the greatest todos id used by now
// curl -H "Content-Type: application/json" 'http://localhost:8080/todos/maxid'
func (s goTodoServer) GetMaxId(ctx echo.Context) error {
	s.logger.Println("# Entering GetMaxId()")
	var maxTodoId int32 = 0
	s.data.lock.RLock()
	defer s.data.lock.RUnlock()
	for myTodoId, _ := range s.data.Todos {
		if myTodoId > maxTodoId {
			maxTodoId = myTodoId
		}
	}
	s.logger.Printf("# Exit GetMaxId() maxTodoId: %d", maxTodoId)
	return ctx.JSON(http.StatusOK, maxTodoId)
}

//GetTodos will retrieve all Todos in the store and return then
//to test it with curl you can try :
//curl -H "Content-Type: application/json" 'http://localhost:8080/todos' |json_pp

func (s goTodoServer) GetTodos(ctx echo.Context, params GetTodosParams) error {
	s.logger.Println("# Entering GetTodos()")
	s.data.lock.RLock()
	defer s.data.lock.RUnlock()
	return ctx.JSON(http.StatusOK, s.data.Todos)
}

//CreateTodo will store the NewTodo task in the store
//to test it with curl you can try :
//curl -XPOST -H "Content-Type: application/json" -d '{"task":"learn Linux"}'  'http://localhost:8080/todos'
//curl -XPOST -H "Content-Type: application/json" -d '{"task":""}'  'http://localhost:8080/todos'
func (s goTodoServer) CreateTodo(ctx echo.Context) error {
	s.logger.Println("# Entering CreateTodo()")
	s.data.lock.Lock()
	defer s.data.lock.Unlock()
	now := time.Now()
	s.data.maxId++
	t := &Todo{
		Id:        s.data.maxId,
		CreatedAt: &now,
	}
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
	s.logger.Printf("# CreateTodo() newTodo : %#v\n", newTodo)
	t.Task = newTodo.Task
	s.logger.Printf("# CreateTodo() Todo %#v\n", t)
	s.data.Todos[t.Id] = t
	return ctx.JSON(http.StatusCreated, t)

}

// UpdateTodo will store the modified information in the store for the given todoId
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": true}'  'http://localhost:8080/todos/3'
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": false}'  'http://localhost:8080/todos/3'
func (s goTodoServer) UpdateTodo(ctx echo.Context, todoId int32) error {
	s.logger.Printf("# Entering UpdateTodo(%d)", todoId)
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
func (s goTodoServer) DeleteTodo(ctx echo.Context, todoId int32) error {
	s.logger.Printf("# Entering DeleteTodo(%d)", todoId)
	if s.data.Todos[todoId] == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		s.data.lock.Lock()
		defer s.data.lock.Unlock()
		delete(s.data.Todos, todoId)
		return ctx.NoContent(http.StatusNoContent)
	}
}

// initializeStorage initialize some dummy data to get some results back
func initializeStorage() memoryStore {

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

	return memoryStore{
		Todos: defaultInitialData,
		maxId: defaultMaxId,
		lock:  sync.RWMutex{},
	}
}

// getNewServer initialize a new Echo server and returns it
func getNewServer(discardLog bool) *echo.Echo {
	var l *log.Logger
	if discardLog == true {
		l = log.New(ioutil.Discard, "todo-api_", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		l = log.New(os.Stdout, "todo-api_", log.Ldate|log.Ltime|log.Lshortfile)
	}
	DBinMemory := initializeStorage()
	myApi := goTodoServer{logger: l, data: &DBinMemory}
	e := echo.New()
	RegisterHandlers(e, &myApi)
	// add a route for maxId
	e.GET("/todos/maxid", myApi.GetMaxId)
	return e
}

//getListenAddrFromEnv returns a valid TCP/IP listening address string based on the values of ENV SERVERIP:PORT
func getListenAddrFromEnv(defaultIP string, defaultPort int) string {
	srvIP := defaultIP
	srvPort := defaultPort
	var err error
	val, exist := os.LookupEnv("PORT")
	if exist {
		srvPort, err = strconv.Atoi(val)
		if err != nil {
			log.Fatal("ERROR: CONFIG ENV PORT should contain a valid integer value !")
		}
	}
	val, exist = os.LookupEnv("SERVERIP")
	if exist {
		srvIP = val
	}
	return fmt.Sprintf("%s:%d", srvIP, srvPort)
}

// main is the entry point of your Todos service you can execute it with :
// SERVERIP=192.168.50.6 PORT=3333 go run main.go todo_*.go
func main() {
	listenAddress := getListenAddrFromEnv(defaultServerIp, defaultServerPort)
	e := getNewServer(true)
	e.Logger.Fatal(e.Start(listenAddress))
}
