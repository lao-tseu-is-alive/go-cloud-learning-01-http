package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_goTodoServer_CreateTodo(t *testing.T) {

	l := log.New(ioutil.Discard, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	DBinMemory := initializeStorage()
	myApi := goTodoServer{logger: l, data: &DBinMemory}
	e := echo.New()
	RegisterHandlers(e, &myApi)

	myNewTodo := Todo{
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   nil,
		Id:          3,
		Task:        "learn Linux",
	}
	myNewTodoJson := `{"task":"learn Linux"}`
	myInvalidNewTodoJson := `{"task":"xyz"}`
	myEmptyNewTodoJson := `{"task":""}`
	myTodoCreated := `{"completed":false,"id":3,"task":"learn Linux"}`

	tests := []struct {
		name           string
		url            string
		method         string
		paramKeyValues map[string]string
		jsonData       string
		want           string
		statusCode     int
	}{
		{
			name:           "CreateTodo receiving a valid NewTodo, should return a valid Todo with a status code 201 ",
			url:            "/todos",
			method:         http.MethodPost,
			paramKeyValues: nil,
			jsonData:       myNewTodoJson,
			want:           myTodoCreated,
			statusCode:     http.StatusCreated,
		},
		{
			name:           "CreateTodo with Task length < 6, should return an Error Bad request",
			url:            "/todos",
			method:         http.MethodPost,
			paramKeyValues: nil,
			jsonData:       myInvalidNewTodoJson,
			want:           "",
			statusCode:     http.StatusBadRequest,
		},
		{
			name:           "CreateTodo with a blank Task, should return Error Bad request",
			url:            "/todos",
			method:         http.MethodPost,
			paramKeyValues: nil,
			jsonData:       myEmptyNewTodoJson,
			want:           "",
			statusCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.jsonData))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := &myApi

			// Assertions for normal cases
			//if tt.statusCode == http.StatusCreated {
			if assert.NoError(t, h.CreateTodo(c)) {
				/*err := h.CreateTodo(c)
				if err != nil {
					fmt.Printf("### Notice in CreateTodo err: %#v \n### Status Code : %d \n### Response : %s\n", err, rec.Code, rec.Body.String())
					return
				}*/
				assert.Equal(t, tt.statusCode, rec.Code)
				// add the created at field from the response so that test passes
				createdTodo := new(Todo)
				err := json.Unmarshal(rec.Body.Bytes(), &createdTodo)
				if err != nil {
					panic(err)
				}
				myNewTodo.CreatedAt = createdTodo.CreatedAt
				/*
					fmt.Printf("%T - %#v\n", createdTodo, createdTodo)
					fmt.Printf("%T - %#v\n", myNewTodo, createdTodo)
				*/
				wantedTodoJson, _ := json.Marshal(myNewTodo)
				assert.JSONEqf(t, string(wantedTodoJson), rec.Body.String(), "CreateTodo Response was not equal to expected.")
			}
		})
	}
}
