package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pb "github.com/dacostaabraham/cv-generator/proto"
)

// Handler contient le client gRPC
type Handler struct {
	client pb.CVGeneratorClient
}

// CVRequestJSON : le JSON recu depuis le browser
type CVRequestJSON struct {
	FullName    string           `json:"full_name"`
	Email       string           `json:"email"`
	Phone       string           `json:"phone"`
	Location    string           `json:"location"`
	Summary     string           `json:"summary"`
	Skills      []string         `json:"skills"`
	Experiences []ExperienceJSON `json:"experiences"`
	Education   []EducationJSON  `json:"education"`
	Template    string           `json:"template"`
	Photo       string           `json:"photo"` // base64
}

type ExperienceJSON struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type EducationJSON struct {
	School string `json:"school"`
	Degree string `json:"degree"`
	Year   int32  `json:"year"`
}

// toProto convertit le JSON en message Protobuf
func (r CVRequestJSON) toProto() (*pb.CVRequest, error) {
	req := &pb.CVRequest{
		FullName: r.FullName,
		Email:    r.Email,
		Phone:    r.Phone,
		Location: r.Location,
		Summary:  r.Summary,
		Skills:   r.Skills,
		Template: r.Template,
	}

	// Decoder la photo base64
	if r.Photo != "" {
		photoBytes, err := base64.StdEncoding.DecodeString(r.Photo)
		if err != nil {
			// Photo invalide : on continue sans photo
			fmt.Printf("[warn] photo base64 invalide: %v\n", err)
		} else {
			req.Photo = photoBytes
		}
	}

	// Convertir les experiences
	for _, e := range r.Experiences {
		req.Experiences = append(req.Experiences, &pb.Experience{
			Company:     e.Company,
			Position:    e.Position,
			StartDate:   e.StartDate,
			EndDate:     e.EndDate,
			Description: e.Description,
		})
	}

	// Convertir les formations
	for _, e := range r.Education {
		req.Education = append(req.Education, &pb.Education{
			School: e.School,
			Degree: e.Degree,
			Year:   e.Year,
		})
	}

	return req, nil
}

// GenerateCV : POST /api/generate → retourne le PDF en telechargement
func (h *Handler) GenerateCV(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "methode non autorisee", http.StatusMethodNotAllowed)
		return
	}

	// Decoder le JSON
	var body CVRequestJSON
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Convertir en Protobuf
	req, err := body.toProto()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Appeler le serveur gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := h.client.GenerateCV(ctx, req)
	if err != nil {
		http.Error(w, "Erreur gRPC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, resp.ErrorMsg, http.StatusInternalServerError)
		return
	}

	// Retourner le PDF comme fichier telechargeable
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, resp.Filename))
	w.Write(resp.PdfData)
}

// StreamProgress : GET /api/stream → SSE (Server-Sent Events)
func (h *Handler) StreamProgress(w http.ResponseWriter, r *http.Request) {
	// Headers SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming non supporte", http.StatusInternalServerError)
		return
	}

	// Decoder le JSON depuis le query param ou body
	var body CVRequestJSON
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Fprintf(w, "data: {\"error\": \"JSON invalide\"}\n\n")
		flusher.Flush()
		return
	}

	req, _ := body.toProto()

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	// Ouvrir le stream gRPC
	stream, err := h.client.StreamProgress(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\": \"%v\"}\n\n", err)
		flusher.Flush()
		return
	}

	// Lire et transmettre chaque evenement au browser
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "data: {\"error\": \"%v\"}\n\n", err)
			flusher.Flush()
			return
		}

		// Encoder l'evenement en JSON et l'envoyer via SSE
		payload := map[string]interface{}{
			"step":    event.Step,
			"percent": event.Percent,
			"done":    event.Done,
		}
		if event.Done && len(event.PdfData) > 0 {
			payload["pdf"] = base64.StdEncoding.EncodeToString(event.PdfData)
		}

		data, _ := json.Marshal(payload)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}
}

// withCORS : middleware qui ajoute les headers CORS
func (h *Handler) withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}
