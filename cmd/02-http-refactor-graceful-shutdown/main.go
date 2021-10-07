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
	defaultServerPath   = "/"
	shutDownTimeout     = 15 * time.Second // number of second to wait before closing server
	defaultReadTimeout  = 2 * time.Minute  // max time to read request from the client
	defaultWriteTimeout = 2 * time.Minute  // max time to write response to the client
	defaultIdleTimeout  = 2 * time.Minute  // max time for connections using TCP Keep-Alive
	defaultUserName     = "𝕎𝕠𝕣𝕝𝕕"
	defaultMessage      = "🅆🄴🄻🄲🄾🄼🄴 🄷🄾🄼🄴 🏠"
	defaultNotFound     = "🤔 ℍ𝕞𝕞... 𝕤𝕠𝕣𝕣𝕪 :【𝟜𝟘𝟜 : ℙ𝕒𝕘𝕖 ℕ𝕠𝕥 𝔽𝕠𝕦𝕟𝕕】🕳️ 🔥"
	secondsToSleep      = 5
	htmlHeaderStart     = `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css"/>`
)

//getListenAddrFromEnv returns a valid TCP/IP listening address string based on the values of ENV SERVERIP:PORT
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

func getHtmlHeader(title string) string {
	return fmt.Sprintf("%s<title>%s</title></head>", htmlHeaderStart, title)
}

func getHtmlPage(title string) string {
	return getHtmlHeader(title) +
		fmt.Sprintf("\n<body><div class=\"container\"><h3>%s</h3></div></body></html>", title)
}

func getHelloMsg(name string) (string, error) {
	const helloMsg = htmlHeaderStart +
		`<title>𝙂𝙤 𝙃𝙚𝙡𝙡𝙤 {{.UserName}} 👋 🖖 🫂</title></head>
<body><div class="container"><h3>ℍ𝕖𝕝𝕝𝕠, {{.UserName}} 👋 🖖 🫂 </h3></div></body></html>`

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

//GoHttpServer is a struct type to store information related to all handlers of web server
type GoHttpServer struct {
	listenAddress string
	// later we will store here the connection to database
	//DB  *db.Conn
	logger *log.Logger
	router *http.ServeMux
}

//NewGoHttpServer is a constructor that initializes the server mux (routes) and all fields of the  GoHttpServer type
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

// (*GoHttpServer) routes initializes all the handlers paths of this web server, it is called inside the NewGoHttpServer constructor
func (s *GoHttpServer) routes() {
	s.router.Handle("/", s.getMyDefaultHandler())
	s.router.Handle("/hello", s.getHelloHandler())
	s.router.Handle("/slowHello", s.getSlowHelloHandler())
}

// StartServer initializes all the handlers paths of this web server, it is called inside the NewGoHttpServer constructor
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
				s.logger.Printf("ERROR:[HelloHandler] Bad request. Error in ParseQuery. Error: %q . Request was: \n%#v\n", err, r)
				http.Error(w, fmt.Sprintf("Bad request. Error in ParseQuery: %q", err), http.StatusBadRequest)
				return
			}
			username := defaultUserName
			if query.Has("username") {
				username = query.Get("username")
			}
			username = strings.TrimSpace(username)
			if len(username) == 0 {
				s.logger.Printf("ERROR:[HelloHandler] Bad request. Username cannot be empty. Request was: \n%#v\n", r)
				http.Error(w, fmt.Sprintf("Bad request. In query.Get('username'): username cannot be empty or spaces only"), http.StatusBadRequest)
				return
			}
			res, err := getHelloMsg(username)
			if err != nil {
				s.logger.Printf("ERROR:[HelloHandler] Internal server error. Request was: \n%#v\n", r)
				http.Error(w, fmt.Sprintf("Internal server error. Error: %q", err), http.StatusInternalServerError)
				return
			}
			n, err := fmt.Fprintf(w, res)
			if err != nil {
				s.logger.Printf("ERROR:[HelloHandler] was unable to Fprintf. requestURI:'%s', from IP: [%s], send_bytes:%d\n", requestedUrl, remoteIp, n)
				http.Error(w, "Internal server error. myDefaultHandler was unable to Fprintf", http.StatusInternalServerError)
				return
			}
			s.logger.Printf("SUCCESS:[HelloHandler]  requestURI:'%s', from IP: [%s], send_bytes:%d\n", requestedUrl, remoteIp, n)
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
		requestedUrlPath := r.URL.Path
		s.logger.Printf("TRACE:[myDefaultHandler] %s  path:'%s', from IP: [%s]\n", r.Method, requestedUrlPath, remoteIp)
		switch r.Method {
		case http.MethodGet:
			if len(strings.TrimSpace(requestedUrlPath)) == 0 || requestedUrlPath == defaultServerPath {
				log.Printf("DEBUG getMyDefaultHandler requested Path %#v\n", requestedUrlPath)
				n, err := fmt.Fprintf(w, getHtmlPage(defaultMessage))
				if err != nil {
					s.logger.Printf("ERROR:[myDefaultHandler] was unable to Fprintf. path:'%s', from IP: [%s], send_bytes:%d\n", requestedUrlPath, remoteIp, n)
					http.Error(w, "Internal server error. myDefaultHandler was unable to Fprintf", http.StatusInternalServerError)
					return
				}
				s.logger.Printf("SUCCESS:[myDefaultHandler]  path:'%s', from IP: [%s], send_bytes:%d\n", requestedUrlPath, remoteIp, n)
			} else {
				w.WriteHeader(http.StatusNotFound)
				n, err := fmt.Fprintf(w, getHtmlPage(defaultNotFound))
				if err != nil {
					s.logger.Printf("ERROR:[myDefaultHandler] Not Found was unable to Fprintf. path:'%s', from IP: [%s], send_bytes:%d\n", requestedUrlPath, remoteIp, n)
					http.Error(w, "Internal server error. myDefaultHandler was unable to Fprintf", http.StatusInternalServerError)
					return
				}
			}
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

	/* example code to log inside a file instead of stdout
	const logFilename = "server.log"
	f, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("FATAL ERROR : Unable to open log file : %s for writing. Error : %s", logFilename, err)
	}
	defer f.Close()
	l := log.New(os.Stdout, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	*/

	// ioutil.Discard(https://golang.org/pkg/io/ioutil/#pkg-variables)
	// ioutil.Discard is a writer on which all calls succeed without doing anything.
	//l := log.New(ioutil.Discard, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	// you can also disable output for the logger anytime by doing a simple
	//l.SetOutput(ioutil.Discard)

	l := log.New(os.Stdout, "hello-api_", log.Ldate|log.Ltime|log.Lshortfile)
	server := NewGoHttpServer(defaultServerIp, defaultServerPort, l)
	server.StartServer()

}
