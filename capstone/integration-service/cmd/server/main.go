package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/adapters/external"
	httpadapter "github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/adapters/http"
	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/adapters/memory"
	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/service"
)

func main() {
	// 1. Structured Logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("initializing integration service components")

	// 2. Adapters
	storage := memory.NewInMemoryJobStorage()
	queue := memory.NewChannelJobQueue(200)
	extClient := external.NewResilientExternalClient()

	// 3. Service Layer
	svc := service.NewIntegrationService(storage, queue)

	// 4. Worker Pool
	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	workerPool := service.NewWorkerPool(4, storage, queue, extClient, 3*time.Second)
	workerPool.Start(workerCtx)
	slog.Info("worker pool started with 4 concurrent routines")

	// 5. HTTP Adapter
	httpHandler := httpadapter.NewServerHandler(svc)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      httpHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("HTTP server listening", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// 6. Graceful Shutdown Signal Handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		slog.Error("server startup failed", slog.String("error", err.Error()))
		os.Exit(1)
	case sig := <-shutdown:
		slog.Info("shutdown signal received", slog.String("signal", sig.String()))

		// Step A: Shutdown HTTP listener
		shutdownCtx, cancelHTTP := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelHTTP()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("graceful HTTP shutdown failed", slog.String("error", err.Error()))
			_ = server.Close()
		}

		// Step B: Signal workers and drain queue
		cancelWorkers()
		_ = queue.Close()
		workerPool.Stop()

		slog.Info("integration service stopped gracefully")
	}
}
