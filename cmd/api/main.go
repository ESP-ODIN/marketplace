package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"marketplace/config"
	"marketplace/db"
	"marketplace/internal/handler"
	"marketplace/internal/repository"
	"marketplace/internal/router"
	"marketplace/internal/service"
)

// @title           Marketplace API
// @version         1.0
// @description     API backend pour la marketplace Odin.
// @host            localhost:8081
// @BasePath        /api/v1
func main() {
	if err := run(); err != nil {
		log.Fatalf("Erreur API : %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("erreur de configuration : %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := db.Connect(connectCtx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("erreur de connexion : %w", err)
	}
	defer pool.Close()

	agentRepo := repository.NewAgentRepository(pool)
	agentService := service.NewAgentService(agentRepo)

	r := router.Setup(router.Config{
		HealthHandler:  handler.NewHealthHandler(pool),
		AgentHandler:   handler.NewAgentHandler(agentService),
		PackageHandler: handler.NewPackageHandler(agentService),
	})

	addr := ":" + cfg.Port

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("erreur serveur : %w", err)
	}
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(listener) }()
	fmt.Printf("🚀 Serveur démarré sur http://localhost%s\n", addr)

	select {
	case err := <-serveErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("erreur serveur : %w", err)
	case <-ctx.Done():
		stop()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			_ = srv.Close()
			return fmt.Errorf("erreur lors de l'arrêt du serveur : %w", err)
		}
		log.Print("Serveur arrêté")
		return nil
	}
}
