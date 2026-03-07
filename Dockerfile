FROM golang:latest AS builder

WORKDIR /app

# Copier SEULEMENT go.mod et go.sum d'abord
COPY go.mod go.sum ./

# Télécharger les dépendances (sans tidy)
RUN go mod download

# Puis copier le reste du code
COPY . .

# Builder
RUN CGO_ENABLED=0 GOOS=linux go build -o cv-generator ./cmd/

# ── Image finale légère ──
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/cv-generator .
COPY --from=builder /app/fonts ./fonts
COPY --from=builder /app/web ./web

EXPOSE 8080
CMD ["./cv-generator"]