package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/config"
	delivery "github.com/nachsyas/arham-porto/apps/api/internal/delivery/http"
	"github.com/nachsyas/arham-porto/apps/api/internal/usecase"
)

func main() {
	cfg := config.Load()

	// Initialize stub use cases for Phase 0
	profileUC := usecase.NewProfileUseCase(nil)
	projectUC := usecase.NewProjectUseCase(nil)
	evidenceUC := usecase.NewEvidenceUseCase(nil)

	// Initialize delivery layer with standard net/http router
	handler := delivery.NewHandler(profileUC, projectUC, evidenceUC)
	router := delivery.NewRouter(handler, cfg.WebOrigin)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[arham-porto-api] Server starting on :%s (env: %s)", cfg.Port, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server ListenAndServe error: %v", err)
		}
	}()

	<-stop
	log.Println("[arham-porto-api] Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}

	log.Println("[arham-porto-api] Server stopped cleanly")
}
