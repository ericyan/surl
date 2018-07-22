package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ericyan/surl/internal/handler"
)

func main() {
	addr := os.Getenv("SURL_ADDR")
	if addr == "" {
		addr = "0.0.0.0:3000"
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: handler.New(),
	}

	shutdown := make(chan struct{})
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		s := <-sig
		log.Printf("Signal %v received, exiting...\n", s)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Try to shutdown gracefully first. Only Interrupt and close active
		// connections when the deadline exceeded.
		err := srv.Shutdown(ctx)
		if err == context.DeadlineExceeded {
			err = srv.Close()
		}

		if err != nil {
			log.Fatal(err)
		}
		close(shutdown)
	}()

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Printf("Listening on %s...", addr)
	<-shutdown
}
