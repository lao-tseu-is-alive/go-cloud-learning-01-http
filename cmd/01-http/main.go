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

const helloMsg = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css" integrity="sha512-EZLkOqwILORob+p0BXZc+Vm3RgJBOe1Iq/0fiI7r/wJgzOFZMlsqTa29UEl6v6U6gsV4uIpsNZoV32YZqrCRCQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
    <title>GOLANG Hello World</title>
  </head>
  <body>
	<h3>Hello, world!</h3>
  </body>
</html>
`

// helloWorldHandler is a simple http handler to give some greetings with a valid html
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, helloMsg)
}

// let's improve our http server : we will allow it to read
// the listening port to use from an environment  variable
// and also
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

	http.HandleFunc("/hello", helloWorldHandler)
	log.Printf("Starting server... try navigating to http://localhost%v/hello to be greeted", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
