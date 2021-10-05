package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const defaultPort = 8080
const defaultUserName = "World"

func getHelloMsg(name string) (string, error) {
	const helloMsg = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css" integrity="sha512-EZLkOqwILORob+p0BXZc+Vm3RgJBOe1Iq/0fiI7r/wJgzOFZMlsqTa29UEl6v6U6gsV4uIpsNZoV32YZqrCRCQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
    <title>GOLANG Hello {{.UserName}}</title>
  </head>
  <body>
	<h3>Hello, {{.UserName}}!</h3>
  </body>
</html>
`

	data := struct {
		UserName string
	}{UserName: name}
	var tpl bytes.Buffer
	t, err := template.New("hello-page").Parse(helloMsg)
	if err != nil {
		return "", err
	}
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func getHelloHandler() http.HandlerFunc {
	// this part is executed only once at the initial launch of your server
	// os it is a good place to run initialisation stuff

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// how to retrieve a parameter from query with standard library
			query, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				log.Printf("ERROR - Bad request. Error in ParseQuery. Error: %q . Request was: \n%#v\n", err, r)
				http.Error(w, fmt.Sprintf("Bad request. Error in ParseQuery: %q", err), http.StatusBadRequest)
				return
			}
			username := defaultUserName
			if query.Has("username") {
				username = query.Get("username")
			}
			username = strings.TrimSpace(username)
			if len(username) == 0 {
				log.Printf("ERROR - Bad request. Username cannot be empty. Request was: \n%#v\n", r)
				http.Error(w, fmt.Sprintf("Bad request. In query.Get('username'): username cannot be empty or spaces only"), http.StatusBadRequest)
				return
			}
			res, err := getHelloMsg(username)
			if err != nil {
				log.Printf("ERROR - Internal server error. Request was: \n%#v\n", r)
				http.Error(w, fmt.Sprintf("Internal server error. Error: %q", err), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, res)
		default:
			log.Printf("ERROR - Method not allowed. Request: %#v", r)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getListenAddrFromEnv(defaultPort int) string {
	listenAddress := fmt.Sprintf(":%v", defaultPort)
	val, exist := os.LookupEnv("PORT")
	if exist {
		port, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal("ERROR: CONFIG ENV PORT should contain a valid integer value !")
		}
		listenAddress = fmt.Sprintf(":%v", port)
	}
	return listenAddress
}

func main() {
	listenAddr := getListenAddrFromEnv(defaultPort)

	http.HandleFunc("/hello", getHelloHandler())
	log.Printf("Starting server... try navigating to http://localhost%v/hello to be greeted", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
