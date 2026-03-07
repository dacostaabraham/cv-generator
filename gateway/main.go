package main

import (
	"fmt"
	"log"
	"net/http"

	pb "github.com/dacostaabraham/cv-generator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Connexion au serveur gRPC (localhost:50051)
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Connexion gRPC impossible: %v", err)
	}
	defer conn.Close()

	// 2. Creer le client gRPC (genere par protoc)
	client := pb.NewCVGeneratorClient(conn)

	// 3. Creer le handler avec le client injecte
	h := &Handler{client: client}

	// 4. Definir les routes HTTP
	mux := http.NewServeMux()

	// Servir les fichiers statiques (web/)
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	// Route API : generer un CV
	mux.HandleFunc("/api/generate", h.withCORS(h.GenerateCV))

	// Route API : streaming progression
	mux.HandleFunc("/api/stream", h.withCORS(h.StreamProgress))

	// 5. Demarrer le serveur HTTP
	fmt.Println("Gateway HTTP demarre sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Erreur serveur: %v", err)
	}
}
