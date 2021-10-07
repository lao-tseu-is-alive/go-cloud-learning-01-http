## Your first version of a golang web server

The most basic web server you can imagine in Go:

*In a few lines we define a web server, 
that listen on a fixed port and returns a classic "Hello World !" 
when a specific url request is entered*

```go
package main

import (
	"io"
	"net/http"
)
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!")
}

func main() {
	http.HandleFunc("/hello", helloWorldHandler)
	http.ListenAndServe(":8080", nil)
}
```

## What to remember :
- Go comes with batteries included : the standard library [net/http](https://pkg.go.dev/net/http) provides HTTP client and server implementations.
- The net/http package uses [Handlers](https://pkg.go.dev/net/http#Handler) to handle HTTP requests which are sent to a specific path.
- Creating a function with the signature func (w ResponseWriter, r *Request) and use **http.HandleFunc(path, function)** to register it.


### How to run it ?

to run it just type :
```bash
go run main.go &
curl http://localhost:8080/hello
```

to build a binary executable called "mywebserver", just type :
```bash
go build -o mywebserver main.go 
```

to unit test the handler just type :
```bash
go test 
```

to check on your server wich process is listening on a specific port :
```bash
ss -lntp
```


### More information :
- [Demystifying HTTP Handlers in Golang](https://medium.com/geekculture/demystifying-http-handlers-in-golang-a363e4222756)