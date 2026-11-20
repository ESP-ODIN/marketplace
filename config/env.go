package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	Port         string
	AuthJWKSURL  string
	AuthIssuer   string
	AuthAudience string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("variable d'environnement DATABASE_URL manquante")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwksURL := os.Getenv("AUTH_JWKS_URL")
	if jwksURL == "" {
		return nil, fmt.Errorf("variable d'environnement AUTH_JWKS_URL manquante")
	}

	issuer := os.Getenv("AUTH_ISSUER")
	if issuer == "" {
		return nil, fmt.Errorf("variable d'environnement AUTH_ISSUER manquante")
	}

	audience := os.Getenv("AUTH_AUDIENCE")
	if audience == "" {
		audience = "odin-api"
	}

	return &Config{
		DatabaseURL:  dbURL,
		Port:         port,
		AuthJWKSURL:  jwksURL,
		AuthIssuer:   issuer,
		AuthAudience: audience,
	}, nil
}
