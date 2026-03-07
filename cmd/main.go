package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/dacostaabraham/cv-generator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// ── 1. Démarrer le serveur gRPC en goroutine ──────────
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("gRPC listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterCVGeneratorServer(grpcServer, &CVService{})
		reflection.Register(grpcServer)
		fmt.Println("gRPC server :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC serve: %v", err)
		}
	}()

	// ── 2. Connecter le gateway au gRPC interne ───────────
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("gRPC dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewCVGeneratorClient(conn)
	h := &Handler{client: client}

	// ── 3. Routes HTTP ────────────────────────────────────
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/api/generate", h.withCORS(h.GenerateCV))
	mux.HandleFunc("/api/stream", h.withCORS(h.StreamProgress))

	// ── 4. Démarrer le gateway HTTP ───────────────────────
	port := ":8080"
	fmt.Printf("Gateway HTTP http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("HTTP serve: %v", err)
	}
}
