FROM golang:latest AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o cv-generator ./cmd/

# ── Image finale légère ──
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/cv-generator .
COPY --from=builder /app/fonts ./fonts
COPY --from=builder /app/web ./web

EXPOSE 8080
CMD ["./cv-generator"]