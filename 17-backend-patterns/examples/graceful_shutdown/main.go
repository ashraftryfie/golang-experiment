package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Root context cancelled on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[Server] Processing long request...")
		time.Sleep(1 * time.Second)
		_, _ = w.Write([]byte("Work complete!\n"))
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8088",
		Handler: mux,
	}

	// Run server in background goroutine
	go func() {
		fmt.Printf("[Server] Listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[Server] Listen error: %v\n", err)
		}
	}()

	// Wait for OS interrupt signal
	<-ctx.Done()
	fmt.Println("\n[Server] Shutdown signal received, draining active connections...")

	// 5-second graceful shutdown timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("[Server] Forced shutdown due to error: %v\n", err)
	} else {
		fmt.Println("[Server] Server stopped cleanly. All requests completed.")
	}
}
