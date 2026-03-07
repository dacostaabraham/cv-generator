package main

import (
	"fmt"

	"github.com/dacostaabraham/cv-generator/internal/models"
)

func main() {
	// Créer un CV de test
	cv := models.NewCVData("Abraham Kouassi", "dacostaabraham556@gmail.com")
	cv.Phone = "+225 07 00 00 00"
	cv.Location = "Abidjan, Côte d'Ivoire"
	cv.Summary = "Lead Engineer passionné par les systèmes distribués."
	cv.Skills = []string{"Go", "PHP", "Flutter", "gRPC", "MySQL"}

	cv.Experiences = []models.Experience{
		{
			Company:   "FAKODROP",
			Position:  "Lead Technical Engineer",
			StartDate: "2022",
			EndDate:   "", // présent
		},
	}

	// Afficher les infos
	fmt.Printf("✅ CV créé : %s\n", cv.FullName)
	fmt.Printf("📧 Email   : %s\n", cv.Email)
	fmt.Printf("🛠  Skills  : %s\n", cv.SkillsString())
	fmt.Printf("✔  Valide  : %v\n", cv.IsValid())
	fmt.Printf("💼 Expériences : %d\n", len(cv.Experiences))
}
