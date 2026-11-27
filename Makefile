.PHONY: run build test vet fmt docker-build docker-dev swagger
-include .env
PORT ?= 8080

run:
	go run ./cmd/api

build:
	go build -o bin/ ./cmd/...

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

docker-build:
	docker build -t marketplace-api .

# Lance l'API en local avec hot-reload (bind mount du code) et l'expose sur le port PORT (défini dans .env).
docker-dev:
	docker run --rm -it \
		-p $(PORT):8080 \
		--env-file .env \
		-e PORT=8080 \
		-v "$(CURDIR)":/app \
		-v /app/tmp \
		marketplace-api

swagger:
	swag init -g cmd/api/main.go
