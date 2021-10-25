package main

import (
	"github.com/lao-tseu-is-alive/go-cloud-learning-01-http/internal/todos"
	"github.com/lao-tseu-is-alive/go-cloud-learning-01-http/pkg/config"
	"time"
)

const (
	appName             = "basicUpload"
	defaultServerPort   = 8080
	defaultServerIp     = "127.0.0.1"
	shutDownTimeout     = 2 * time.Second // number of second to wait before closing server
	defaultReadTimeout  = 2 * time.Minute
	defaultWriteTimeout = 2 * time.Minute
	defaultWebRootDir   = "./web/dist"
)

// main is the entry point of your basicUpload service
// you can try it with :
// SERVERIP=127.0.0.1 PORT=3333 make go-run
func main() {
	listenAddress, err := config.GetListenAddrFromEnv(defaultServerIp, defaultServerPort)
	if err != nil {
		panic(err)
	}
	e := todos.GetNewServer(false)
	e.Logger.Fatal(e.Start(listenAddress))
}
