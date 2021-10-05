package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	const (
		secondsToSleep         = 5
		secondsShutDownTimeout = 15 // number of second to wait before closing server
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("New request: '%s' from RemoteAddr %s, sleeping %d seconds\n", r.RequestURI, r.RemoteAddr, secondsToSleep)
		time.Sleep(secondsToSleep * time.Second)
		fmt.Fprintf(w, "Hello world after sleeping for %d  seconds!\n", secondsToSleep)
	})

	srv := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error when trying ListenAndServe: %s\n", err)
		}
	}()

	log.Println("Server started and listening on : http://localhost" + srv.Addr)
	log.Printf("From another terminal run : curl http://localhost%s/username=titi ", srv.Addr)
	log.Printf("And just after run CTRL+C on this server terminal, to experiment  a gracefull shutdown of running request")

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		secondsShutDownTimeout*time.Second)
	// gracefully shuts down the server without interrupting any active connections
	// https://pkg.go.dev/net/http#Server.Shutdown
	srv.Shutdown(ctx)
	<-ctx.Done()
	cancel()
	log.Println("Server gracefully stopped")
}
