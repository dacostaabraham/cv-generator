package pdf

import (
	"fmt"

	templates "github.com/dacostaabraham/cv-generator.git/internal/template"
	pb "github.com/dacostaabraham/cv-generator.git/proto"
)

// CVTemplate : tout type qui implemente Render() et Name()
// est un template valide — Classic et Modern implementent cette interface
type CVTemplate interface {
	Render(req *pb.CVRequest) ([]byte, error)
	Name() string
}

// Generate choisit le bon template et retourne le PDF en bytes
func Generate(req *pb.CVRequest) ([]byte, error) {
	// Choisir le template selon la requete
	var tmpl CVTemplate

	switch req.Template {
	case "modern":
		tmpl = &templates.ModernTemplate{}
	default: // "classic" ou valeur vide
		tmpl = &templates.ClassicTemplate{}
	}

	fmt.Printf("[PDF] Generation avec template '%s' pour %s\n",
		tmpl.Name(), req.FullName)

	// Deleguer le rendu au template choisi
	return tmpl.Render(req)
}
