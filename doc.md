# Documentation API Marketplace & Registre ODIN

Ce document résume le fonctionnement de l'API Marketplace ODIN : connexion à la base managée Neon (PostgreSQL), architecture en couches, documentation Swagger, exposition des endpoints catalogue et du registre de packages pour le CLI, gestion d'erreurs normalisée et tests automatisés.

---

## 1. Stack Technique

* **Langage :** Go
* **Framework Web :** Gin (`[github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)`)
* **Driver PostgreSQL :** pgx v5 (`[github.com/jackc/pgx/v5/pgxpool](https://github.com/jackc/pgx/v5/pgxpool)`)
* **Documentation API :** Swaggo (`[github.com/swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)`, `[github.com/swaggo/files](https://github.com/swaggo/files)`)
* **Variables d'environnement :** godotenv (`[github.com/joho/godotenv](https://github.com/joho/godotenv)`)
* **Identifiants :** google/uuid (`[github.com/google/uuid](https://github.com/google/uuid)`)
* **Tests & Assertions :** testify (`[github.com/stretchr/testify/assert](https://github.com/stretchr/testify/assert)`)

---

## 2. Architecture & Flux de Données

Le projet respecte les principes de la Clean Architecture / Standard Go Layout :

* **Router** (`internal/router`) : achemine la requête et regroupe sous le préfixe `/api/v1`.
* **Handler** (`internal/handler`) : valide l'entrée (query/params/JSON), applique les DTOs et formate les erreurs CLI.
* **Service** (`internal/service`) : contient la logique métier, résout les versions et mappe les erreurs sentinelles.
* **Repository** (`internal/repository`) : exécute les requêtes SQL paramétrées (`pgxpool`).
* **Base Neon** : PostgreSQL managé.

### Rôle de chaque dossier :

* `config/` : chargement et validation des variables d'environnement (`.env`).
* `db/` : initialisation et gestion du pool de connexions PostgreSQL (`pgxpool.Pool`).
* `docs/` : spécifications OpenAPI générées par `swag` (`docs.go`, `swagger.json`, `swagger.yaml`).
* `internal/dto/` : structures d'entrée/sortie API, réponses standardisées et format d'erreur machine-readable (`CLIErrorResponse`).
* `internal/model/` : entités métier reflétant les tables (`Agent`, `AgentVersion`).
* `internal/repository/` : requêtes SQL directes (recherche partielle `ILIKE`, résolution des versions).
* `internal/service/` : règles métier et erreurs sentinelles (`ErrPackageNotFound`, `ErrVersionNotFound`).
* `internal/handler/` : contrôleurs HTTP (`CatalogHandler`, `PackageHandler`, `HealthHandler`).
* `internal/router/` : déclaration et organisation modulaire des routes (`health`, `catalog`, `packages`).
* `cmd/api/main.go` : injection des dépendances, configuration Swagger et démarrage du serveur.

---

## 3. Endpoints Implémentés

### Système & Santé

| Méthode | Route | Description | Statuts HTTP |
| --- | --- | --- | --- |
| `GET` | `/swagger/*any` | Interface interactive Swagger UI | 200 |
| `GET` | `/health/live` | Liveness probe (le serveur répond) | 200 |
| `GET` | `/health/ready` | Readiness probe (la base Neon est joignable) | 200, 503 |

### Catalogue Web (`/api/v1/catalog`)

| Méthode | Route | Description | Statuts HTTP |
| --- | --- | --- | --- |
| `GET` | `/api/v1/catalog/agents` | Liste tous les agents enregistrés | 200, 500 |
| `GET` | `/api/v1/catalog/agents/{id}` | Détail d'un agent par son UUID | 200, 400, 404, 500 |
| `POST` | `/api/v1/catalog/agents` | Crée un nouvel agent via son DTO | 201, 400, 500 |

### Registre de Packages CLI (`/api/v1/packages`)

| Méthode | Route | Description | Statuts HTTP |
| --- | --- | --- | --- |
| `GET` | `/api/v1/packages/search?q=<query>` | Recherche insensible à la casse sur nom et description | 200, 500 |
| `GET` | `/api/v1/packages/{name}/{version}` | Métadonnées du package et de sa version (supporte `latest`) | 200, 404, 500 |

---

## 4. Contrat d'Erreurs Spécifique CLI

Les routes sous `/api/v1/packages` renvoient une charge utile structurée afin de permettre au CLI ODIN de gérer les sorties de commandes et les codes de retour système.

Codes machines retournés :

* `PACKAGE_NOT_FOUND` : le nom du package n'existe pas en base (404).
* `VERSION_NOT_FOUND` : le package existe, mais la version spécifiée est introuvable (404).
* `INTERNAL_SERVER_ERROR` : erreur inattendue ou indisponibilité base de données (500).

---

## 5. Configuration Requise (`.env`)

* `PORT=8081` (configuré par défaut pour éviter tout conflit avec le service `auth` du projet ODIN)
* `DATABASE_URL=postgres://<user>:<password>@<endpoint>.neon.tech/<dbname>?sslmode=require`

---

## 6. Démarrage Local & Tests

* Téléchargement et synchronisation des modules : `go mod tidy`
* Régénération des définitions Swagger : `swag init -g cmd/api/main.go`
* Lancement des tests automatisés : `go test -v ./...`
* Lancement du serveur API : `go run cmd/api/main.go`
* Swagger UI accessible sur : `