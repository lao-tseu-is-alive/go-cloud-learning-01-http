## Your third version of a golang web server

In this next iteration of our hello web server we will introduce those features:

1. Refactor the code using dedicated struct type to store information related to all handlers of web server
2. Implement graceful shuts down the server without interrupting any active connections.
3. Implement a default handler to intercepts and handle all request to nonexistent endpoints.
4. Learn how to configure and use a logger


### Refactor the code around a dedicated struct type

Using a dedicated struct type, will allow you to store information useful  to all the handlers of web server. 

```go
package main

import (
	"log"
	"net/http"
	"os"
)

const (
	defaultServerPort   = 8080
	defaultServerIp     = "127.0.0.1" //safe default for listening ip 
)

//getListenAddrFromEnv returns a valid TCP/IP listening address string based on
// the values of ENV SERVERIP:PORT
func getListenAddrFromEnv(defaultIP string, defaultPort int) string {
	// some code goes her to retrieve the values of env variables or use the default
	return ""
}


//GoHttpServer is a struct type to store information related to all web handlers 
type GoHttpServer struct {
	listenAddress string
	// later we will store here the connection to database
	//DB  *db.Conn
	logger *log.Logger
	router *http.ServeMux
}

//NewGoHttpServer is a constructor that initializes the server mux (routes) 
//and all fields of the  GoHttpServer type
func NewGoHttpServer(ServerIpDefault string, ServerPortDefault int, logger *log.Logger) *GoHttpServer {
	
	myServer := GoHttpServer{
		listenAddress: getListenAddrFromEnv(ServerIpDefault, ServerPortDefault),
		logger:        logger,
		router:        http.NewServeMux(),
	}
	myServer.routes()
	return &myServer
}

// (*GoHttpServer) routes initializes all the handlers paths of this web server, it is called inside the NewGoHttpServer constructor
func (s *GoHttpServer) routes() {
	// ALL your routes are defined in one place easy to find  
	s.router.Handle("/", s.getMyDefaultHandler())
	s.router.Handle("/hello", s.getHelloHandler())
	s.router.Handle("/slowHello", s.getSlowHelloHandler())
}

// StartServer initializes all the handlers paths of this web server, it is called inside the NewGoHttpServer constructor
func (s *GoHttpServer) StartServer() { 	
	// all the code to start the server handling graceful shutdown etc...
}

func (s *GoHttpServer) getHelloHandler() http.HandlerFunc {
	// this part is executed only once at the initial launch of your server
	//os it is a good place to run initialisation stuff
	s.logger.Println("INITIAL CALL TO getHelloHandler()")
	return func(w http.ResponseWriter, r *http.Request) {
		// more code goes here ...
		return
	}
}

func (s *GoHttpServer) getMyDefaultHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// handler code goes here ...
		return
	}
}

func (s *GoHttpServer) getSlowHelloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}



func main() {
	// much clean place now no ?
	l := log.New(os.Stdout, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	server := NewGoHttpServer(defaultServerIp, defaultServerPort, l)
	server.StartServer()
}

```

### What to remember :
- It is important to handle graceful shutdown for long queries to close cleanly
- Use a dedicated type to store all information related to your web server, like DB connection, 
 and share them with your handlers 
- A logger allows you to get more information from your code "on demand". 
It can be easily configured to send output to a file, or disable output completely. 
**In production log output (particularly on file) can have a big impact on performance 
because of concurrent execution of request connections** 
Have a look on the logger implementation and you probably will find a Mutex
to control concurrency access...   



### How to run it ?

to run the server and listen on port 3333  just type :
```bash
PORT=3333 go run main.go 
#in another terminal test it with curl or open a browser
curl http://localhost:3333/hello
curl  http://localhost:3333/hello?username=Rob%20Pike
#check what happens if you use another HTTP verb like POST
curl -XPOST  http://localhost:3333/hello?username=Rob%20Pike
curl -XPOST -d '{"username":"toto"}' http://localhost:3333/hello?username=Rob%20Pike
curl -XPUT -d '{"username":"toto"}' http://localhost:3333/nohandlersForThisRoute
curl  http://localhost:3333/nohandlersForThisRoute
# now that we have a default handler what appears now in the log for the last line ?

# testing the graceful shutdown :
curl  http://localhost:3333/slowHello
# And just after run CTRL+C on the server terminal,
# to experiment  a gracefull shutdown of running request
# if you are fast enough you can from a third terminal try to connect after the CTRL+C :
curl  http://localhost:3333/slowHello
# the new connections are refused but the active ones (before the CTRL-C signal) 
# will terminate as long as they does not last more then the timeout


```

to build a binary executable called "mywebserver", just type :
```bash
go build -o mywebserver main.go 
```

to unit test the handler just type, by the way have a look at the main_test.go code :
```bash
go test -race -covermode=atomic -coverprofile=coverage.out 
```

to check on your server which process is listening on a specific port :
```bash
ss -lntp
```


### More information :
- [Go http server Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)
- [Graceful shutdown in Go http server](https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a)
- [Logging in Go: Choosing a System and Using it](https://www.honeybadger.io/blog/golang-logging/)
- [The http.Handler wrapper technique in GO](https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702)