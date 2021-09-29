package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const defaultPort = 8080

// your first golang web server
// to try it just type : go run main.go
func main() {
	listenAddr := fmt.Sprintf(":%v", defaultPort)

	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	})
	log.Printf("Starting server... try navigating to http://localhost%v/hello to be greeted", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
