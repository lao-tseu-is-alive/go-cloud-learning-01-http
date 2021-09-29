package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = 8080

// just a simple http handler to give some greetings
func handleHello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

// let's improve our http server : we will allow it to read
// the listening port to use from an environment  variable
// to try it just type 		: go run main.go
// or with a different port :  WEB_PORT=3333 go run main.go
// you can also try 		:  WEB_PORT=XXX3333 go run main.go
func main() {
	listenAddr := fmt.Sprintf(":%v", defaultPort)
	// check ENV WEB_PORT for the PORT to use for listening to connection
	val, exist := os.LookupEnv("WEB_PORT")
	if exist {
		port, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal("ERROR: WEB_PORT ENV should contain a valid integer value !")
		}
		listenAddr = fmt.Sprintf(":%v", port)
	}

	http.HandleFunc("/hello", handleHello)
	log.Printf("Starting server... try navigating to http://localhost%v/hello to be greeted", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
