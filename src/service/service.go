package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartHttpService(port int, handler http.Handler) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
		}
		serverStopCtx()
	}()

	// Run the server

	_, err := os.Stat("./localhost.key")
	if err != nil {
		fmt.Println(fmt.Sprintf("Starting server on http://127.0.0.1:%d", port))
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
		}
	} else {
		fmt.Println(fmt.Sprintf("Starting server on https://127.0.0.1:%d", port))
		err := server.ListenAndServeTLS("localhost.crt", "localhost.key")
		if err != nil && err != http.ErrServerClosed {
		}
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
