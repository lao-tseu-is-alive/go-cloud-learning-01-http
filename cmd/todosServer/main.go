package main

import (
	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-learning-01-http/internal/todos"
	"github.com/lao-tseu-is-alive/go-cloud-learning-01-http/pkg/config"
	"log"
	"os"
)

const (
	appName           = "todos"
	defaultServerPort = 8080
	defaultServerIp   = "127.0.0.1"
	defaultDBPort     = 5432
	defaultDBIp       = "127.0.0.1"
	defaultDBPassword = "todos_password"
	/*  defaultDBSslMode: in dev env it can be ok to disable SSL mode but in prod it is another story
			it depends on various factor. is your service (go) running in the same host as the db (localhost ?)
	 		if not, is the network between your server and your db trusted ?? read the doc and ask your security officer:
	 		https://www.postgresql.org/docs/11/libpq-ssl.html#LIBPQ-SSL-PROTECTION
	*/
	defaultDBSslMode = "disable"
	defaultDBDriver  = "postgres"
	/*
		shutDownTimeout     = 2 * time.Second // number of second to wait before closing server
		defaultReadTimeout  = 2 * time.Minute
		defaultWriteTimeout = 2 * time.Minute
		defaultWebRootDir   = "./web/dist"
	*/
)

// GetNewServer initialize a new Echo server and returns it
func GetNewServer(l *log.Logger, store todos.Storage) *echo.Echo {
	e := echo.New()
	myTodosApi := todos.Service{
		Log:   l,
		Store: store,
	}

	todos.RegisterHandlers(e, &myTodosApi)
	// add a route for maxId
	e.GET("/todos/maxid", myTodosApi.GetMaxId)
	return e
}

// main is the entry point of your todos Api TodosService service
func main() {
	//l := log.New(ioutil.Discard, appName, 0)
	l := log.New(os.Stdout, appName, log.Ldate|log.Ltime|log.Lshortfile)

	listenAddress, err := config.GetListenAddrFromEnv(defaultServerIp, defaultServerPort)
	if err != nil {
		log.Fatalf("error doing config.GetListenAddrFromEnv. error: %v", err)
	}

	driver := config.GetDbDriverFromEnv(defaultDBDriver)
	if !todos.IsDriverSupported(driver) {
		log.Fatalf("error the driver : %s is not supported yet.\n", driver)
	}

	var dbDsn = ""
	if driver == "postgres" {
		dbDsn, err = config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
			appName, appName, defaultDBPassword, defaultDBSslMode)
		if err != nil {
			log.Fatalf("error doing config.GetPgDbDsnUrlFromEnv. error: %v", err)
		}
	}

	s, err := todos.GetStorageInstance(driver, dbDsn, l)
	if err != nil {
		l.Fatalf("error getting Storage Instance for driver %s. error: %v", driver, err)
	}
	defer s.Close()

	e := GetNewServer(l, s)

	e.Logger.Fatal(e.Start(listenAddress))
}
