package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-learning-01-http/internal/todos"
	"github.com/lao-tseu-is-alive/go-cloud-learning-01-http/pkg/config"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	defaultNewTask = "Learn Linux"
	DEBUG          = false
)

type idCounter struct {
	currentMaxId int32
}

func (c *idCounter) increment() int32 {
	c.currentMaxId = c.currentMaxId + 1
	if DEBUG {
		fmt.Printf("# IN increment() currentMaxId: %d\n", c.currentMaxId)
	}
	return c.currentMaxId
}

func (c *idCounter) decrement() int32 {
	c.currentMaxId = c.currentMaxId - 1
	if DEBUG {
		fmt.Printf("# IN decrement() currentMaxId: %d\n", c.currentMaxId)
	}
	return c.currentMaxId
}

func (c *idCounter) current() int32 {
	if DEBUG {
		fmt.Printf("# IN current() currentMaxId: %d\n", c.currentMaxId)
	}
	return c.currentMaxId
}

func (c *idCounter) currentAsString() string {
	if DEBUG {
		fmt.Printf("# IN currentAsString() currentMaxId: %d\n", c.currentMaxId)
	}
	return fmt.Sprintf("%d", c.currentMaxId)
}

func getUrlForId(myIdCounter idCounter) string {
	return fmt.Sprintf("/todos/%d", myIdCounter.current())
}

type testScenario struct {
	name           string
	wantStatusCode int
	wantBody       string
	r              *http.Request
}

func getTestTable(t *testing.T, ts *httptest.Server, myId idCounter, jsonInitialData, firstTodo []byte) []testScenario {
	newRequest := func(method, url string, body string) *http.Request {
		r, err := http.NewRequest(method, ts.URL+url, strings.NewReader(body))
		if err != nil {
			t.Fatalf("### ERROR http.NewRequest %s on [%s] \n", method, url)
		}
		return r
	}

	return []testScenario{
		{
			name:           "1: GetTodos , should return all the existing Todos as json",
			wantStatusCode: http.StatusOK,
			wantBody:       string(jsonInitialData),
			r:              newRequest(http.MethodGet, "/todos", ""),
		},
		{
			name:           "2: CreateTodo with a valid new Todo task, should return a valid Todo",
			wantStatusCode: http.StatusCreated,
			wantBody:       `{"completed":false,"id":` + fmt.Sprintf("%d", myId.increment()) + ` ,"task":"` + defaultNewTask + `"}`,
			r:              newRequest(http.MethodPost, "/todos", `{"task":"`+defaultNewTask+`"}`),
		},
		{
			name:           "3: After a successful GetTodos maxid should be one increment higher",
			wantStatusCode: http.StatusOK,
			wantBody:       myId.currentAsString(),
			r:              newRequest(http.MethodGet, "/todos/maxid", ""),
		},
		{
			name:           "4: CreateTodo with another valid new Todo task, should return a valid Todo",
			wantStatusCode: http.StatusCreated,
			wantBody:       `{"completed":false,"id":` + fmt.Sprintf("%d", myId.increment()) + ` ,"task":"` + defaultNewTask + `"}`,
			r:              newRequest(http.MethodPost, "/todos", `{"task":"`+defaultNewTask+`"}`),
		},
		{
			name:           "5: CreateTodo with task field of wrong type in body, should return Bad request)",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo has invalid format",
			r:              newRequest(http.MethodPost, "/todos", `{"task":123}`),
		},
		{
			name:           "6: CreateTodo with an empty body, should return Bad request",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task cannot be empty\"}\n",
			r:              newRequest(http.MethodPost, "/todos", `{}`),
		},
		{
			name:           "7: CreateTodo with a missing task field in body, should return Bad request",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task cannot be empty\"}",
			r:              newRequest(http.MethodPost, "/todos", `{"nope":"should fail"}`),
		},
		{
			name:           "8: CreateTodo with an empty task, should return Bad request",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task cannot be empty\"}",
			r:              newRequest(http.MethodPost, "/todos", `{"task":""}`),
		},
		{
			name:           "9: CreateTodo with a task too short(<6), should return Bad request",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task minLength is 5\"}",
			r:              newRequest(http.MethodPost, "/todos", `{"task":"123"}`),
		},
		{
			name:           "10: DeleteTodo with an id that does not exist, should return Not Found",
			wantStatusCode: http.StatusNotFound,
			wantBody:       "",
			r:              newRequest(http.MethodDelete, "/todos/123456789", ""),
		},
		{
			name:           "11: DeleteTodo with an existing id, should return No Content",
			wantStatusCode: http.StatusNoContent,
			wantBody:       "",
			r:              newRequest(http.MethodDelete, getUrlForId(myId), ""),
		},
		{
			name:           "12: DeleteTodo with an id that was just deleted, should return Not Found",
			wantStatusCode: http.StatusNotFound,
			wantBody:       "",
			r:              newRequest(http.MethodDelete, getUrlForId(myId), ""),
		},
		{
			name:           "13: GetMaxId after a successful DeleteTodo, should return one less",
			wantStatusCode: http.StatusOK,
			wantBody:       fmt.Sprintf("%d", myId.decrement()),
			r:              newRequest(http.MethodGet, "/todos/maxid", ""),
		},
		{
			name:           "14: UpdateTodo with an id that does not exist, should return Not Found",
			wantStatusCode: http.StatusNotFound,
			wantBody:       "",
			r:              newRequest(http.MethodPut, "/todos/123456789", ""),
		},
		{
			name:           "15: UpdateTodo with completed=true, should return a Todo updated with completed=true",
			wantStatusCode: http.StatusOK,
			wantBody:       `"completed":true`,
			r:              newRequest(http.MethodPut, getUrlForId(myId), `{"completed":true,"id":`+myId.currentAsString()+` ,"task":"`+defaultNewTask+`"}`),
		},
		{
			name:           "16: UpdateTodo with completed=false, should return a Todo updated with completed=false",
			wantStatusCode: http.StatusOK,
			wantBody:       `"completed":false`,
			r:              newRequest(http.MethodPut, getUrlForId(myId), `{"completed":false,"id":`+myId.currentAsString()+` ,"task":"`+defaultNewTask+`"}`),
		},
		{
			name:           "17: UpdateTodo with empty task, will return a Bad request",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task cannot be empty\"}",
			r:              newRequest(http.MethodPut, getUrlForId(myId), `{"completed":false,"id":`+myId.currentAsString()+` ,"task":""}`),
		},
		{
			name:           "18: UpdateTodo with task id different form id in body, will return a Bad request",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task cannot be empty\"}",
			r:              newRequest(http.MethodPut, getUrlForId(myId), `{"completed":false,"id": 1 ,"task":""}`),
		},
		{
			name:           "19: GetTodo 1 , should return the first Todo as json",
			wantStatusCode: http.StatusOK,
			wantBody:       string(firstTodo),
			r:              newRequest(http.MethodGet, "/todos/1", ""),
		},
		{
			name:           "20: GetTodo 99 , should return the first Todo as json",
			wantStatusCode: http.StatusNotFound,
			wantBody:       "todo id : 99 does not exist",
			r:              newRequest(http.MethodGet, "/todos/99", ""),
		},
		{
			name:           "99:  invalid path, should return 404 not found",
			wantStatusCode: http.StatusNotFound,
			wantBody:       "{\"message\":\"Not Found\"}",
			r:              newRequest(http.MethodGet, "/nothing_available_here", `{"task":"123"}`),
		},
	}
}

func Test_goTodoServer_TodosMemory(t *testing.T) {
	// Create server using the router initialized elsewhere. The router
	// can be a net/http ServeMux a http.DefaultServeMux or
	// any value that satisfies the net/http Handler interface.
	l := log.New(ioutil.Discard, appName, 0)
	InitialDB, _ := todos.GetStorageInstance("memory", "", l)
	myServer := GetNewServer(true, l, InitialDB)
	ts := httptest.NewServer(myServer)
	defer ts.Close()

	res, _ := InitialDB.List(0, 100)
	jsonInitialData, _ := json.Marshal(res)
	firstTodo, _ := json.Marshal(res[0])

	tests := getTestTable(t, ts, idCounter{currentMaxId: todos.DefaultMaxId}, jsonInitialData, firstTodo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			resp, err := http.DefaultClient.Do(tt.r)
			if DEBUG {
				fmt.Printf("### %s : %s on %s\n", tt.name, tt.r.Method, tt.r.URL)
			}
			defer resp.Body.Close()
			if err != nil {
				fmt.Printf("### GOT ERROR : %s\n%s", err, resp.Body)
				t.Fatal(err)
			}
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode, "expected status code should be returned")
			receivedJson, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusCreated {
				// In case of Created we need to modify the CreatedAt value from the response for equality test to pass
				myNewTodo := todos.Todo{
					Completed:   false,
					CompletedAt: nil,
					CreatedAt:   nil,
					Id:          0,
					Task:        defaultNewTask,
				}
				var createdTodo todos.Todo
				err := json.Unmarshal(receivedJson, &createdTodo)
				if err != nil {
					fmt.Printf("FATAL ERROR doing json.Unmarshall of %#v   - Error: %s", string(receivedJson), err)
				}
				myNewTodo.CreatedAt = createdTodo.CreatedAt
				myNewTodo.Id = createdTodo.Id
				wantedTodoJson, _ := json.Marshal(myNewTodo)
				if DEBUG {
					fmt.Printf("WANTED   :%T - %#v\n", wantedTodoJson, string(wantedTodoJson))
					fmt.Printf("RECEIVED :%T - %#v\n", receivedJson, string(receivedJson))
				}
				assert.JSONEqf(t, string(wantedTodoJson), string(receivedJson), "CreateTodo Response was not equal to expected.")
			} else {
				// here are all other cases (except Created above)
				if DEBUG {
					fmt.Printf("WANTED   :%T - %#v\n", tt.wantBody, tt.wantBody)
					fmt.Printf("RECEIVED :%T - %#v\n", receivedJson, string(receivedJson))
				}
				// check that receivedJson contains the specified tt.wantBody substring . https://pkg.go.dev/github.com/stretchr/testify/assert#Contains
				assert.Contains(t, string(receivedJson), tt.wantBody, "Response should contain what was expected.")
			}
		})
	}
}

func Test_goTodoServer_TodosPostgres(t *testing.T) {
	// Create server using the router initialized elsewhere. The router
	// can be a net/http ServeMux a http.DefaultServeMux or
	// any value that satisfies the net/http Handler interface.
	l := log.New(ioutil.Discard, appName, 0)
	dbDsn, err := config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
		appName, appName, defaultDBPassword, defaultDBSslMode)
	if err != nil {
		panic(fmt.Sprintf("error doing config.GetPgDbDsnUrlFromEnv. error: %v", err))
	}

	InitialDB, err := todos.GetStorageInstance("postgres", dbDsn, l)
	if err != nil {
		t.Fatalf(fmt.Sprintf("error getting storage. is postgres available ? error : %v ", err))
	}
	defer InitialDB.Close()
	myServer := GetNewServer(true, l, InitialDB)
	ts := httptest.NewServer(myServer)
	defer ts.Close()

	res, _ := InitialDB.List(0, 100)

	jsonInitialData, _ := json.Marshal(res)
	fmt.Printf("%#v", res)
	firstTodo, _ := json.Marshal(res[0])
	maxId, _ := InitialDB.GetMaxId()

	tests := getTestTable(t, ts, idCounter{currentMaxId: maxId}, jsonInitialData, firstTodo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			resp, err := http.DefaultClient.Do(tt.r)
			if DEBUG {
				fmt.Printf("### %s : %s on %s\n", tt.name, tt.r.Method, tt.r.URL)
			}
			defer resp.Body.Close()
			if err != nil {
				fmt.Printf("### GOT ERROR : %s\n%s", err, resp.Body)
				t.Fatal(err)
			}
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode, "expected status code should be returned")
			receivedJson, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusCreated {
				// In case of Created we need to modify the CreatedAt value from the response for equality test to pass
				myNewTodo := todos.Todo{
					Completed:   false,
					CompletedAt: nil,
					CreatedAt:   nil,
					Id:          0,
					Task:        defaultNewTask,
				}
				var createdTodo todos.Todo
				err := json.Unmarshal(receivedJson, &createdTodo)
				if err != nil {
					fmt.Printf("FATAL ERROR doing json.Unmarshall of %#v   - Error: %s", string(receivedJson), err)
				}
				myNewTodo.CreatedAt = createdTodo.CreatedAt
				myNewTodo.Id = createdTodo.Id
				wantedTodoJson, _ := json.Marshal(myNewTodo)
				if DEBUG {
					fmt.Printf("WANTED   :%T - %#v\n", wantedTodoJson, string(wantedTodoJson))
					fmt.Printf("RECEIVED :%T - %#v\n", receivedJson, string(receivedJson))
				}
				assert.JSONEqf(t, string(wantedTodoJson), string(receivedJson), "CreateTodo Response was not equal to expected.")
			} else {
				// here are all other cases (except Created above)
				if DEBUG {
					fmt.Printf("WANTED   :%T - %#v\n", tt.wantBody, tt.wantBody)
					fmt.Printf("RECEIVED :%T - %#v\n", receivedJson, string(receivedJson))
				}
				// check that receivedJson contains the specified tt.wantBody substring . https://pkg.go.dev/github.com/stretchr/testify/assert#Contains
				assert.Contains(t, string(receivedJson), tt.wantBody, "Response should contain what was expected.")
			}
		})
	}
}
