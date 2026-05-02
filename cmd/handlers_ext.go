package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	pdfext "github.com/dacostaabraham/cv-generator/internal/pdf"
	"github.com/dacostaabraham/cv-generator/internal/models"
)

// ════════════════════════════════════════════════════════
// HANDLER V2 — bypass gRPC, utilise le modèle étendu
// ════════════════════════════════════════════════════════

// CVRequestExtJSON : JSON reçu du browser (version étendue)
type CVRequestExtJSON struct {
	// Champs de base
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

	// Nouveaux champs v2
	JobTitle       string   `json:"job_title"`
	LinkedIn       string   `json:"linkedin"`
	Github         string   `json:"github"`
	Website        string   `json:"website"`
	Languages      []string `json:"languages"`
	Certifications []string `json:"certifications"`
}

// toExtended convertit le JSON en ExtendedCVRequest
func (r CVRequestExtJSON) toExtended() (*models.ExtendedCVRequest, error) {
	req := &models.ExtendedCVRequest{
		FullName:       r.FullName,
		Email:          r.Email,
		Phone:          r.Phone,
		Location:       r.Location,
		Summary:        r.Summary,
		Skills:         r.Skills,
		Template:       r.Template,
		JobTitle:       r.JobTitle,
		LinkedIn:       r.LinkedIn,
		Github:         r.Github,
		Website:        r.Website,
		Languages:      r.Languages,
		Certifications: r.Certifications,
	}

	// Décoder la photo base64
	if r.Photo != "" {
		photoBytes, err := base64.StdEncoding.DecodeString(r.Photo)
		if err != nil {
			fmt.Printf("[warn] photo base64 invalide: %v\n", err)
		} else {
			req.Photo = photoBytes
		}
	}

	// Convertir les expériences
	for _, e := range r.Experiences {
		req.Experiences = append(req.Experiences, &models.ExperienceExt{
			Company:     e.Company,
			Position:    e.Position,
			StartDate:   e.StartDate,
			EndDate:     e.EndDate,
			Description: e.Description,
		})
	}

	// Convertir les formations
	for _, e := range r.Education {
		req.Education = append(req.Education, &models.EducationExt{
			School: e.School,
			Degree: e.Degree,
			Year:   e.Year,
		})
	}

	if req.Template == "" {
		req.Template = "classic"
	}

	return req, nil
}

// StreamExtended : POST /api/v2/stream — SSE avec modèle étendu (no gRPC)
// Requiert un token de paiement valide dans le header Authorization: Bearer <token>
func (h *Handler) StreamExtended(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming non supporté", http.StatusInternalServerError)
		return
	}

	sendEvent := func(step string, percent int, done bool, pdfData []byte) {
		payload := map[string]interface{}{
			"step":    step,
			"percent": percent,
			"done":    done,
		}
		if done && len(pdfData) > 0 {
			payload["pdf"] = base64.StdEncoding.EncodeToString(pdfData)
		}
		data, _ := json.Marshal(payload)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	// ── Vérification du token de paiement ────────────────
	token := r.Header.Get("X-Payment-Token")
	if token == "" {
		// Fallback : query param (utile pour les EventSource qui ne supportent pas les headers)
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		sendEvent("Erreur: Paiement requis pour générer un CV", 0, true, nil)
		return
	}
	if _, err := VerifyPaymentToken(token); err != nil {
		sendEvent(fmt.Sprintf("Erreur: Token invalide (%s)", err.Error()), 0, true, nil)
		return
	}
	fmt.Printf("[v2] ✅ Token paiement valide\n")

	// Décoder la requête
	var body CVRequestExtJSON
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendEvent("Erreur: JSON invalide", 0, true, nil)
		return
	}

	// Valider
	if body.FullName == "" || body.Email == "" {
		sendEvent("Erreur: Nom et email obligatoires", 0, true, nil)
		return
	}

	fmt.Printf("[v2] Stream pour: %s (template: %s)\n", body.FullName, body.Template)

	// Progression simulée
	steps := []struct {
		step    string
		percent int
		delay   time.Duration
	}{
		{"Validation des données...", 12, 250 * time.Millisecond},
		{"Chargement du template...", 28, 350 * time.Millisecond},
		{"Mise en page du header...", 45, 400 * time.Millisecond},
		{"Rendu des sections...", 65, 450 * time.Millisecond},
		{"Génération du PDF...", 82, 300 * time.Millisecond},
		{"Finalisation...", 95, 200 * time.Millisecond},
	}

	for _, s := range steps {
		select {
		case <-r.Context().Done():
			return
		default:
			sendEvent(s.step, s.percent, false, nil)
			time.Sleep(s.delay)
		}
	}

	// Convertir et générer
	req, err := body.toExtended()
	if err != nil {
		sendEvent(fmt.Sprintf("Erreur: %v", err), 100, true, nil)
		return
	}

	pdfData, err := pdfext.GenerateExtended(req)
	if err != nil {
		sendEvent(fmt.Sprintf("Erreur PDF: %v", err), 100, true, nil)
		return
	}

	sendEvent("CV généré avec succès !", 100, true, pdfData)
}
