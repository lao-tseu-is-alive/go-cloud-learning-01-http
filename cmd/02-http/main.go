package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const (
	defaultServerPort   = 8080
	defaultServerIp     = "127.0.0.1"
	shutDownTimeout     = 15 * time.Second // number of second to wait before closing server
	defaultReadTimeout  = 2 * time.Minute
	defaultWriteTimeout = 2 * time.Minute
	defaultIdleTimeout  = 2 * time.Minute
	defaultUserName     = "World"
	secondsToSleep      = 5
)

func getListenAddrFromEnv(defaultIP string, defaultPort int) string {
	srvIP := defaultIP
	srvPort := defaultPort
	var err error
	val, exist := os.LookupEnv("PORT")
	if exist {
		srvPort, err = strconv.Atoi(val)
		if err != nil {
			log.Fatal("ERROR: CONFIG ENV PORT should contain a valid integer value !")
		}
	}
	val, exist = os.LookupEnv("SERVERIP")
	if exist {
		srvIP = val
	}
	return fmt.Sprintf("%s:%d", srvIP, srvPort)
}

type GoHttpServer struct {
	listenAddress string
	// later we will store here the connection to database
	//DB  *db.Conn
	logger *log.Logger
	router *http.ServeMux
}

func NewGoHttpServer(ServerIpDefault string, ServerPortDefault int, logger *log.Logger) *GoHttpServer {
	listenAddress := getListenAddrFromEnv(ServerIpDefault, ServerPortDefault)
	myServer := GoHttpServer{
		listenAddress: listenAddress,
		logger:        logger,
		router:        http.NewServeMux(),
	}
	myServer.routes()
	return &myServer
}

func (s *GoHttpServer) routes() {
	s.router.Handle("/", s.getMyDefaultHandler())
	s.router.Handle("/hello", s.getHelloHandler())
	s.router.Handle("/slowHello", s.getSlowHelloHandler())
}

func (s *GoHttpServer) StartServer() {

	// create a new http server
	srv := http.Server{
		Addr:         s.listenAddress,     // configure the bind address
		Handler:      s.router,            // set the default handler
		ErrorLog:     s.logger,            // set the logger for the server
		ReadTimeout:  defaultReadTimeout,  // max time to read request from the client
		WriteTimeout: defaultWriteTimeout, // max time to write response to the client
		IdleTimeout:  defaultIdleTimeout,  // max time for connections using TCP Keep-Alive
	}

	// Starting the web server in his own goroutine
	go func() {
		s.logger.Printf("## Starting server... try navigating to http://localhost%v/hello to be greeted", s.listenAddress)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Error starting server: %s\n", err)
		}
	}()
	s.logger.Printf("Server listening on : %s PID:[%d]", srv.Addr, os.Getpid())

	// Wait for interrupt signal and gracefully shutdown the server with a timeout of n seconds.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	signal.Notify(stopChan, os.Kill)

	// Block until a signal is received.
	// wait for SIGINT (interrupt) 	: ctrl + C keypress, or in a shell : kill -SIGINT processId
	sig := <-stopChan

	s.logger.Printf("SIGINT %d interrupt signal received, about to shut down server after max %v sec...\n", sig, shutDownTimeout.Seconds())

	ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)

	// gracefully shuts down the server without interrupting any active connections
	// as long as the actives connections last less than shutDownTimeout
	// https://pkg.go.dev/net/http#Server.Shutdown
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Printf("Problem doing Shutdown %v", err)
	}
	<-ctx.Done()
	cancel()
	s.logger.Print("Server gracefully stopped. Bye-bye")
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	//ctx, _ := context.WithTimeout(context.Background(), shutDownTimeout)
	//srv.Shutdown(ctx)

}

func getHelloMsg(name string) (string, error) {
	const helloMsg = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
	<link rel="stylesheet" 
		href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css"
		integrity="sha512-EZLkOqwILORob+p0BXZc+Vm3RgJBOe1Iq/0fiI7r/wJgzOFZMlsqTa29UEl6v6U6gsV4uIpsNZoV32YZqrCRCQ==" 
		crossorigin="anonymous" referrerpolicy="no-referrer" />
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

func (s *GoHttpServer) getHelloHandler() http.HandlerFunc {
	// this part is executed only once at the initial launch of your server
	// os it is a good place to run initialisation stuff
	s.logger.Println("INITIAL CALL TO getHelloHandler()")

	return func(w http.ResponseWriter, r *http.Request) {
		remoteIp := r.RemoteAddr
		requestedUrl := r.RequestURI
		switch r.Method {
		case http.MethodGet:
			// how to retrieve a parameter from query with standard library
			query, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				s.logger.Printf("ERROR: Bad request. Error in ParseQuery. Error: %q . Request was: \n%#v\n", err, r)
				http.Error(w, fmt.Sprintf("Bad request. Error in ParseQuery: %q", err), http.StatusBadRequest)
				return
			}
			username := defaultUserName
			if query.Has("username") {
				username = query.Get("username")
			}
			username = strings.TrimSpace(username)
			if len(username) == 0 {
				s.logger.Printf("ERROR: Bad request. Username cannot be empty. Request was: \n%#v\n", r)
				http.Error(w, fmt.Sprintf("Bad request. In query.Get('username'): username cannot be empty or spaces only"), http.StatusBadRequest)
				return
			}
			res, err := getHelloMsg(username)
			if err != nil {
				s.logger.Printf("ERROR: Internal server error. Request was: \n%#v\n", r)
				http.Error(w, fmt.Sprintf("Internal server error. Error: %q", err), http.StatusInternalServerError)
				return
			}
			n, err := fmt.Fprintf(w, res)
			if err != nil {
				s.logger.Printf("ERROR: helloHandler was unable to Fprintf. requestURI:'%s', from IP: [%s], send_bytes:%d\n", requestedUrl, remoteIp, n)
				http.Error(w, "Internal server error. myDefaultHandler was unable to Fprintf", http.StatusInternalServerError)
				return
			}
		default:
			s.logger.Printf("ERROR - Method not allowed. Request: %#v", r)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (s *GoHttpServer) getMyDefaultHandler() http.HandlerFunc {
	s.logger.Println("INITIAL CALL TO getMyDefaultHandler()")
	return func(w http.ResponseWriter, r *http.Request) {
		remoteIp := r.RemoteAddr
		requestedUrl := r.RequestURI
		s.logger.Printf("TRACE:[myDefaultHandler] %s  requestURI:'%s', from IP: [%s]\n", r.Method, requestedUrl, remoteIp)
		switch r.Method {
		case http.MethodGet:
			n, err := fmt.Fprintf(w, "myDefaultHandler requestURI:'%s', from IP: [%s]\n", requestedUrl, remoteIp)
			if err != nil {
				s.logger.Printf("ERROR:[myDefaultHandler] was unable to Fprintf. requestURI:'%s', from IP: [%s], send_bytes:%d\n", requestedUrl, remoteIp, n)
				http.Error(w, "Internal server error. myDefaultHandler was unable to Fprintf", http.StatusInternalServerError)
				return
			}
			s.logger.Printf("SUCCESS:[myDefaultHandler]  requestURI:'%s', from IP: [%s], send_bytes:%d\n", requestedUrl, remoteIp, n)
		default:
			s.logger.Printf("ERROR: Method not allowed. Request: %#v", r)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (s *GoHttpServer) getSlowHelloHandler() http.HandlerFunc {
	s.logger.Println("INITIAL CALL TO getSlowHelloHandler()")
	return func(w http.ResponseWriter, r *http.Request) {
		remoteIp := r.RemoteAddr
		requestedUrl := r.RequestURI
		s.logger.Printf("TRACE:[getSlowHelloHandler] %s  requestURI:'%s', from IP: [%s]\n", r.Method, requestedUrl, remoteIp)
		switch r.Method {
		case http.MethodGet:
			s.logger.Printf("INFO:[getSlowHelloHandler] sleeping %d seconds requestURI:'%s', from IP: [%s]\n", secondsToSleep, requestedUrl, remoteIp)
			time.Sleep(secondsToSleep * time.Second)
			n, err := fmt.Fprintf(w, "getSlowHelloHandler requestURI:'%s', from IP: [%s]\n", requestedUrl, remoteIp)
			if err != nil {
				s.logger.Printf("ERROR:[getSlowHelloHandler] was unable to Fprintf. requestURI:'%s', from IP: [%s], send_bytes:%d\n", requestedUrl, remoteIp, n)
				http.Error(w, "Internal server error. getSlowHelloHandler was unable to Fprintf", http.StatusInternalServerError)
				return
			}
			s.logger.Printf("SUCCESS:  requestURI:'%s', from IP: [%s], send_bytes:%d\n", requestedUrl, remoteIp, n)
		default:
			s.logger.Printf("ERROR:[getSlowHelloHandler] Method not allowed. Request: %#v", r)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func main() {
	server := NewGoHttpServer(
		defaultServerIp,
		defaultServerPort,
		log.New(os.Stdout, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile),
	)
	server.StartServer()

}
