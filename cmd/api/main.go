package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"marketplace/config"
	"marketplace/db"
	"marketplace/internal/handler"
	"marketplace/internal/repository"
	"marketplace/internal/router"
	"marketplace/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erreur de configuration : %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erreur de connexion : %v", err)
	}
	defer pool.Close()

	agentRepo := repository.NewAgentRepository(pool)
	agentService := service.NewAgentService(agentRepo)

	r := router.Setup(router.Config{
		HealthHandler: handler.NewHealthHandler(pool),
		AgentHandler:  handler.NewAgentHandler(agentService),
	})

	addr := ":8080"
	if cfg.Port != "" {
		addr = ":" + cfg.Port
	}

	fmt.Printf("🚀 Serveur démarré sur http://localhost%s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Erreur serveur : %v", err)
	}
}
