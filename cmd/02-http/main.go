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

func handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}
}

// your first golang web server
// to try it just type : go run main.go
func main() {
	listenAddr := fmt.Sprintf(":%v", defaultPort)
	// check ENV PORT for the good PORT TO USE
	val, exist := os.LookupEnv("PORT")
	if exist {
		port, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal("ERROR: PORT ENV should contain a valid integer value !")
		}
		listenAddr = fmt.Sprintf(":%v", port)
	}

	http.HandleFunc("/hello", handleHello())
	log.Printf("Starting server... try navigating to http://localhost%v/hello to be greeted", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
