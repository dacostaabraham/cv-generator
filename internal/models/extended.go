package models

// ════════════════════════════════════════════════════════
// MODÈLE ÉTENDU — utilisé par les nouveaux templates
// Contient tous les champs y compris les nouveaux
// ════════════════════════════════════════════════════════

// ExperienceExt = expérience professionnelle
type ExperienceExt struct {
	Company     string
	Position    string
	StartDate   string
	EndDate     string
	Description string
}

// EducationExt = formation / diplôme
type EducationExt struct {
	School string
	Degree string
	Year   int32
}

// ExtendedCVRequest = toutes les données du CV (version complète)
type ExtendedCVRequest struct {
	// ── Champs de base (compatibles avec proto) ────────
	FullName    string
	Email       string
	Phone       string
	Location    string
	Summary     string
	Skills      []string
	Experiences []*ExperienceExt
	Education   []*EducationExt
	Template    string // classic | modern | executive | minimal | creative
	Photo       []byte

	// ── Nouveaux champs ────────────────────────────────
	JobTitle       string   // Ex: "Senior Fullstack Engineer"
	LinkedIn       string   // URL LinkedIn
	Github         string   // URL/pseudo GitHub
	Website        string   // Portfolio ou site perso
	Languages      []string // Ex: ["Français", "Anglais", "Dioula"]
	Certifications []string // Ex: ["AWS Solutions Architect", "CKAD"]
}
