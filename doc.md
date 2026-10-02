# Initialisation API Marketplace & Connexion Neon

Ce document résume le travail effectué sur l'API Marketplace : connexion de l'API Go à la base de données managée Neon (PostgreSQL), mise en place de l'architecture en couches et intégration de la documentation interactive Swagger.

---

## 1. Stack Technique

- **Langage :** Go
- **Framework Web :** Gin (`[github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)`)
- **Driver PostgreSQL :** pgx v5 (`[github.com/jackc/pgx/v5/pgxpool](https://github.com/jackc/pgx/v5/pgxpool)`)
- **Documentation API :** Swaggo (`[github.com/swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)`, `[github.com/swaggo/files](https://github.com/swaggo/files)`)
- **Variables d'environnement :** godotenv (`[github.com/joho/godotenv](https://github.com/joho/godotenv)`)
- **Identifiants :** google/uuid (`[github.com/google/uuid](https://github.com/google/uuid)`)

---

## 2. Architecture & Flux de Données

Le projet respecte les principes de la Clean Architecture / Standard Go Layout :

```text
Requête HTTP  
   │  
   ▼  
 [Router] (internal/router)  
   │  Achemine la requête et regroupe sous le préfixe /api/v1  
   ▼  
 [Handler] (internal/handler)  
   │  Valide le format d'entrée (JSON / paramètres URL) via les DTOs  
   ▼  
 [Service] (internal/service)  
   │  Contient la logique métier et convertit les DTOs en Modèles  
   ▼  
 [Repository] (internal/repository)  
   │  Exécute les requêtes SQL (pgxpool)  
   ▼  
 [Base Neon] (PostgreSQL managé)
```

### Rôle de chaque dossier :

- `config/` : chargement et validation des variables d'environnement (`.env`).
- `db/` : initialisation et gestion du pool de connexions PostgreSQL (`pgxpool.Pool`).
- `docs/` : spécifications OpenAPI générées par `swag` (`docs.go`, `swagger.json`, `swagger.yaml`).
- `internal/dto/` : structures d'entrée/sortie de l'API avec tags de validation Gin (`binding:"required"`).
- `internal/model/` : entités métier reflétant les tables en base de données.
- `internal/repository/` : couche d'accès aux données (requêtes SQL directes).
- `internal/service/` : couche de logique applicative et cas d'usage.
- `internal/handler/` : contrôleurs HTTP (traitement de `gin.Context`, codes de statut et réponses JSON).
- `internal/router/` : déclaration et organisation modulaire des routes (`health`, `catalog`).
- `cmd/api/main.go` : point d'entrée de l'application (injection des dépendances, annotations Swagger et démarrage).

---

## 3. Endpoints Implémentés

| Méthode | Route | Description | Statuts HTTP |
| :--- | :--- | :--- | :--- |
| `GET` | `/swagger/*any` | Interface interactive Swagger UI | 200 |
| `GET` | `/health/live` | Liveness probe (le serveur répond) | 200 |
| `GET` | `/health/ready` | Readiness probe (la base Neon est joignable) | 200, 503 |
| `GET` | `/api/v1/catalog/agents` | Récupère la liste de tous les agents | 200, 500 |
| `GET` | `/api/v1/catalog/agents/{id}` | Récupère le détail d'un agent par son UUID | 200, 400, 404, 500 |
| `POST` | `/api/v1/catalog/agents` | Enregistre un nouvel agent via son DTO | 201, 400, 500 |

---

## 4. Configuration Requise (`.env`)

```env
PORT=8081
DATABASE_URL=postgres://<user>:<password>@<endpoint>.neon.tech/<dbname>?sslmode=require
```

> **Note :** le port `8081` est configuré par défaut pour éviter tout conflit avec le service `auth` du projet ODIN.

---

## 5. Démarrage Local

```bash
# Téléchargement et synchronisation des modules Go
go mod tidy

# Régénération des définitions Swagger
swag init -g cmd/api/main.go

# Lancement du serveur API
go run cmd/api/main.go
```

L'interface Swagger UI est accessible sur : `http://localhost:8081/swagger/index.html`