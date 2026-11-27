# Dockerfile de développement pour le service API (Go)
# Objectif : lancer l'API avec hot-reload (via air) dans un conteneur,
# pour être branché tel quel dans le docker-compose global de l'archi micro-service.

FROM golang:1.27-alpine

WORKDIR /app

# air : hot-reload pour Go (rebuild + restart au changement de fichier),
# l'équivalent de `next dev` côté Node.
RUN go install github.com/air-verse/air@v1.67.4

# L'utilisateur de développement doit pouvoir compiler et écrire les caches Go.
RUN addgroup -g 10001 app \
    && adduser -D -u 10001 -G app app \
    && mkdir -p /app/tmp /home/app/.cache/go-build \
    && chown -R app:app /app /go/pkg /home/app

ENV GOCACHE=/home/app/.cache/go-build
USER app:app

# Copier les manifests avant le code permet de réutiliser le cache Docker
# tant que go.mod et go.sum ne changent pas.
COPY --chown=app:app go.mod go.sum ./
RUN go mod download

# Le reste du code sera de toute façon écrasé par le bind mount en dev,
# mais on le copie pour que l'image soit utilisable seule (ex: premier build, CI).
COPY --chown=app:app . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
