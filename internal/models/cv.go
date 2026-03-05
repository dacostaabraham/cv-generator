package models

import "strings"

// ════════════════════════════════════════════════════════
// MODÈLES DE DONNÉES
// Ces structs seront traduits en messages Protobuf à l'étape 2
// ════════════════════════════════════════════════════════

type Experience struct {
	Company     string
	Position    string
	StartDate   string
	EndDate     string
	Description string
}

type Education struct {
	School string
	Degree string
	Year   int
}

type CVData struct {
	FullName    string
	Email       string
	Phone       string
	Location    string
	Summary     string
	Skills      []string
	Experiences []Experience
	Education   []Education
	Template    string // "classic" | "modern"
}

// ════════════════════════════════════════════════════════
// MÉTHODES
// ════════════════════════════════════════════════════════

func (cv CVData) IsValid() bool {
	return cv.FullName != "" && cv.Email != ""
}

func (cv CVData) SkillsString() string {
	return strings.Join(cv.Skills, ", ")
}

// ════════════════════════════════════════════════════════
// INTERFACE
// ════════════════════════════════════════════════════════

type CVTemplate interface {
	Render(data CVData) ([]byte, error)
	Name() string
}

// ════════════════════════════════════════════════════════
// CONSTRUCTEUR (pattern Go)
// ════════════════════════════════════════════════════════

// NewCVData crée un CVData avec des valeurs par défaut
func NewCVData(name, email string) CVData {
	return CVData{
		FullName: name,
		Email:    email,
		Template: "classic",
		Skills:   []string{},
	}
}
