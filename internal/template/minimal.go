package templates

import (
	"fmt"
	"strings"

	"github.com/dacostaabraham/cv-generator/internal/models"
	"github.com/signintech/gopdf"
)

// MinimalTemplate — ultra-épuré, beaucoup d'espace blanc, accent teal
// Palette: #0f766e (teal foncé), #14b8a6 (teal clair), #1e293b (texte), #94a3b8 (gris)
type MinimalTemplate struct{}

func (t *MinimalTemplate) TemplateName() string { return "minimal" }

func (t *MinimalTemplate) RenderExtended(req *models.ExtendedCVRequest) ([]byte, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	if err := pdf.AddTTFFont("regular", "./fonts/DejaVuSans.ttf"); err != nil {
		return nil, fmt.Errorf("police regular: %w", err)
	}
	if err := pdf.AddTTFFont("bold", "./fonts/DejaVuSans-Bold.ttf"); err != nil {
		return nil, fmt.Errorf("police bold: %w", err)
	}

	pdf.AddPage()

	const (
		mTealDR, mTealDG, mTealDB = 15, 118, 110  // #0f766e
		mTealR, mTealG, mTealB    = 20, 184, 166  // #14b8a6
		mDarkR, mDarkG, mDarkB    = 30, 41, 59    // #1e293b
		mGrayR, mGrayG, mGrayB    = 148, 163, 184 // #94a3b8
		mLightR, mLightG, mLightB = 241, 245, 249 // fond très léger
	)

	pageW := 595.0
	marginL := 55.0
	marginR := 45.0

	// ════════════════════════════════════════════════════
	// HEADER — nom à gauche, contacts à droite
	// ════════════════════════════════════════════════════
	// Barre teal fine à gauche (toute la hauteur de la page)
	pdf.SetFillColor(mTealR, mTealG, mTealB)
	pdf.RectFromUpperLeftWithStyle(0, 0, 5, 842, "F")

	// ── Photo ───────────────────────────────────────────
	photoSize := 72.0
	hasPhoto := len(req.Photo) > 0
	nameX := marginL
	if hasPhoto {
		renderPhotoExt(pdf, req.Photo, marginL, 22, photoSize, photoSize,
			mLightR, mLightG, mLightB, mTealR, mTealG, mTealB)
		nameX = marginL + photoSize + 16
	}

	// ── Nom ─────────────────────────────────────────────
	pdf.SetFont("bold", "", 28)
	pdf.SetTextColor(mDarkR, mDarkG, mDarkB)
	pdf.SetX(nameX)
	pdf.SetY(22)
	pdf.Cell(nil, req.FullName)

	// ── Job Title ───────────────────────────────────────
	titleY := 54.0
	if req.JobTitle != "" {
		pdf.SetFont("regular", "", 11)
		pdf.SetTextColor(mTealDR, mTealDG, mTealDB)
		pdf.SetX(nameX)
		pdf.SetY(titleY)
		pdf.Cell(nil, req.JobTitle)
		titleY = 70.0
	}

	// ── Ligne teal ───────────────────────────────────────
	pdf.SetLineWidth(2.0)
	pdf.SetStrokeColor(mTealR, mTealG, mTealB)
	pdf.Line(nameX, titleY, nameX+140, titleY)

	// ── Contacts (droite) ────────────────────────────────
	contacts := []struct{ label, val string }{
		{"", req.Email},
		{"", req.Phone},
		{"", req.Location},
	}
	if req.LinkedIn != "" {
		contacts = append(contacts, struct{ label, val string }{"", truncate(req.LinkedIn, 28)})
	}
	if req.Github != "" {
		contacts = append(contacts, struct{ label, val string }{"", req.Github})
	}
	if req.Website != "" {
		contacts = append(contacts, struct{ label, val string }{"", truncate(req.Website, 28)})
	}

	cy := 20.0
	for _, c := range contacts {
		if c.val == "" {
			continue
		}
		pdf.SetFont("regular", "", 8.5)
		pdf.SetTextColor(mGrayR, mGrayG, mGrayB)
		cw := float64(len(c.val)) * 4.2
		pdf.SetX(pageW - marginR - cw)
		pdf.SetY(cy)
		pdf.Cell(nil, c.val)
		cy += 13
	}

	y := titleY + 18

	// Ligne séparateur header/body
	pdf.SetLineWidth(0.5)
	pdf.SetStrokeColor(mLightR, mLightG, mLightB)
	pdf.Line(marginL, y, pageW-marginR, y)
	y += 14

	// ════════════════════════════════════════════════════
	// HELPERS
	// ════════════════════════════════════════════════════
	sectionTitle := func(title string, yPos float64) float64 {
		// Trait teal court
		pdf.SetFillColor(mTealR, mTealG, mTealB)
		pdf.RectFromUpperLeftWithStyle(marginL, yPos+5, 18, 2, "F")
		pdf.SetFont("bold", "", 8.5)
		pdf.SetTextColor(mTealDR, mTealDG, mTealDB)
		pdf.SetX(marginL + 24)
		pdf.SetY(yPos + 1)
		pdf.Cell(nil, strings.ToUpper(title))
		return yPos + 18
	}

	// ════════════════════════════════════════════════════
	// LAYOUT 2 COLONNES
	// ════════════════════════════════════════════════════
	leftW := (pageW - marginL - marginR) * 0.62
	rightX := marginL + leftW + 16
	rightW := pageW - rightX - marginR

	leftY := y
	rightY := y

	// ── COLONNE GAUCHE — Summary + Expériences + Formation ──

	// About
	if req.Summary != "" {
		leftY = sectionTitle("Profil", leftY)
		for _, line := range wrapStr(req.Summary, 72) {
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(mDarkR, mDarkG, mDarkB)
			pdf.SetX(marginL)
			pdf.SetY(leftY)
			pdf.Cell(nil, line)
			leftY += 12
		}
		leftY += 10
	}

	// Expériences
	if len(req.Experiences) > 0 {
		leftY = sectionTitle("Experience", leftY)
		for _, exp := range req.Experiences {
			// Date — très léger, à droite
			dates := exp.StartDate
			if exp.EndDate != "" {
				dates += " – " + exp.EndDate
			} else {
				dates += " – Présent"
			}
			dw := float64(len(dates)) * 4.0
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(mGrayR, mGrayG, mGrayB)
			pdf.SetX(marginL + leftW - dw)
			pdf.SetY(leftY)
			pdf.Cell(nil, dates)

			// Entreprise
			pdf.SetFont("bold", "", 9.5)
			pdf.SetTextColor(mDarkR, mDarkG, mDarkB)
			pdf.SetX(marginL)
			pdf.SetY(leftY)
			pdf.Cell(nil, exp.Company)

			// Poste
			leftY += 13
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(mTealDR, mTealDG, mTealDB)
			pdf.SetX(marginL)
			pdf.SetY(leftY)
			pdf.Cell(nil, exp.Position)

			// Description
			if exp.Description != "" {
				leftY += 12
				for _, line := range wrapStr(exp.Description, 70) {
					pdf.SetFont("regular", "", 8)
					pdf.SetTextColor(mGrayR, mGrayG, mGrayB)
					pdf.SetX(marginL + 4)
					pdf.SetY(leftY)
					pdf.Cell(nil, "– "+line)
					leftY += 11
				}
			} else {
				leftY += 8
			}
			// Ligne fine
			pdf.SetLineWidth(0.3)
			pdf.SetStrokeColor(mLightR, mLightG, mLightB)
			pdf.Line(marginL, leftY+4, marginL+leftW, leftY+4)
			leftY += 12
		}
	}

	// Formation
	if len(req.Education) > 0 {
		leftY = sectionTitle("Formation", leftY)
		for _, edu := range req.Education {
			pdf.SetFont("bold", "", 9.5)
			pdf.SetTextColor(mDarkR, mDarkG, mDarkB)
			pdf.SetX(marginL)
			pdf.SetY(leftY)
			pdf.Cell(nil, edu.Degree)

			if edu.Year > 0 {
				yr := fmt.Sprintf("%d", edu.Year)
				yw := float64(len(yr)) * 4.5
				pdf.SetFont("bold", "", 8)
				pdf.SetTextColor(mTealR, mTealG, mTealB)
				pdf.SetX(marginL + leftW - yw)
				pdf.SetY(leftY)
				pdf.Cell(nil, yr)
			}

			leftY += 13
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(mGrayR, mGrayG, mGrayB)
			pdf.SetX(marginL)
			pdf.SetY(leftY)
			pdf.Cell(nil, edu.School)
			leftY += 18
		}
	}

	// ── COLONNE DROITE — Skills + Langues + Certifs ──────

	// Skills
	if len(req.Skills) > 0 {
		rightY = sectionTitle("Competences", rightY)
		_ = rightW
		for _, skill := range req.Skills {
			// Tag minimal avec bordure teal
			sw := float64(len(skill))*4.5 + 14
			if sw > rightW {
				sw = rightW
			}
			pdf.SetFillColor(mLightR, mLightG, mLightB)
			pdf.RectFromUpperLeftWithStyle(rightX, rightY-1, sw, 13, "F")
			// Bordure gauche teal
			pdf.SetFillColor(mTealR, mTealG, mTealB)
			pdf.RectFromUpperLeftWithStyle(rightX, rightY-1, 2, 13, "F")
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(mDarkR, mDarkG, mDarkB)
			pdf.SetX(rightX + 7)
			pdf.SetY(rightY)
			pdf.Cell(nil, truncate(skill, 18))
			rightY += 16
		}
		rightY += 6
	}

	// Langues
	if len(req.Languages) > 0 {
		rightY = sectionTitle("Langues", rightY)
		for _, lang := range req.Languages {
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(mDarkR, mDarkG, mDarkB)
			pdf.SetX(rightX)
			pdf.SetY(rightY)
			pdf.Cell(nil, "· "+lang)
			rightY += 14
		}
		rightY += 6
	}

	// Certifications
	if len(req.Certifications) > 0 {
		rightY = sectionTitle("Certifications", rightY)
		for _, cert := range req.Certifications {
			for _, line := range wrapStr(cert, 22) {
				pdf.SetFont("regular", "", 7.5)
				pdf.SetTextColor(mGrayR, mGrayG, mGrayB)
				pdf.SetX(rightX)
				pdf.SetY(rightY)
				pdf.Cell(nil, "✓ "+line)
				rightY += 13
			}
		}
	}

	// Séparateur vertical
	maxH := leftY
	if rightY > maxH {
		maxH = rightY
	}
	pdf.SetLineWidth(0.4)
	pdf.SetStrokeColor(mLightR, mLightG, mLightB)
	pdf.Line(rightX-8, y-14, rightX-8, maxH+10)

	return pdf.GetBytesPdf(), nil
}
