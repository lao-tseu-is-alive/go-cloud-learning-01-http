package main

import (
	"io"
	"log"
	"net/http"
)

const helloMsg = "Hello, world!"

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, helloMsg)
}

// basic golang web server
// to try it just type : go run main.go
func main() {
	listenAddr := ":8080"
	http.HandleFunc("/hello", helloWorldHandler)
	log.Printf("Starting server... try navigating to http://localhost%v/hello to be greeted", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
