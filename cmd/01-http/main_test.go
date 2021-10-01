package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHelloWorldHandler allows to check that the HelloWorldHandler works as expected
// just run : go test
func TestHelloWorldHandler(t *testing.T) {

	defaultMsg, _ := getHelloMsg(defaultUserName)

	tt := []struct {
		name           string
		method         string
		paramKeyValues map[string]string
		want           string
		statusCode     int
	}{
		{
			name:           "without any username parameter",
			method:         http.MethodGet,
			paramKeyValues: make(map[string]string, 0),
			want:           defaultMsg,
			statusCode:     http.StatusOK,
		},
		{
			name:           "with username having a valid value",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "Carlos"},
			want:           "", // let's calculate the result later based on given userName
			statusCode:     http.StatusOK,
		},
		{
			name:           "with username having an empty value",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": ""},
			want:           "Bad request: Error in query.Get('username'): username cannot be empty or spaces only\n",
			statusCode:     http.StatusBadRequest,
		},
		{
			name:           "with bad method",
			method:         http.MethodPost,
			paramKeyValues: map[string]string{"username": "WhatEverYouWant", "param2": "nobody is here"},
			want:           "Method not allowed\n",
			statusCode:     http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(tc.method, "/hello", nil)
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
			helloWorldHandler(response, request)
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
