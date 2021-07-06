package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/autocorrectoff/banners/internal/dto"
	"github.com/autocorrectoff/banners/internal/engine"
	"github.com/autocorrectoff/banners/internal/handler"
	"github.com/autocorrectoff/banners/internal/repo"
	"github.com/minus5/svckit/log"
	"github.com/minus5/svckit/nsq"
)

const (
	port = ":8000"
)

func main() {
	fmt.Println(os.Getpid())
	rpo, err := repo.New()
	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(rpo)
	n := engine.New(rpo)

	sub := nsq.Sub("ponuda.req", func(m *nsq.Message) error {
		var err error
		req := &dto.MatchRequest{}
		mb := m.Body
		if err = json.Unmarshal(mb, req); err != nil {
			log.S("ponuda.req", string(mb)).Error(err)
			return nil
		}
		matchErr := n.HandleMatch(req)
		if matchErr != nil {
			log.Error(err)
			return err
		}

		fmt.Printf("%+v\n", req)
		return nil
	})
	defer sub.Close()

	sm := http.NewServeMux()
	sm.Handle("/games/details", h)

	s := http.Server{
		Addr:         port,
		Handler:      sm,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("Starting server on port %s", port)

		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	sig := <-c
	log.Info("Got signal: %s", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	s.Shutdown(ctx)
}