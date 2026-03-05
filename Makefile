# Variables
PROTO_DIR := proto
OUT_DIR   := .

# Commandes disponibles
.PHONY: proto build run clean help

## help: Affiche cette aide
help:
	@echo "Commandes disponibles:"
	@grep -h '##' $(MAKEFILE_LIST) | sed 's/## //'

## proto: Génère le code Go depuis les fichiers .proto
proto:
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/*.proto

## build: Compile le projet
build:
	go build ./...

## run-server: Lance le serveur gRPC
run-server:
	go run ./server/main.go

## run-gateway: Lance le gateway HTTP
run-gateway:
	go run ./gateway/main.go

## deps: Installe les dépendances
deps:
	go get google.golang.org/grpc
	go get google.golang.org/protobuf
	go get github.com/signintech/gopdf

## clean: Supprime les fichiers générés
clean:
	rm -f proto/*.pb.go