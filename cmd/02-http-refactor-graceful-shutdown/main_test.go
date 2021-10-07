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
	//l := log.New(ioutil.Discard, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	l := log.New(os.Stdout, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	goSrv := NewGoHttpServer(defaultServerIp, defaultServerPort, l)
	helloWorldHandler := goSrv.getHelloHandler()

	tt := []struct {
		name           string
		url            string
		method         string
		paramKeyValues map[string]string
		want           string
		statusCode     int
	}{
		{
			name:           "without any username parameter, we want default message",
			url:            "/hello",
			method:         http.MethodGet,
			paramKeyValues: make(map[string]string, 0),
			want:           defaultMsg,
			statusCode:     http.StatusOK,
		},
		{
			name:           "with username having a valid value, we want greeting with username",
			url:            "/hello",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "Carlos"},
			want:           "", // let's calculate the result later based on given userName
			statusCode:     http.StatusOK,
		},
		{
			name:           "with username having an empty value, we want 400 Bad request",
			url:            "/hello",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": ""},
			want:           "Bad request. In query.Get('username'): username cannot be empty or spaces only\n",
			statusCode:     http.StatusBadRequest,
		},
		{
			name:           "with unsupported http method, we want 405 Method not allowed",
			method:         http.MethodPost,
			paramKeyValues: map[string]string{"username": "WhatEverYouWant"},
			want:           "Method not allowed\n",
			statusCode:     http.StatusMethodNotAllowed,
		},
		{ // curl  'http://localhost:8080/hello?username=WhatEver;YouWant will not work, use instead :
			// curl -G --data-urlencode 'username=WhatEver;YouWant' http://localhost:8080/hello
			// translates to curl  'http://localhost:8080/hello?username=WhatEver%3YouWant
			name:           "with username containing semicolon url encoded, we want 200 with default message",
			url:            "/hello",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "What Ever;You Want"},
			want:           "", // let's calculate the result later based on given userName
			statusCode:     http.StatusOK,
		},
		{ // curl  'http://localhost:8080/hello?username=ℂ𝕠𝕦𝕔𝕠𝕦 𝕝𝕖 𝕙𝕚𝕓𝕠𝕦' will not work, use instead :
			// curl -G --data-urlencode 'username=ℂ𝕠𝕦𝕔𝕠𝕦 𝕝𝕖 𝕙𝕚𝕓𝕠𝕦' http://localhost:8080/hello
			name:           "with username containing unicode, we want 200 with default message",
			url:            "/hello",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "ℂ𝕠𝕦𝕔𝕠𝕦 𝕝𝕖 𝕙𝕚𝕓𝕠𝕦"},
			want:           "", // let's calculate the result later based on given userName
			statusCode:     http.StatusOK,
		},
		{ // curl -G --data-urlencode 'username=╚»☯💥⚡✌ℂ𝔾𝕀𝕃✌⚡💥☯«╝' http://localhost:8080/hello
			name:           "with username containing unicode, we want 200 with default message",
			url:            "/hello",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "╚»☯💥⚡✌ℂ𝔾𝕀𝕃✌⚡💥☯«╝"},
			want:           "", // let's calculate the result later based on given userName
			statusCode:     http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(tc.method, tc.url, nil)
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
			log.Printf("## IN TEST %s\n## path : %#v", tc.name, request.URL.Path)
			helloWorldHandler(response, request)
			got := response.Body.String()
			want := tc.want
			if len(want) == 0 {
				want, _ = getHelloMsg(username)
			}
			// testing the response body
			if got != want {
				t.Errorf("### in response.Body got :\n%q\n### want :\n %q", got, want)
			}
			// testing the server status code
			if response.Code != tc.statusCode {
				t.Errorf("in response.Code got status : %d, want status %d", response.Code, tc.statusCode)
			}

		})
	}
}

func TestMyDefaultHandler(t *testing.T) {

	okResponse := getHtmlPage(defaultMessage)
	notFoundResponse := getHtmlPage(defaultNotFound)
	l := log.New(ioutil.Discard, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	goSrv := NewGoHttpServer(defaultServerIp, defaultServerPort, l)
	defaultHandler := goSrv.getMyDefaultHandler()
	//serverTest := httptest.NewServer(goSrv.router)
	//defer serverTest.Close()

	tt := []struct {
		name           string
		url            string
		method         string
		paramKeyValues map[string]string
		want           string
		statusCode     int
	}{
		{
			name:           "without any parameter, we want default message",
			url:            "/",
			method:         http.MethodGet,
			paramKeyValues: make(map[string]string, 0),
			want:           okResponse,
			statusCode:     http.StatusOK,
		},
		{
			name:           "with username having a valid value, we want default message",
			url:            "/",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "Carlos"},
			want:           okResponse,
			statusCode:     http.StatusOK,
		},
		{
			name:           "with unsupported http method, we want 405 Method not allowed",
			url:            "/",
			method:         http.MethodPost,
			paramKeyValues: map[string]string{"username": "WhatEverYouWant"},
			want:           "Method not allowed\n",
			statusCode:     http.StatusMethodNotAllowed,
		},
		{
			name:           "with nonexistent url query,we want 404 page not found",
			url:            "/nope",
			method:         http.MethodGet,
			paramKeyValues: map[string]string{"username": "WhatEverYouWant"},
			want:           notFoundResponse,
			statusCode:     http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(tc.method, tc.url, nil)
			response := httptest.NewRecorder()
			log.Printf("## IN TEST %s\n## path : %#v", tc.name, request.URL.Path)
			defaultHandler(response, request)
			got := response.Body.String()
			want := tc.want
			// testing the response body
			if got != want {
				t.Errorf("### in response.Body got :\n%q\n### want :\n %q", got, want)
			}
			// testing the server status code
			if response.Code != tc.statusCode {
				t.Errorf("in response.Code got status : %d, want status %d", response.Code, tc.statusCode)
			}

		})
	}
}
