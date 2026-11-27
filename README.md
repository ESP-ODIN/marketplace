# Marketplace API

Go API for the Odin marketplace, using Gin for HTTP routing and PostgreSQL for storage.

## Prerequisites

- **Go 1.27.1 or later**, as specified in `go.mod`, for local development.
- **PostgreSQL** accessible through `DATABASE_URL`, with the application schema initialized.
- **Make** to use the provided development commands.
- **Docker** with Linux containers for the Docker workflow. Go does not need to be installed on the host for this workflow.

## Environment setup

Copy the example configuration at the project root:

```sh
cp .env.exemple .env
```

In PowerShell, use `Copy-Item .env.exemple .env`.

Set `DATABASE_URL` to your PostgreSQL connection URL. `PORT` controls the local HTTP port:
the example sets it to `8081`, while the application defaults to `8080` when it is unset.

## Local development

```sh
make run
```

With the example port, check liveness at `http://localhost:8081/health/live`, database readiness
at `http://localhost:8081/health/ready`, and API documentation at `http://localhost:8081/swagger/index.html`.
The agent catalog is available at `GET /api/v1/catalog/agents`.

Use `make fmt` to format all Go source files, including `internal/`, `make test` to run tests,
and `make vet` for static analysis.

## Docker development

The development image runs the API with Air for automatic rebuilds and restarts when Go files change.
It runs as the non-root user `app` (UID/GID `10001`).

Configure `.env` first, then build and start the container:

```sh
make docker-build
make docker-dev
```

The container always listens on port `8080`: `docker-dev` passes `-e PORT=8080`, overriding the
value supplied by `--env-file .env`. The host port comes from `PORT` in `.env` (`8081` in the example,
or `8080` if unset). With the example configuration, Docker publishes `8081:8080`.

```sh
curl http://localhost:8081/health/live
curl http://localhost:8081/health/ready
```

The source directory is mounted at `/app`, and a separate volume stores build output in `/app/tmp`.
On Linux, the mounted source must be readable by UID `10001` and its directories must be traversable.

`DATABASE_URL` must be reachable from inside the container. With Docker Desktop, use
`host.docker.internal` instead of `localhost` when PostgreSQL runs on the host. For PostgreSQL
in another container, use its service name and connect both containers to the same Docker network.
The `docker-dev` target does not create a database or join a custom network automatically.
