package main

import (
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
	/*
		shutDownTimeout     = 2 * time.Second // number of second to wait before closing server
		defaultReadTimeout  = 2 * time.Minute
		defaultWriteTimeout = 2 * time.Minute
		defaultWebRootDir   = "./web/dist"
	*/
)

// main is the entry point of your todos Api Server service
func main() {
	//log := log.New(ioutil.Discard, appName, log.Ldate|log.Ltime|log.Lshortfile)
	l := log.New(os.Stdout, appName, log.Ldate|log.Ltime|log.Lshortfile)

	listenAddress, err := config.GetListenAddrFromEnv(defaultServerIp, defaultServerPort)
	if err != nil {
		log.Fatalf("error doing config.GetListenAddrFromEnv. error: %v", err)
	}
	dbDsn, err := config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort, appName, appName, defaultDBPassword)
	if err != nil {
		log.Fatalf("error doing config.GetPgDbDsnUrlFromEnv. error: %v", err)
	}
	e := todos.GetNewServer(l, dbDsn)
	e.Logger.Fatal(e.Start(listenAddress))
}
