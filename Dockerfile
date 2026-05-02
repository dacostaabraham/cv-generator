# ─────────────────────────────────────────────
# Stage 1 — Build
# ─────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copier les dépendances en premier (cache Docker optimal)
COPY go.mod go.sum ./
RUN go mod download

# Copier le code source
COPY . .

# Compiler — binaire statique, stripped (taille minimale)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o cv-generator ./cmd/.

# ─────────────────────────────────────────────
# Stage 2 — Image finale légère
# ─────────────────────────────────────────────
FROM debian:bookworm-slim

WORKDIR /app

# Certificats TLS
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Binaire + assets
COPY --from=builder /app/cv-generator .
COPY --from=builder /app/fonts ./fonts
COPY --from=builder /app/web ./web

# Render injecte $PORT dynamiquement — on documente 8080 comme défaut local
EXPOSE 8080

CMD ["./cv-generator"]