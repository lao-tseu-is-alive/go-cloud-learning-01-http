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

// helloWorldHandler is a simple http handler to give some greetings with a valid html
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// retrieve a parameter from query
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, fmt.Sprintf("Bad request: Error in ParseQuery: %q", err), http.StatusBadRequest)
			return
		}
		username := defaultUserName
		if query.Has("username") {
			username = query.Get("username")
		}
		username = strings.TrimSpace(username)
		if len(username) == 0 {
			http.Error(w, fmt.Sprintf("Bad request: Error in query.Get('username'): username cannot be empty or spaces only"), http.StatusBadRequest)
			return
		}
		res, err := getHelloMsg(username)
		if err != nil {
			http.Error(w, fmt.Sprintf("Internal server error:%q", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, res)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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

	/* section above is to illustrate and test the template
	msg, err := getHelloMsg("toto")
	if err != nil {
		fmt.Errorf("ERROR doing getHelloMsg : %q", err)
	}
	fmt.Println(msg)
	*/

	http.HandleFunc("/hello", helloWorldHandler)
	log.Printf("Starting server... try navigating to http://localhost%v/hello to be greeted", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
