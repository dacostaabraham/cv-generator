package templates

import (
	"fmt"
	"os"
	"strings"

	pb "github.com/dacostaabraham/cv-generator/proto"
	"github.com/signintech/gopdf"
)

// Palette Classic
const (
	navyR, navyG, navyB       = 28, 54, 120   // #1c3678 — bleu marine
	lightR, lightG, lightB    = 232, 240, 252 // #e8f0fc — fond header
	accentR, accentG, accentB = 59, 130, 246  // #3b82f6 — bleu accent
	darkR, darkG, darkB       = 15, 23, 42    // #0f172a — texte principal
	grayR, grayG, grayB       = 100, 116, 139 // #64748b — texte secondaire
	lineR, lineG, lineB       = 203, 213, 225 // #cbd5e1 — ligne grise legere
)

type ClassicTemplate struct{}

func (t *ClassicTemplate) Name() string { return "classic" }

func (t *ClassicTemplate) Render(req *pb.CVRequest) ([]byte, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	if err := pdf.AddTTFFont("regular", "./fonts/DejaVuSans.ttf"); err != nil {
		return nil, fmt.Errorf("police regular: %w", err)
	}
	if err := pdf.AddTTFFont("bold", "./fonts/DejaVuSans-Bold.ttf"); err != nil {
		return nil, fmt.Errorf("police bold: %w", err)
	}

	pdf.AddPage()

	pageW := 595.0
	margin := 40.0

	// ══════════════════════════════════════════════════════
	// HEADER — fond bleu clair avec bande bleue marine en bas
	// ══════════════════════════════════════════════════════
	headerH := 128.0

	// Fond principal header
	pdf.SetFillColor(lightR, lightG, lightB)
	pdf.RectFromUpperLeftWithStyle(0, 0, pageW, headerH, "F")

	// Bande bleu marine en bas du header
	pdf.SetFillColor(navyR, navyG, navyB)
	pdf.RectFromUpperLeftWithStyle(0, headerH-6, pageW, 6, "F")

	// ── Photo ─────────────────────────────────────────────
	photoX, photoY, photoW, photoH := 28.0, 14.0, 96.0, 100.0
	t.renderPhoto(pdf, req.Photo, photoX, photoY, photoW, photoH)

	// ── Nom ───────────────────────────────────────────────
	pdf.SetFont("bold", "", 26)
	pdf.SetTextColor(navyR, navyG, navyB)
	pdf.SetX(140)
	pdf.SetY(20)
	pdf.Cell(nil, req.FullName)

	// ── Ligne fine sous le nom ────────────────────────────
	pdf.SetLineWidth(0.8)
	pdf.SetStrokeColor(accentR, accentG, accentB)
	pdf.Line(140, 50, pageW-margin, 50)

	// ── Contacts ──────────────────────────────────────────
	contacts := []struct{ label, value string }{
		{"Tel", req.Phone},
		{"Email", req.Email},
		{"Lieu", req.Location},
	}
	cy := 58.0
	for _, c := range contacts {
		if c.value == "" {
			continue
		}
		// Label bleu marine bold
		pdf.SetFont("bold", "", 8)
		pdf.SetTextColor(navyR, navyG, navyB)
		pdf.SetX(140)
		pdf.SetY(cy)
		pdf.Cell(nil, c.label+":")

		// Valeur
		pdf.SetFont("regular", "", 8)
		pdf.SetTextColor(darkR, darkG, darkB)
		pdf.SetX(168)
		pdf.SetY(cy)
		pdf.Cell(nil, c.value)
		cy += 14
	}

	y := headerH + 14

	// ══════════════════════════════════════════════════════
	// ABOUT ME
	// ══════════════════════════════════════════════════════
	if req.Summary != "" {
		y = t.sectionTitle(pdf, "A PROPOS", margin, y)
		pdf.SetFont("regular", "", 9)
		pdf.SetTextColor(grayR, grayG, grayB)
		lines := wrapText(req.Summary, 95)
		for _, line := range lines {
			pdf.SetX(margin)
			pdf.SetY(y)
			pdf.Cell(nil, line)
			y += 13
		}
		y += 8
	}

	// ══════════════════════════════════════════════════════
	// EXPERIENCE
	// ══════════════════════════════════════════════════════
	if len(req.Experiences) > 0 {
		y = t.sectionTitle(pdf, "EXPERIENCE PROFESSIONNELLE", margin, y)

		for _, exp := range req.Experiences {
			// Bande coloree a gauche (timeline)
			pdf.SetFillColor(accentR, accentG, accentB)
			pdf.RectFromUpperLeftWithStyle(margin, y-2, 3, 10, "F")

			// Entreprise + poste
			pdf.SetFont("bold", "", 10)
			pdf.SetTextColor(darkR, darkG, darkB)
			pdf.SetX(margin + 12)
			pdf.SetY(y)
			pdf.Cell(nil, strings.ToUpper(exp.Company))

			// Poste en accent
			pdf.SetFont("regular", "", 9)
			pdf.SetTextColor(accentR, accentG, accentB)
			pdf.SetX(margin + 12)
			pdf.SetY(y + 13)
			pdf.Cell(nil, exp.Position)

			// Dates — aligne a droite
			dates := exp.StartDate
			if exp.EndDate != "" {
				dates += " — " + exp.EndDate
			} else {
				dates += " — present"
			}
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(grayR, grayG, grayB)
			dateW := float64(len(dates)) * 4.5
			pdf.SetX(pageW - margin - dateW)
			pdf.SetY(y)
			pdf.Cell(nil, dates)

			// Description
			if exp.Description != "" {
				descLines := wrapText(exp.Description, 85)
				descY := y + 26
				for _, line := range descLines {
					pdf.SetFont("regular", "", 8.5)
					pdf.SetTextColor(grayR, grayG, grayB)
					pdf.SetX(margin + 12)
					pdf.SetY(descY)
					pdf.Cell(nil, "· "+line)
					descY += 12
				}
				y = descY + 10
			} else {
				y += 32
			}

			// Ligne separatrice legere
			pdf.SetLineWidth(0.4)
			pdf.SetStrokeColor(lineR, lineG, lineB)
			pdf.Line(margin+12, y-4, pageW-margin, y-4)
		}
		y += 4
	}

	// ══════════════════════════════════════════════════════
	// EDUCATION
	// ══════════════════════════════════════════════════════
	if len(req.Education) > 0 {
		y = t.sectionTitle(pdf, "FORMATION", margin, y)

		for _, edu := range req.Education {
			pdf.SetFillColor(navyR, navyG, navyB)
			pdf.RectFromUpperLeftWithStyle(margin, y-2, 3, 10, "F")

			pdf.SetFont("bold", "", 10)
			pdf.SetTextColor(darkR, darkG, darkB)
			pdf.SetX(margin + 12)
			pdf.SetY(y)
			pdf.Cell(nil, edu.Degree)

			pdf.SetFont("regular", "", 9)
			pdf.SetTextColor(grayR, grayG, grayB)
			pdf.SetX(margin + 12)
			pdf.SetY(y + 13)
			pdf.Cell(nil, edu.School)

			if edu.Year > 0 {
				yr := fmt.Sprintf("%d", edu.Year)
				yrW := float64(len(yr)) * 4.5
				pdf.SetFont("bold", "", 8)
				pdf.SetTextColor(accentR, accentG, accentB)
				pdf.SetX(pageW - margin - yrW)
				pdf.SetY(y)
				pdf.Cell(nil, yr)
			}

			y += 30
			pdf.SetLineWidth(0.4)
			pdf.SetStrokeColor(lineR, lineG, lineB)
			pdf.Line(margin+12, y-4, pageW-margin, y-4)
		}
		y += 4
	}

	// ══════════════════════════════════════════════════════
	// SKILLS : grille de badges
	// ══════════════════════════════════════════════════════
	if len(req.Skills) > 0 {
		y = t.sectionTitle(pdf, "COMPETENCES", margin, y)

		// Deux colonnes : liste a gauche, barres a droite
		colMid := pageW/2 + 10
		leftY := y
		rightY := y

		half := (len(req.Skills) + 1) / 2
		percents := []int{90, 85, 80, 75, 88, 72, 94, 68, 82, 78}

		for i, skill := range req.Skills {
			if i < half {
				// Colonne gauche — badge
				pdf.SetFillColor(lightR, lightG, lightB)
				pdf.RectFromUpperLeftWithStyle(margin, leftY-1, 110, 14, "F")
				pdf.SetFillColor(accentR, accentG, accentB)
				pdf.RectFromUpperLeftWithStyle(margin, leftY-1, 3, 14, "F")
				pdf.SetFont("regular", "", 8.5)
				pdf.SetTextColor(darkR, darkG, darkB)
				pdf.SetX(margin + 8)
				pdf.SetY(leftY)
				pdf.Cell(nil, skill)
				leftY += 18
			} else {
				// Colonne droite — barre de progression
				j := i - half
				pct := 75
				if j < len(percents) {
					pct = percents[j]
				}
				pdf.SetFont("regular", "", 8.5)
				pdf.SetTextColor(darkR, darkG, darkB)
				pdf.SetX(colMid)
				pdf.SetY(rightY)
				pdf.Cell(nil, skill)

				barX := colMid + 90.0
				barW := 110.0
				barH := 6.0
				pdf.SetFillColor(lineR, lineG, lineB)
				pdf.RectFromUpperLeftWithStyle(barX, rightY+3, barW, barH, "F")
				pdf.SetFillColor(accentR, accentG, accentB)
				pdf.RectFromUpperLeftWithStyle(barX, rightY+3, barW*float64(pct)/100, barH, "F")
				pdf.SetFont("regular", "", 7)
				pdf.SetTextColor(grayR, grayG, grayB)
				pdf.SetX(barX + barW + 4)
				pdf.SetY(rightY + 1)
				pdf.Cell(nil, fmt.Sprintf("%d%%", pct))
				rightY += 18
			}
		}
	}

	return pdf.GetBytesPdf(), nil
}

// sectionTitle dessine le titre de section avec barre bleu marine
func (t *ClassicTemplate) sectionTitle(pdf *gopdf.GoPdf, title string, x, y float64) float64 {
	// Rectangle bleu marine a gauche du titre
	pdf.SetFillColor(navyR, navyG, navyB)
	pdf.RectFromUpperLeftWithStyle(x, y, 4, 14, "F")

	pdf.SetFont("bold", "", 10)
	pdf.SetTextColor(navyR, navyG, navyB)
	pdf.SetX(x + 10)
	pdf.SetY(y)
	pdf.Cell(nil, title)

	// Ligne grise legere
	pdf.SetLineWidth(0.5)
	pdf.SetStrokeColor(lineR, lineG, lineB)
	pdf.Line(x, y+16, 555, y+16)
	return y + 24
}

// renderPhoto affiche la photo ou un placeholder
func (t *ClassicTemplate) renderPhoto(pdf *gopdf.GoPdf, photoBytes []byte, x, y, w, h float64) {
	fmt.Printf("[DEBUG] renderPhoto: %d bytes\n", len(photoBytes)) // AJOUTER
	if len(photoBytes) > 0 {
		ext := ".jpg"
		if len(photoBytes) >= 4 &&
			photoBytes[0] == 0x89 && photoBytes[1] == 0x50 &&
			photoBytes[2] == 0x4E && photoBytes[3] == 0x47 {
			ext = ".png"
		}
		tmpFile, err := os.CreateTemp("", "cv-photo-*"+ext)
		if err == nil {
			_, writeErr := tmpFile.Write(photoBytes)
			tmpFile.Close()
			if writeErr == nil {
				imgErr := pdf.Image(tmpFile.Name(), x, y, &gopdf.Rect{W: w, H: h})
				os.Remove(tmpFile.Name())
				if imgErr == nil {
					return // Photo affichee avec succes
				}
				fmt.Printf("[warn] image: %v\n", imgErr)
			} else {
				os.Remove(tmpFile.Name())
			}
		}
	}
	// Placeholder
	pdf.SetFillColor(lightR, lightG, lightB)
	pdf.RectFromUpperLeftWithStyle(x, y, w, h, "F")
	pdf.SetFillColor(navyR, navyG, navyB)
	pdf.RectFromUpperLeftWithStyle(x, y+h-4, w, 4, "F")
	pdf.SetFont("regular", "", 8)
	pdf.SetTextColor(grayR, grayG, grayB)
	pdf.SetX(x + w/2 - 14)
	pdf.SetY(y + h/2 - 5)
	pdf.Cell(nil, "PHOTO")
}

// drawPhotoPlaceholder — compatibilite
func (t *ClassicTemplate) drawPhotoPlaceholder(pdf *gopdf.GoPdf, x, y, w, h float64) {
	t.renderPhoto(pdf, nil, x, y, w, h)
}

// wrapText coupe le texte en lignes de maxLen caracteres
func wrapText(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}
	words := strings.Fields(text)
	var lines []string
	current := ""
	for _, word := range words {
		if len(current)+len(word)+1 > maxLen {
			if current != "" {
				lines = append(lines, current)
			}
			current = word
		} else {
			if current == "" {
				current = word
			} else {
				current += " " + word
			}
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
