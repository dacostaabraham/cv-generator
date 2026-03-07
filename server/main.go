package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/dacostaabraham/cv-generator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":50051"

func main() {
	// 1. Creer un listener TCP sur le port 50051
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Impossible d'ecouter sur %s: %v", port, err)
	}

	// 2. Creer le serveur gRPC
	grpcServer := grpc.NewServer(
	// On peut ajouter des intercepteurs ici plus tard (auth, logs...)
	)

	// 3. Enregistrer notre service CVService
	pb.RegisterCVGeneratorServer(grpcServer, &CVService{})

	// 4. Activer la reflection gRPC
	// Permet a grpcurl de decouvrir les services sans avoir le .proto
	reflection.Register(grpcServer)

	// 5. Demarrer le serveur
	fmt.Printf("Serveur gRPC demarre sur le port %s\n", port)
	fmt.Println("En attente de connexions...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Erreur serveur: %v", err)
	}
}
