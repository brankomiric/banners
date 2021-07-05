package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/autocorrectoff/banners/internal/engine"
)

const (
	port = ":8000"
)

func main() {
	fmt.Println(os.Getpid())

	h := engine.New()

	sm := http.NewServeMux()
	sm.Handle("/games/details/:id", h)

	s := http.Server{
		Addr:         port,
		Handler:      sm,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		// log.Info("Starting server on port %s", port)
		fmt.Printf("Starting server on port %s", port)

		err := s.ListenAndServe()
		if err != nil {
			// log.Fatal(err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	// log.Info("Got signal: %s", sig)
	fmt.Printf("Got signal: %s", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}