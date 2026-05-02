package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

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
	conn, err := grpc.NewClient(
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

	// Routes legacy (gRPC) — compatibilité
	mux.HandleFunc("/api/generate", h.withCORS(h.GenerateCV))
	mux.HandleFunc("/api/stream", h.withCORS(h.StreamProgress))

	// Route v2 — direct, extended model, no gRPC overhead
	mux.HandleFunc("/api/v2/stream", h.withCORS(h.StreamExtended))

	// Route paiement — vérification Chariow + émission token
	mux.HandleFunc("/api/verify-payment", h.withCORS(h.VerifyPayment))

	// ── 4. Démarrer le gateway HTTP ───────────────────────
	// Utilise $PORT si défini (Render, Railway, etc.), sinon 8080
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	addr := ":" + httpPort
	fmt.Printf("CV Generator v2 — http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("HTTP serve: %v", err)
	}
}
