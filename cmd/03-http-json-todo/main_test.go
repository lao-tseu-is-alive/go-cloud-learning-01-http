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

func Test_goTodoServer_Todos(t *testing.T) {
	// Create server using the router initialized elsewhere. The router
	// can be a net/http ServeMux a http.DefaultServeMux or
	// any value that satisfies the net/http Handler interface.
	ts := httptest.NewServer(getNewServer(true))
	defer ts.Close()

	const (
		defaultNewTask = "Learn Linux"
		DEBUG          = false
	)
	currentMaxId := defaultMaxId

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
			name:           "1: CreateTodo valid new Todo task, should return a valid Todo",
			wantStatusCode: http.StatusCreated,
			wantBody:       `{"completed":false,"id":` + fmt.Sprintf("%d", currentMaxId+1) + ` ,"task":"` + defaultNewTask + `"}`,
			r:              newRequest(http.MethodPost, "/todos", `{"task":"`+defaultNewTask+`"}`),
		},
		{
			name:           "1: CreateTodo valid new Todo task",
			wantStatusCode: http.StatusCreated,
			wantBody:       `{"completed":false,"id":` + fmt.Sprintf("%d", currentMaxId+1) + ` ,"task":"` + defaultNewTask + `"}`,
			r:              newRequest(http.MethodPost, "/todos", `{"task":"`+defaultNewTask+`"}`),
		},
		{
			name:           "2: CreateTodo invalid request (task field is of wrong type in body)",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo has invalid format",
			r:              newRequest(http.MethodPost, "/todos", `{"task":123}`),
		},
		{
			name:           "2: CreateTodo invalid request (empty body)",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task cannot be empty\"}\n",
			r:              newRequest(http.MethodPost, "/todos", `{}`),
		},
		{
			name:           "2: CreateTodo invalid request (missing task field in body)",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task cannot be empty\"}",
			r:              newRequest(http.MethodPost, "/todos", `{"nope":"should fail"}`),
		},
		{
			name:           "4: CreateTodo empty new Todo task",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task cannot be empty\"}",
			r:              newRequest(http.MethodPost, "/todos", `{"task":""}`),
		},
		{
			name:           "3: CreateTodo invalid new Todo task (too short)",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "{\"message\":\"CreateTodo task minLength is 5\"}",
			r:              newRequest(http.MethodPost, "/todos", `{"task":"123"}`),
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
			fmt.Printf("### %s : %s on %s\n", tt.name, tt.r.Method, tt.r.URL)
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
				if DEBUG {
					fmt.Printf("WANTED   :%T - %#v\n", tt.wantBody, tt.wantBody)
					fmt.Printf("RECEIVED :%T - %#v\n", receivedJson, string(receivedJson))
				}
				assert.Contains(t, string(receivedJson), tt.wantBody, "CreateTodo Response contains what was expected.")
			}
		})
	}
}
