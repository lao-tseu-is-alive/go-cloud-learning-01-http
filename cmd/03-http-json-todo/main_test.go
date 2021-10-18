package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type idCounter struct {
	currentMaxId int32
}

func (c *idCounter) increment() int32 {
	c.currentMaxId = c.currentMaxId + 1
	fmt.Printf("# IN increment() currentMaxId: %d\n", c.currentMaxId)
	return c.currentMaxId
}

func (c *idCounter) decrement() int32 {
	c.currentMaxId = c.currentMaxId - 1
	fmt.Printf("# IN decrement() currentMaxId: %d\n", c.currentMaxId)
	return c.currentMaxId
}

func (c *idCounter) current() int32 {
	fmt.Printf("# IN current() currentMaxId: %d\n", c.currentMaxId)
	return c.currentMaxId
}

func (c *idCounter) currentAsString() string {
	fmt.Printf("# IN currentAsString() currentMaxId: %d\n", c.currentMaxId)
	return fmt.Sprintf("%d", c.currentMaxId)
}

func getUrlForId(myIdCounter idCounter) string {
	return fmt.Sprintf("/todos/%d", myIdCounter.current())
}

func Test_goTodoServer_Todos(t *testing.T) {
	// Create server using the router initialized elsewhere. The router
	// can be a net/http ServeMux a http.DefaultServeMux or
	// any value that satisfies the net/http Handler interface.
	ts := httptest.NewServer(getNewServer(true))
	defer ts.Close()

	const (
		defaultNewTask = "Learn Linux"
		DEBUG          = true
	)
	myId := idCounter{currentMaxId: defaultMaxId}
	InitialDB := initializeStorage()
	jsonInitialData, _ := json.Marshal(InitialDB.Todos)

	newRequest := func(method, url string, body string) *http.Request {
		r, err := http.NewRequest(method, ts.URL+url, strings.NewReader(body))
		if err != nil {
			t.Fatalf("### ERROR http.NewRequest %s on [%s] \n", method, url)
		}
		return r
	}

	tests := []struct {
		name           string
		wantStatusCode int
		wantBody       string
		r              *http.Request
	}{
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
			name:           "15: UpdateTodo with completed=false, should return a Todo updated with completed=false",
			wantStatusCode: http.StatusOK,
			wantBody:       `"completed":false`,
			r:              newRequest(http.MethodPut, getUrlForId(myId), `{"completed":false,"id":`+myId.currentAsString()+` ,"task":"`+defaultNewTask+`"}`),
		},
		{
			name:           "15: UpdateTodo with empty task, will return a Bad request",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task cannot be empty\"}",
			r:              newRequest(http.MethodPut, getUrlForId(myId), `{"completed":false,"id":`+myId.currentAsString()+` ,"task":""}`),
		},
		{
			name:           "99:  invalid path, should return 404 not found",
			wantStatusCode: http.StatusNotFound,
			wantBody:       "{\"message\":\"Not Found\"}",
			r:              newRequest(http.MethodGet, "/nothing_available_here", `{"task":"123"}`),
		},
	}
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
				myNewTodo := Todo{
					Completed:   false,
					CompletedAt: nil,
					CreatedAt:   nil,
					Id:          0,
					Task:        defaultNewTask,
				}
				var createdTodo Todo
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
				assert.Contains(t, string(receivedJson), tt.wantBody, "CreateTodo Response contains what was expected.")
			}
		})
	}
}
