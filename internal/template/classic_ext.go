package templates

import (
	"fmt"
	"strings"

	"github.com/dacostaabraham/cv-generator/internal/models"
	"github.com/signintech/gopdf"
)

// ClassicTemplateExt — version étendue du template Classic
// Supporte les nouveaux champs (JobTitle, LinkedIn, GitHub, etc.)
type ClassicTemplateExt struct{}

func (t *ClassicTemplateExt) TemplateName() string { return "classic" }

func (t *ClassicTemplateExt) RenderExtended(req *models.ExtendedCVRequest) ([]byte, error) {
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
		cNavyR, cNavyG, cNavyB   = 28, 54, 120
		cLightR, cLightG, cLightB = 232, 240, 252
		cAccR, cAccG, cAccB       = 59, 130, 246
		cDarkR, cDarkG, cDarkB   = 15, 23, 42
		cGrayR, cGrayG, cGrayB   = 100, 116, 139
		cLineR, cLineG, cLineB   = 203, 213, 225
	)

	pageW := 595.0
	margin := 40.0
	headerH := 135.0

	// ── Header fond bleu clair ──────────────────────────
	pdf.SetFillColor(cLightR, cLightG, cLightB)
	pdf.RectFromUpperLeftWithStyle(0, 0, pageW, headerH, "F")
	pdf.SetFillColor(cNavyR, cNavyG, cNavyB)
	pdf.RectFromUpperLeftWithStyle(0, headerH-6, pageW, 6, "F")

	// ── Photo ───────────────────────────────────────────
	renderPhotoExt(pdf, req.Photo, 28, 14, 96, 100,
		cLightR, cLightG, cLightB, cNavyR, cNavyG, cNavyB)

	// ── Nom ─────────────────────────────────────────────
	pdf.SetFont("bold", "", 24)
	pdf.SetTextColor(cNavyR, cNavyG, cNavyB)
	pdf.SetX(140)
	pdf.SetY(18)
	pdf.Cell(nil, req.FullName)

	// ── Job Title ───────────────────────────────────────
	if req.JobTitle != "" {
		pdf.SetFont("regular", "", 10)
		pdf.SetTextColor(cAccR, cAccG, cAccB)
		pdf.SetX(140)
		pdf.SetY(46)
		pdf.Cell(nil, req.JobTitle)
	}

	// ── Ligne ───────────────────────────────────────────
	pdf.SetLineWidth(0.8)
	pdf.SetStrokeColor(cAccR, cAccG, cAccB)
	pdf.Line(140, 60, pageW-margin, 60)

	// ── Contacts ────────────────────────────────────────
	contacts := []struct{ label, value string }{
		{"Tel", req.Phone},
		{"Email", req.Email},
		{"Lieu", req.Location},
	}
	if req.LinkedIn != "" {
		contacts = append(contacts, struct{ label, value string }{"LinkedIn", truncate(req.LinkedIn, 35)})
	}
	if req.Github != "" {
		contacts = append(contacts, struct{ label, value string }{"GitHub", req.Github})
	}
	if req.Website != "" {
		contacts = append(contacts, struct{ label, value string }{"Web", truncate(req.Website, 35)})
	}

	cy := 68.0
	colBreak := 3
	col2X := 310.0
	col2Y := 68.0
	for i, c := range contacts {
		if c.value == "" {
			continue
		}
		if i < colBreak {
			pdf.SetFont("bold", "", 7.5)
			pdf.SetTextColor(cNavyR, cNavyG, cNavyB)
			pdf.SetX(140)
			pdf.SetY(cy)
			pdf.Cell(nil, c.label+":")
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(cDarkR, cDarkG, cDarkB)
			pdf.SetX(168)
			pdf.SetY(cy)
			pdf.Cell(nil, c.value)
			cy += 13
		} else {
			pdf.SetFont("bold", "", 7.5)
			pdf.SetTextColor(cNavyR, cNavyG, cNavyB)
			pdf.SetX(col2X)
			pdf.SetY(col2Y)
			pdf.Cell(nil, c.label+":")
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(cDarkR, cDarkG, cDarkB)
			pdf.SetX(col2X+42)
			pdf.SetY(col2Y)
			pdf.Cell(nil, c.value)
			col2Y += 13
		}
	}

	y := headerH + 14

	sectionTitle := func(title string, yPos float64) float64 {
		pdf.SetFillColor(cNavyR, cNavyG, cNavyB)
		pdf.RectFromUpperLeftWithStyle(margin, yPos, 4, 13, "F")
		pdf.SetFont("bold", "", 9.5)
		pdf.SetTextColor(cNavyR, cNavyG, cNavyB)
		pdf.SetX(margin + 10)
		pdf.SetY(yPos)
		pdf.Cell(nil, title)
		pdf.SetLineWidth(0.4)
		pdf.SetStrokeColor(cLineR, cLineG, cLineB)
		pdf.Line(margin, yPos+15, pageW-margin, yPos+15)
		return yPos + 23
	}

	// ── About Me ────────────────────────────────────────
	if req.Summary != "" {
		y = sectionTitle("A PROPOS", y)
		for _, line := range wrapStr(req.Summary, 95) {
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(cGrayR, cGrayG, cGrayB)
			pdf.SetX(margin)
			pdf.SetY(y)
			pdf.Cell(nil, line)
			y += 13
		}
		y += 8
	}

	// ── Expériences ─────────────────────────────────────
	if len(req.Experiences) > 0 {
		y = sectionTitle("EXPERIENCE PROFESSIONNELLE", y)
		for _, exp := range req.Experiences {
			pdf.SetFillColor(cAccR, cAccG, cAccB)
			pdf.RectFromUpperLeftWithStyle(margin, y-2, 3, 10, "F")

			pdf.SetFont("bold", "", 9.5)
			pdf.SetTextColor(cDarkR, cDarkG, cDarkB)
			pdf.SetX(margin + 12)
			pdf.SetY(y)
			pdf.Cell(nil, strings.ToUpper(exp.Company))

			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(cAccR, cAccG, cAccB)
			pdf.SetX(margin + 12)
			pdf.SetY(y + 13)
			pdf.Cell(nil, exp.Position)

			dates := exp.StartDate
			if exp.EndDate != "" {
				dates += " — " + exp.EndDate
			} else {
				dates += " — Présent"
			}
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(cGrayR, cGrayG, cGrayB)
			dateW := float64(len(dates)) * 4.2
			pdf.SetX(pageW - margin - dateW)
			pdf.SetY(y)
			pdf.Cell(nil, dates)

			if exp.Description != "" {
				descY := y + 26
				for _, line := range wrapStr(exp.Description, 85) {
					pdf.SetFont("regular", "", 8)
					pdf.SetTextColor(cGrayR, cGrayG, cGrayB)
					pdf.SetX(margin + 12)
					pdf.SetY(descY)
					pdf.Cell(nil, "· "+line)
					descY += 11
				}
				y = descY + 8
			} else {
				y += 30
			}
			pdf.SetLineWidth(0.3)
			pdf.SetStrokeColor(cLineR, cLineG, cLineB)
			pdf.Line(margin+12, y-4, pageW-margin, y-4)
		}
		y += 4
	}

	// ── Formation ────────────────────────────────────────
	if len(req.Education) > 0 {
		y = sectionTitle("FORMATION", y)
		for _, edu := range req.Education {
			pdf.SetFillColor(cNavyR, cNavyG, cNavyB)
			pdf.RectFromUpperLeftWithStyle(margin, y-2, 3, 10, "F")

			pdf.SetFont("bold", "", 9.5)
			pdf.SetTextColor(cDarkR, cDarkG, cDarkB)
			pdf.SetX(margin + 12)
			pdf.SetY(y)
			pdf.Cell(nil, edu.Degree)

			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(cGrayR, cGrayG, cGrayB)
			pdf.SetX(margin + 12)
			pdf.SetY(y + 13)
			pdf.Cell(nil, edu.School)

			if edu.Year > 0 {
				yr := fmt.Sprintf("%d", edu.Year)
				yrW := float64(len(yr)) * 4.5
				pdf.SetFont("bold", "", 8)
				pdf.SetTextColor(cAccR, cAccG, cAccB)
				pdf.SetX(pageW - margin - yrW)
				pdf.SetY(y)
				pdf.Cell(nil, yr)
			}
			y += 30
		}
		y += 4
	}

	// ── Skills ───────────────────────────────────────────
	if len(req.Skills) > 0 {
		y = sectionTitle("COMPETENCES", y)
		percents := []int{92, 85, 80, 75, 88, 72, 94, 68, 82, 78, 90, 76}
		half := (len(req.Skills) + 1) / 2
		leftY := y
		rightY := y
		for i, skill := range req.Skills {
			pct := percents[i%len(percents)]
			if i < half {
				pdf.SetFillColor(cLightR, cLightG, cLightB)
				pdf.RectFromUpperLeftWithStyle(margin, leftY-1, 115, 14, "F")
				pdf.SetFillColor(cAccR, cAccG, cAccB)
				pdf.RectFromUpperLeftWithStyle(margin, leftY-1, 3, 14, "F")
				pdf.SetFont("regular", "", 8)
				pdf.SetTextColor(cDarkR, cDarkG, cDarkB)
				pdf.SetX(margin + 8)
				pdf.SetY(leftY)
				pdf.Cell(nil, skill)
				leftY += 18
			} else {
				colMid := pageW/2 + 10
				pdf.SetFont("regular", "", 8)
				pdf.SetTextColor(cDarkR, cDarkG, cDarkB)
				pdf.SetX(colMid)
				pdf.SetY(rightY)
				pdf.Cell(nil, skill)
				barX := colMid + 95.0
				barW := 105.0
				pdf.SetFillColor(cLineR, cLineG, cLineB)
				pdf.RectFromUpperLeftWithStyle(barX, rightY+3, barW, 5, "F")
				pdf.SetFillColor(cAccR, cAccG, cAccB)
				pdf.RectFromUpperLeftWithStyle(barX, rightY+3, barW*float64(pct)/100, 5, "F")
				pdf.SetFont("regular", "", 6.5)
				pdf.SetTextColor(cGrayR, cGrayG, cGrayB)
				pdf.SetX(barX + barW + 3)
				pdf.SetY(rightY + 1)
				pdf.Cell(nil, fmt.Sprintf("%d%%", pct))
				rightY += 18
			}
		}
		if leftY > rightY {
			y = leftY
		} else {
			y = rightY
		}
		y += 8
	}

	// ── Langues ──────────────────────────────────────────
	if len(req.Languages) > 0 {
		y = sectionTitle("LANGUES", y)
		lx := margin
		for _, lang := range req.Languages {
			pdf.SetFillColor(cNavyR, cNavyG, cNavyB)
			pdf.RectFromUpperLeftWithStyle(lx, y-1, float64(len(lang))*5+16, 14, "F")
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetX(lx + 8)
			pdf.SetY(y)
			pdf.Cell(nil, lang)
			lx += float64(len(lang))*5 + 24
		}
		y += 22
	}

	// ── Certifications ───────────────────────────────────
	if len(req.Certifications) > 0 {
		y = sectionTitle("CERTIFICATIONS", y)
		for _, cert := range req.Certifications {
			pdf.SetFillColor(cAccR, cAccG, cAccB)
			pdf.RectFromUpperLeftWithStyle(margin, y+4, 6, 6, "F")
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(cDarkR, cDarkG, cDarkB)
			pdf.SetX(margin + 14)
			pdf.SetY(y)
			pdf.Cell(nil, cert)
			y += 16
		}
	}

	return pdf.GetBytesPdf(), nil
}
