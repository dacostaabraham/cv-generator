package templates

import (
	"fmt"
	"strings"

	"github.com/dacostaabraham/cv-generator/internal/models"
	"github.com/signintech/gopdf"
)

// ExecutiveTemplate — Corporate, élégant, noir & or
// Palette: #1a1a2e (navy noir), #c9a84c (or), #f5f5f0 (crème), #555 (gris)
type ExecutiveTemplate struct{}

func (t *ExecutiveTemplate) TemplateName() string { return "executive" }

func (t *ExecutiveTemplate) RenderExtended(req *models.ExtendedCVRequest) ([]byte, error) {
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
		eNR, eNG, eNB   = 26, 26, 46  // #1a1a2e navy noir
		eGR, eGG, eGB   = 201, 168, 76 // #c9a84c or/gold
		eCrR, eCrG, eCrB = 245, 245, 240 // crème fond
		eTextR, eTextG, eTextB = 40, 40, 40 // texte principal
		eGrayR, eGrayG, eGrayB = 110, 110, 120 // gris
		eLR, eLG, eLB   = 220, 215, 200 // ligne séparateur
	)

	pageW := 595.0
	marginL := 50.0
	marginR := 40.0
	contentW := pageW - marginL - marginR

	// ════════════════════════════════════════════════════
	// HEADER — Bande noire pleine en haut
	// ════════════════════════════════════════════════════
	headerH := 120.0
	pdf.SetFillColor(eNR, eNG, eNB)
	pdf.RectFromUpperLeftWithStyle(0, 0, pageW, headerH, "F")

	// Bande or fine en bas du header
	pdf.SetFillColor(eGR, eGG, eGB)
	pdf.RectFromUpperLeftWithStyle(0, headerH-3, pageW, 3, "F")

	// ── Photo (cercle dans le header) ───────────────────
	photoSize := 82.0
	photoX := marginL
	renderPhotoExt(pdf, req.Photo, photoX, (headerH-photoSize)/2, photoSize, photoSize,
		50, 50, 80, eGR, eGG, eGB)

	// ── Nom ─────────────────────────────────────────────
	nameX := photoX + photoSize + 22
	pdf.SetFont("bold", "", 26)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetX(nameX)
	pdf.SetY(22)
	pdf.Cell(nil, req.FullName)

	// ── Job Title ───────────────────────────────────────
	if req.JobTitle != "" {
		pdf.SetFont("regular", "", 11)
		pdf.SetTextColor(eGR, eGG, eGB)
		pdf.SetX(nameX)
		pdf.SetY(52)
		pdf.Cell(nil, req.JobTitle)
	}

	// ── Ligne or sous le nom ─────────────────────────────
	pdf.SetLineWidth(1.0)
	pdf.SetStrokeColor(eGR, eGG, eGB)
	pdf.Line(nameX, 67, pageW-marginR, 67)

	// ── Contacts inline ──────────────────────────────────
	contacts := []string{}
	if req.Phone != "" {
		contacts = append(contacts, req.Phone)
	}
	if req.Email != "" {
		contacts = append(contacts, req.Email)
	}
	if req.Location != "" {
		contacts = append(contacts, req.Location)
	}
	if req.LinkedIn != "" {
		contacts = append(contacts, truncate(req.LinkedIn, 30))
	}
	if req.Github != "" {
		contacts = append(contacts, req.Github)
	}

	if len(contacts) > 0 {
		cx := nameX
		for i, c := range contacts {
			if i > 0 {
				// séparateur or
				pdf.SetFont("regular", "", 8)
				pdf.SetTextColor(eGR, eGG, eGB)
				pdf.SetX(cx)
				pdf.SetY(74)
				pdf.Cell(nil, " | ")
				cx += 12
			}
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(200, 200, 210)
			pdf.SetX(cx)
			pdf.SetY(74)
			pdf.Cell(nil, c)
			cx += float64(len(c))*4.2 + 2
		}
	}

	y := headerH + 18

	// ════════════════════════════════════════════════════
	// LAYOUT 2 COLONNES
	// ════════════════════════════════════════════════════
	leftW := contentW * 0.62
	rightW := contentW * 0.35
	leftX := marginL
	rightX := marginL + leftW + 14

	// ── Helper section ───────────────────────────────────
	sectionHeader := func(title string, x, yPos, width float64) float64 {
		// Fond crème
		pdf.SetFillColor(eCrR, eCrG, eCrB)
		pdf.RectFromUpperLeftWithStyle(x, yPos, width, 16, "F")
		// Barre or gauche
		pdf.SetFillColor(eGR, eGG, eGB)
		pdf.RectFromUpperLeftWithStyle(x, yPos, 4, 16, "F")
		pdf.SetFont("bold", "", 8.5)
		pdf.SetTextColor(eNR, eNG, eNB)
		pdf.SetX(x + 10)
		pdf.SetY(yPos + 3)
		pdf.Cell(nil, strings.ToUpper(title))
		return yPos + 22
	}

	// ════════════════════════════════════════════════════
	// COLONNE GAUCHE — About + Expérience
	// ════════════════════════════════════════════════════
	leftY := y

	// ── About ────────────────────────────────────────────
	if req.Summary != "" {
		leftY = sectionHeader("Profil Professionnel", leftX, leftY, leftW)
		for _, line := range wrapStr(req.Summary, 70) {
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(eTextR, eTextG, eTextB)
			pdf.SetX(leftX)
			pdf.SetY(leftY)
			pdf.Cell(nil, line)
			leftY += 12
		}
		leftY += 10
	}

	// ── Expériences ─────────────────────────────────────
	if len(req.Experiences) > 0 {
		leftY = sectionHeader("Experience Professionnelle", leftX, leftY, leftW)
		for _, exp := range req.Experiences {
			// Badge or à gauche
			pdf.SetFillColor(eGR, eGG, eGB)
			pdf.RectFromUpperLeftWithStyle(leftX, leftY+2, 2, 8, "F")

			// Entreprise
			pdf.SetFont("bold", "", 9.5)
			pdf.SetTextColor(eNR, eNG, eNB)
			pdf.SetX(leftX + 8)
			pdf.SetY(leftY)
			pdf.Cell(nil, exp.Company)

			// Dates
			dates := exp.StartDate
			if exp.EndDate != "" {
				dates += " — " + exp.EndDate
			} else {
				dates += " — Présent"
			}
			dw := float64(len(dates)) * 4.0
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(eGR, eGG, eGB)
			pdf.SetX(leftX + leftW - dw)
			pdf.SetY(leftY)
			pdf.Cell(nil, dates)

			// Poste
			leftY += 13
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(eGR, eGG, eGB)
			pdf.SetX(leftX + 8)
			pdf.SetY(leftY)
			pdf.Cell(nil, exp.Position)

			// Description
			if exp.Description != "" {
				leftY += 12
				for _, line := range wrapStr(exp.Description, 68) {
					pdf.SetFont("regular", "", 7.5)
					pdf.SetTextColor(eGrayR, eGrayG, eGrayB)
					pdf.SetX(leftX + 8)
					pdf.SetY(leftY)
					pdf.Cell(nil, "· "+line)
					leftY += 11
				}
			} else {
				leftY += 8
			}

			// Séparateur fin
			pdf.SetLineWidth(0.3)
			pdf.SetStrokeColor(eLR, eLG, eLB)
			pdf.Line(leftX+8, leftY+2, leftX+leftW, leftY+2)
			leftY += 10
		}
	}

	// ── Formation ────────────────────────────────────────
	if len(req.Education) > 0 {
		leftY = sectionHeader("Formation", leftX, leftY, leftW)
		for _, edu := range req.Education {
			pdf.SetFillColor(eGR, eGG, eGB)
			pdf.RectFromUpperLeftWithStyle(leftX, leftY+2, 2, 8, "F")

			pdf.SetFont("bold", "", 9)
			pdf.SetTextColor(eNR, eNG, eNB)
			pdf.SetX(leftX + 8)
			pdf.SetY(leftY)
			pdf.Cell(nil, edu.Degree)

			if edu.Year > 0 {
				yr := fmt.Sprintf("%d", edu.Year)
				yw := float64(len(yr)) * 4.5
				pdf.SetFont("bold", "", 8)
				pdf.SetTextColor(eGR, eGG, eGB)
				pdf.SetX(leftX + leftW - yw)
				pdf.SetY(leftY)
				pdf.Cell(nil, yr)
			}

			leftY += 13
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(eGrayR, eGrayG, eGrayB)
			pdf.SetX(leftX + 8)
			pdf.SetY(leftY)
			pdf.Cell(nil, edu.School)
			leftY += 18
		}
	}

	// ════════════════════════════════════════════════════
	// COLONNE DROITE — Skills + Langues + Certifs
	// ════════════════════════════════════════════════════
	rightY := y

	// ── Skills ───────────────────────────────────────────
	if len(req.Skills) > 0 {
		rightY = sectionHeader("Competences", rightX, rightY, rightW)
		percents := []int{90, 85, 80, 75, 88, 72, 94, 68, 82, 78, 86, 76}
		for i, skill := range req.Skills {
			pct := percents[i%len(percents)]
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(eTextR, eTextG, eTextB)
			pdf.SetX(rightX)
			pdf.SetY(rightY)
			pdf.Cell(nil, truncate(skill, 18))

			// Barre de progression or
			barY := rightY + 10
			barH := 4.0
			pdf.SetFillColor(eLR, eLG, eLB)
			pdf.RectFromUpperLeftWithStyle(rightX, barY, rightW, barH, "F")
			pdf.SetFillColor(eGR, eGG, eGB)
			pdf.RectFromUpperLeftWithStyle(rightX, barY, rightW*float64(pct)/100, barH, "F")

			pdf.SetFont("regular", "", 6.5)
			pdf.SetTextColor(eGrayR, eGrayG, eGrayB)
			pdf.SetX(rightX + rightW + 2)
			pdf.SetY(rightY + 7)
			pdf.Cell(nil, fmt.Sprintf("%d%%", pct))

			rightY += 22
		}
		rightY += 4
	}

	// ── Langues ──────────────────────────────────────────
	if len(req.Languages) > 0 {
		rightY = sectionHeader("Langues", rightX, rightY, rightW)
		for _, lang := range req.Languages {
			pdf.SetFillColor(eNR, eNG, eNB)
			pdf.RectFromUpperLeftWithStyle(rightX, rightY-1, float64(len(lang))*4.5+14, 13, "F")
			pdf.SetFillColor(eGR, eGG, eGB)
			pdf.RectFromUpperLeftWithStyle(rightX, rightY-1, 3, 13, "F")
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetX(rightX + 7)
			pdf.SetY(rightY)
			pdf.Cell(nil, lang)
			rightY += 17
		}
		rightY += 4
	}

	// ── Certifications ───────────────────────────────────
	if len(req.Certifications) > 0 {
		rightY = sectionHeader("Certifications", rightX, rightY, rightW)
		for _, cert := range req.Certifications {
			for _, line := range wrapStr(cert, 24) {
				pdf.SetFillColor(eGR, eGG, eGB)
				pdf.RectFromUpperLeftWithStyle(rightX, rightY+4, 5, 5, "F")
				pdf.SetFont("regular", "", 7.5)
				pdf.SetTextColor(eTextR, eTextG, eTextB)
				pdf.SetX(rightX + 10)
				pdf.SetY(rightY)
				pdf.Cell(nil, line)
				rightY += 14
			}
		}
	}

	// ── Séparateur vertical entre colonnes ───────────────
	maxH := leftY
	if rightY > maxH {
		maxH = rightY
	}
	pdf.SetLineWidth(0.5)
	pdf.SetStrokeColor(eLR, eLG, eLB)
	pdf.Line(rightX-7, y, rightX-7, maxH)

	return pdf.GetBytesPdf(), nil
}
