package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloWorldHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	response := httptest.NewRecorder()

	helloWorldHandler(response, request)
	// what about status code ??
	t.Run("returns the hello message", func(t *testing.T) {
		got := response.Body.String()
		want := helloMsg

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
