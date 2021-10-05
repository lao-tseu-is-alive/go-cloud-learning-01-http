package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestHelloHandler allows to check that the HelloWorldHandler works as expected
// just run : go test -race -covermode=atomic -coverprofile=coverage.out
func TestHelloHandler(t *testing.T) {

	defaultMsg, _ := getHelloMsg(defaultUserName)
	l := log.New(os.Stdout, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	l.SetOutput(ioutil.Discard)
	server := NewGoHttpServer(defaultServerIp, defaultServerPort, l)
	server.StartServer()

	tt := []struct {
		name           string
		method         string
		paramKeyValues map[string]string
		want           string
		statusCode     int
	}{
		{
			name:           "without any username parameter, we want default message",
			method:         http.MethodGet,
			paramKeyValues: make(map[string]string, 0),
			want:           defaultMsg,
			statusCode:     http.StatusOK,
		},
		{
			name:           "with username having a valid value, we want greeting with username",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "Carlos"},
			want:           "", // let's calculate the result later based on given userName
			statusCode:     http.StatusOK,
		},
		{
			name:           "with username having an empty value, we want 400 Bad request",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": ""},
			want:           "Bad request. In query.Get('username'): username cannot be empty or spaces only\n",
			statusCode:     http.StatusBadRequest,
		},
		{
			name:           "with unsupported http method, we want 405 Method not allowed",
			method:         http.MethodPost,
			paramKeyValues: map[string]string{"username": "WhatEverYouWant", "param2": "nobody is here"},
			want:           "Method not allowed\n",
			statusCode:     http.StatusMethodNotAllowed,
		},
		// for next test to pass we need to test the server not only the handler
		{
			name:           "with invalid query, we want 400 Bad request",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "WhatEverYouWant", "param2": "nobody; is here"},
			want:           "Bad request. Error in ParseQuery",
			statusCode:     http.StatusBadRequest,
		},
		{
			name:           "with nonexistent url query,we want 404 page not found",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "WhatEverYouWant", "param2": "nobody; is here"},
			want:           "Bad request. Error in ParseQuery",
			statusCode:     http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(tc.method, "/", nil)
			username := ""
			if len(tc.paramKeyValues) > 0 {
				parameters := request.URL.Query()
				for paramName, paramValue := range tc.paramKeyValues {
					parameters.Add(paramName, paramValue)
					if paramName == "username" {
						username = paramValue
					}
				}
				request.URL.RawQuery = parameters.Encode()
			}
			response := httptest.NewRecorder()
			server.router.ServeHTTP(response, request)
			got := response.Body.String()
			want := tc.want
			if len(want) == 0 {
				// let's calculate the result later based on current username in test array
				want, _ = getHelloMsg(username)
			}
			// testing the response body
			if got != want {
				t.Errorf("in response.Body got %q, want %q", got, want)
			}
			// testing the server status code
			if response.Code != tc.statusCode {
				t.Errorf("in response.Code got %d, want status %d", response.Code, tc.statusCode)
			}

		})
	}
}
