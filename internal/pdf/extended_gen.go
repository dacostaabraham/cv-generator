package pdf

import (
	"fmt"

	"github.com/dacostaabraham/cv-generator/internal/models"
	templates "github.com/dacostaabraham/cv-generator/internal/template"
)

// ExtendedTemplate : interface commune pour tous les templates (v2)
type ExtendedTemplate interface {
	RenderExtended(req *models.ExtendedCVRequest) ([]byte, error)
	TemplateName() string
}

// GenerateExtended choisit le template et génère le PDF
// en utilisant le modèle étendu (sans gRPC)
func GenerateExtended(req *models.ExtendedCVRequest) ([]byte, error) {
	var tmpl ExtendedTemplate

	switch req.Template {
	case "modern":
		tmpl = &templates.ModernTemplateExt{}
	case "executive":
		tmpl = &templates.ExecutiveTemplate{}
	case "minimal":
		tmpl = &templates.MinimalTemplate{}
	case "creative":
		tmpl = &templates.CreativeTemplate{}
	default: // "classic" ou valeur vide
		tmpl = &templates.ClassicTemplateExt{}
	}

	fmt.Printf("[PDF-EXT] Génération avec template '%s' pour %s\n",
		tmpl.TemplateName(), req.FullName)

	return tmpl.RenderExtended(req)
}
