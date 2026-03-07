package main

import (
	"context"
	"fmt"
	"time"

	pdfgen "github.com/dacostaabraham/cv-generator/internal/pdf" // AJOUTER
	pb "github.com/dacostaabraham/cv-generator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CVService implementer l'interface CVGeneratorServer generee par protoc
type CVService struct {
	// Embed pour eviter l'erreur si une methode manque
	pb.UnimplementedCVGeneratorServer
}

// validate verifie que les champs obligatoires sont remplis
func validate(req *pb.CVRequest) error {
	if req.FullName == "" {
		return status.Error(codes.InvalidArgument, "full_name est obligatoire")
	}
	if req.Email == "" {
		return status.Error(codes.InvalidArgument, "email est obligatoire")
	}
	return nil
}

// buildFilename construit le nom du fichier PDF
func buildFilename(fullName string) string {
	// "Abraham Kouassi" -> "abraham-kouassi-cv.pdf"
	result := ""
	lower := false
	for _, ch := range fullName {
		if ch == ' ' {
			result += "-"
		} else if ch >= 'A' && ch <= 'Z' {
			result += string(ch + 32)
			lower = true
		} else {
			result += string(ch)
			lower = true
		}
	}
	_ = lower
	return result + "-cv.pdf"
}

// Recoit les donnees -> retourne le PDF en une seule reponse
func (s *CVService) GenerateCV(
	ctx context.Context,
	req *pb.CVRequest,
) (*pb.CVResponse, error) {

	// 1. Verifier que le contexte n'est pas annule (timeout, cancel)
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "requete annulee")
	default:
	}

	// 2. Valider les donnees
	if err := validate(req); err != nil {
		return nil, err
	}

	// 3. Logger la requete (utile pour debug)
	fmt.Printf("[GenerateCV] Requete recue pour: %s (template: %s)\n",
		req.FullName, req.Template)

	// 4. Choisir le template (defaut: classic)
	tmpl := req.Template
	if tmpl == "" {
		tmpl = "classic"
	}

	// 5. Generer le PDF
	// Pour l'instant: PDF factice avec les donnees du CV
	// A l'etape 4 on remplacera par la vraie generation gofpdf
	pdfData, err := pdfgen.Generate(req)
	if err != nil {
		return &pb.CVResponse{
			Success:  false,
			ErrorMsg: fmt.Sprintf("erreur generation PDF: %v", err),
		}, nil
	}

	// 6. Construire et retourner la reponse
	return &pb.CVResponse{
		Success:  true,
		PdfData:  pdfData,
		Filename: buildFilename(req.FullName),
	}, nil
}

// generateMockPDF : simule un PDF (remplace par gofpdf a l'etape 4)
func generateMockPDF(req *pb.CVRequest, tmpl string) []byte {
	content := fmt.Sprintf(
		"CV - %s | %s | template:%s | skills:%d | exp:%d",
		req.FullName,
		req.Email,
		tmpl,
		len(req.Skills),
		len(req.Experiences),
	)
	return []byte(content)
}

func (s *CVService) StreamProgress(
	req *pb.CVRequest,
	stream pb.CVGenerator_StreamProgressServer,
) error {

	// 1. Valider
	if err := validate(req); err != nil {
		return err
	}

	fmt.Printf("[StreamProgress] Streaming pour: %s\n", req.FullName)

	// 2. Definir les etapes de progression
	steps := []struct {
		step    string
		percent int32
		delay   time.Duration
	}{
		{"Validation des donnees...", 10, 300 * time.Millisecond},
		{"Chargement du template...", 25, 400 * time.Millisecond},
		{"Mise en page du CV...", 50, 600 * time.Millisecond},
		{"Generation du PDF...", 75, 800 * time.Millisecond},
		{"Finalisation...", 90, 300 * time.Millisecond},
	}

	// 3. Envoyer chaque etape au client
	for _, s := range steps {
		// Verifier que le client est toujours connecte
		if err := stream.Context().Err(); err != nil {
			return status.Error(codes.Canceled, "client deconnecte")
		}

		// Envoyer l'evenement de progression
		if err := stream.Send(&pb.ProgressEvent{
			Step:    s.step,
			Percent: s.percent,
			Done:    false,
		}); err != nil {
			return err
		}

		// Simuler le travail
		time.Sleep(s.delay)
	}

	// 4. Generer le PDF final
	pdfData, err := pdfgen.Generate(req)
	if err != nil {
		return status.Errorf(codes.Internal, "erreur PDF: %v", err)
	}

	// 5. Envoyer l'evenement final avec le PDF
	return stream.Send(&pb.ProgressEvent{
		Step:    "CV genere avec succes !",
		Percent: 100,
		Done:    true,
		PdfData: pdfData,
	})
}
