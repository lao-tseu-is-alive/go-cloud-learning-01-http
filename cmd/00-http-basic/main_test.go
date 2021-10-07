package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHelloWorldHandler allows to check that the HelloWorldHandler works as expected
// just run : go test
func TestHelloWorldHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	response := httptest.NewRecorder()

	helloWorldHandler(response, request)

	t.Run("returns the hello message", func(t *testing.T) {
		got := response.Body.String()
		want := helloMsg

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
