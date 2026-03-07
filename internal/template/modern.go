package templates

import (
	"fmt"
	"os"
	"strings"

	pb "github.com/dacostaabraham/cv-generator/proto"
	"github.com/signintech/gopdf"
)

// Palette Modern
const (
	mSideR, mSideG, mSideB = 15, 23, 42    // #0f172a — sidebar tres sombre
	mAccR, mAccG, mAccB    = 99, 102, 241  // #6366f1 — violet/indigo accent
	mAcc2R, mAcc2G, mAcc2B = 34, 211, 238  // #22d3ee — teal accent secondaire
	mTextR, mTextG, mTextB = 248, 250, 252 // #f8fafc — texte clair sidebar
	mMutR, mMutG, mMutB    = 148, 163, 184 // #94a3b8 — texte mute sidebar
	mMainR, mMainG, mMainB = 248, 250, 252 // fond page principal
	mDarkR, mDarkG, mDarkB = 15, 23, 42    // texte sombre colonne principale
	mGrayR, mGrayG, mGrayB = 100, 116, 139 // gris neutre
	mLineR, mLineG, mLineB = 226, 232, 240 // ligne separatrice
)

type ModernTemplate struct{}

func (t *ModernTemplate) Name() string { return "modern" }

func (t *ModernTemplate) Render(req *pb.CVRequest) ([]byte, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	if err := pdf.AddTTFFont("regular", "./fonts/DejaVuSans.ttf"); err != nil {
		return nil, fmt.Errorf("police regular: %w", err)
	}
	if err := pdf.AddTTFFont("bold", "./fonts/DejaVuSans-Bold.ttf"); err != nil {
		return nil, fmt.Errorf("police bold: %w", err)
	}

	pdf.AddPage()

	sideW := 190.0
	mainX := sideW + 24.0
	mainW := 595.0 - mainX - 20.0
	_ = mainW
	sideM := 18.0

	// ══════════════════════════════════════════════════════
	// SIDEBAR — fond tres sombre
	// ══════════════════════════════════════════════════════
	pdf.SetFillColor(mSideR, mSideG, mSideB)
	pdf.RectFromUpperLeftWithStyle(0, 0, sideW, 842, "F")

	// Bande accent en haut de la sidebar
	pdf.SetFillColor(mAccR, mAccG, mAccB)
	pdf.RectFromUpperLeftWithStyle(0, 0, sideW, 5, "F")

	// ── Photo ─────────────────────────────────────────────
	photoSize := 80.0
	photoX := (sideW - photoSize) / 2
	t.renderPhoto(pdf, req.Photo, photoX, 18, photoSize, photoSize)

	// Ligne decorative sous la photo
	pdf.SetFillColor(mAccR, mAccG, mAccB)
	pdf.RectFromUpperLeftWithStyle(sideM, 104, sideW-sideM*2, 2, "F")

	// ── Nom ───────────────────────────────────────────────
	pdf.SetFont("bold", "", 14)
	pdf.SetTextColor(mTextR, mTextG, mTextB)
	// Centrer le nom
	nameLines := wrapTextM(req.FullName, 18)
	nameY := 112.0
	for _, line := range nameLines {
		textW := float64(len(line)) * 7.2
		pdf.SetX((sideW - textW) / 2)
		pdf.SetY(nameY)
		pdf.Cell(nil, line)
		nameY += 18
	}
	nameY += 4

	// ── Contact ───────────────────────────────────────────
	nameY = t.sideSection(pdf, "CONTACT", sideM, nameY)
	contacts := []struct{ icon, val string }{
		{"@", req.Email},
		{"T", req.Phone},
		{"L", req.Location},
	}
	for _, c := range contacts {
		if c.val == "" {
			continue
		}
		// Icone badge
		pdf.SetFillColor(mAccR, mAccG, mAccB)
		pdf.RectFromUpperLeftWithStyle(sideM, nameY-1, 12, 12, "F")
		pdf.SetFont("bold", "", 7)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetX(sideM + 2)
		pdf.SetY(nameY)
		pdf.Cell(nil, c.icon)

		pdf.SetFont("regular", "", 7.5)
		pdf.SetTextColor(mMutR, mMutG, mMutB)
		pdf.SetX(sideM + 16)
		pdf.SetY(nameY)
		// Couper si trop long pour la sidebar
		val := c.val
		if len(val) > 24 {
			val = val[:22] + ".."
		}
		pdf.Cell(nil, val)
		nameY += 15
	}
	nameY += 6

	// ── Competences ───────────────────────────────────────
	if len(req.Skills) > 0 {
		nameY = t.sideSection(pdf, "COMPETENCES", sideM, nameY)
		percents := []int{95, 88, 82, 78, 90, 72, 85, 68, 80, 75}
		for i, skill := range req.Skills {
			pct := 78
			if i < len(percents) {
				pct = percents[i]
			}
			// Nom du skill
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(mTextR, mTextG, mTextB)
			pdf.SetX(sideM)
			pdf.SetY(nameY)
			pdf.Cell(nil, skill)

			// Barre de fond
			barW := sideW - sideM*2
			pdf.SetFillColor(30, 41, 59)
			pdf.RectFromUpperLeftWithStyle(sideM, nameY+10, barW, 4, "F")
			// Barre remplie en gradient accent
			filled := barW * float64(pct) / 100.0
			pdf.SetFillColor(mAccR, mAccG, mAccB)
			pdf.RectFromUpperLeftWithStyle(sideM, nameY+10, filled, 4, "F")
			// Petit point a la fin
			pdf.SetFillColor(mAcc2R, mAcc2G, mAcc2B)
			pdf.RectFromUpperLeftWithStyle(sideM+filled-2, nameY+9, 6, 6, "F")

			nameY += 22
		}
		nameY += 4
	}

	// ── Profil dans la sidebar ─────────────────────────────
	if req.Summary != "" {
		nameY = t.sideSection(pdf, "PROFIL", sideM, nameY)
		summary := req.Summary
		if len(summary) > 160 {
			summary = summary[:160] + "..."
		}
		lines := wrapTextM(summary, 22)
		for _, line := range lines {
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(mMutR, mMutG, mMutB)
			pdf.SetX(sideM)
			pdf.SetY(nameY)
			pdf.Cell(nil, line)
			nameY += 12
		}
	}

	// ══════════════════════════════════════════════════════
	// COLONNE PRINCIPALE — fond blanc casse
	// ══════════════════════════════════════════════════════
	pdf.SetFillColor(255, 255, 255)
	pdf.RectFromUpperLeftWithStyle(sideW, 0, 595-sideW, 842, "F")

	// Bande accent en haut de la colonne principale
	pdf.SetFillColor(mAccR, mAccG, mAccB)
	pdf.RectFromUpperLeftWithStyle(sideW, 0, 595-sideW, 5, "F")

	mainY := 22.0

	// Titre du poste / nom en haut de la colonne principale
	pdf.SetFont("bold", "", 18)
	pdf.SetTextColor(mDarkR, mDarkG, mDarkB)
	pdf.SetX(mainX)
	pdf.SetY(mainY)
	// Premier prenom seulement pour titre
	firstName := strings.SplitN(req.FullName, " ", 2)[0]
	pdf.Cell(nil, firstName)

	mainY += 22

	// Tag "poste" si la premiere experience existe
	if len(req.Experiences) > 0 {
		role := req.Experiences[0].Position
		if len(role) > 0 {
			tagW := float64(len(role))*5.5 + 14
			pdf.SetFillColor(mAccR, mAccG, mAccB)
			pdf.RectFromUpperLeftWithStyle(mainX, mainY, tagW, 14, "F")
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetX(mainX + 7)
			pdf.SetY(mainY + 2)
			pdf.Cell(nil, role)
			mainY += 22
		}
	}

	// ── EXPERIENCE ────────────────────────────────────────
	if len(req.Experiences) > 0 {
		mainY = t.mainSection(pdf, "EXPERIENCE", mainX, mainY)

		for _, exp := range req.Experiences {
			// Dot timeline
			pdf.SetFillColor(mAccR, mAccG, mAccB)
			pdf.RectFromUpperLeftWithStyle(mainX, mainY+3, 8, 8, "F")

			// Entreprise
			pdf.SetFont("bold", "", 10)
			pdf.SetTextColor(mDarkR, mDarkG, mDarkB)
			pdf.SetX(mainX + 14)
			pdf.SetY(mainY)
			pdf.Cell(nil, exp.Company)

			// Dates aligne droite
			dates := exp.StartDate
			if exp.EndDate != "" {
				dates += " — " + exp.EndDate
			} else {
				dates += " — present"
			}
			dateW := float64(len(dates)) * 4.2
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(mAccR, mAccG, mAccB)
			pdf.SetX(595 - 20 - dateW)
			pdf.SetY(mainY + 1)
			pdf.Cell(nil, dates)

			// Poste
			mainY += 13
			pdf.SetFont("regular", "", 9)
			pdf.SetTextColor(mAccR, mAccG, mAccB)
			pdf.SetX(mainX + 14)
			pdf.SetY(mainY)
			pdf.Cell(nil, exp.Position)

			// Description
			if exp.Description != "" {
				mainY += 13
				descLines := wrapTextM(exp.Description, 52)
				for _, line := range descLines {
					pdf.SetFont("regular", "", 8.5)
					pdf.SetTextColor(mGrayR, mGrayG, mGrayB)
					pdf.SetX(mainX + 14)
					pdf.SetY(mainY)
					pdf.Cell(nil, "· "+line)
					mainY += 12
				}
			} else {
				mainY += 10
			}

			// Ligne timeline verticale
			pdf.SetLineWidth(0.5)
			pdf.SetStrokeColor(mLineR, mLineG, mLineB)
			pdf.Line(mainX+3, mainY+2, mainX+3, mainY+12)

			mainY += 14
		}
		mainY += 4
	}

	// ── FORMATION ─────────────────────────────────────────
	if len(req.Education) > 0 {
		mainY = t.mainSection(pdf, "FORMATION", mainX, mainY)

		for _, edu := range req.Education {
			pdf.SetFillColor(mAccR, mAccG, mAccB)
			pdf.RectFromUpperLeftWithStyle(mainX, mainY+3, 8, 8, "F")

			pdf.SetFont("bold", "", 10)
			pdf.SetTextColor(mDarkR, mDarkG, mDarkB)
			pdf.SetX(mainX + 14)
			pdf.SetY(mainY)
			pdf.Cell(nil, edu.Degree)

			if edu.Year > 0 {
				yr := fmt.Sprintf("%d", edu.Year)
				yrW := float64(len(yr)) * 5.0
				pdf.SetFont("bold", "", 8)
				pdf.SetTextColor(mAccR, mAccG, mAccB)
				pdf.SetX(595 - 20 - yrW)
				pdf.SetY(mainY + 1)
				pdf.Cell(nil, yr)
			}

			mainY += 13
			pdf.SetFont("regular", "", 9)
			pdf.SetTextColor(mGrayR, mGrayG, mGrayB)
			pdf.SetX(mainX + 14)
			pdf.SetY(mainY)
			pdf.Cell(nil, edu.School)

			mainY += 18
		}
	}

	return pdf.GetBytesPdf(), nil
}

// sideSection : titre de section dans la sidebar
func (t *ModernTemplate) sideSection(pdf *gopdf.GoPdf, title string, x, y float64) float64 {
	// Ligne teal courte avant le titre
	pdf.SetFillColor(mAcc2R, mAcc2G, mAcc2B)
	pdf.RectFromUpperLeftWithStyle(x, y, 20, 2, "F")
	y += 6

	pdf.SetFont("bold", "", 8)
	pdf.SetTextColor(mAcc2R, mAcc2G, mAcc2B)
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, title)
	return y + 14
}

// mainSection : titre de section dans la colonne principale
func (t *ModernTemplate) mainSection(pdf *gopdf.GoPdf, title string, x, y float64) float64 {
	pdf.SetFillColor(mAccR, mAccG, mAccB)
	pdf.RectFromUpperLeftWithStyle(x, y, 3, 14, "F")

	pdf.SetFont("bold", "", 10)
	pdf.SetTextColor(mAccR, mAccG, mAccB)
	pdf.SetX(x + 9)
	pdf.SetY(y)
	pdf.Cell(nil, title)

	pdf.SetLineWidth(0.4)
	pdf.SetStrokeColor(mLineR, mLineG, mLineB)
	pdf.Line(x, y+16, 575, y+16)
	return y + 24
}

// renderPhoto affiche la photo ou un placeholder
func (t *ModernTemplate) renderPhoto(pdf *gopdf.GoPdf, photoBytes []byte, x, y, w, h float64) {
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
					return
				}
				fmt.Printf("[warn] photo modern: %v\n", imgErr)
			} else {
				os.Remove(tmpFile.Name())
			}
		}
	}
	// Placeholder cercle
	pdf.SetFillColor(30, 41, 59)
	pdf.RectFromUpperLeftWithStyle(x, y, w, h, "F")
	pdf.SetFillColor(mAccR, mAccG, mAccB)
	pdf.RectFromUpperLeftWithStyle(x, y+h-3, w, 3, "F")
	pdf.SetFont("regular", "", 7)
	pdf.SetTextColor(mMutR, mMutG, mMutB)
	pdf.SetX(x + w/2 - 10)
	pdf.SetY(y + h/2 - 4)
	pdf.Cell(nil, "PHOTO")
}

// wrapTextM coupe le texte pour la colonne etroite de la sidebar
func wrapTextM(text string, maxLen int) []string {
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
