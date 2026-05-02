package templates

import (
	"fmt"
	"strings"

	"github.com/dacostaabraham/cv-generator/internal/models"
	"github.com/signintech/gopdf"
)

// CreativeTemplate — Bold, violet/indigo, header gradient simulé
// Palette: #4f46e5 (indigo), #7c3aed (violet), #f97316 (orange), #1e1b4b (dark)
type CreativeTemplate struct{}

func (t *CreativeTemplate) TemplateName() string { return "creative" }

func (t *CreativeTemplate) RenderExtended(req *models.ExtendedCVRequest) ([]byte, error) {
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
		cIndR, cIndG, cIndB     = 79, 70, 229   // #4f46e5 indigo
		cVioR, cVioG, cVioB     = 109, 40, 217  // #6d28d9 violet
		cDark2R, cDark2G, cDark2B = 30, 27, 75  // #1e1b4b
		cOrgR, cOrgG, cOrgB     = 249, 115, 22  // #f97316 orange
		cTextR, cTextG, cTextB  = 30, 27, 75    // texte principal
		cGrayR, cGrayG, cGrayB  = 107, 114, 128 // gris
		cLR, cLG, cLB           = 229, 231, 235 // ligne
		cLightR, cLightG, cLightB = 238, 242, 255 // fond très léger indigo
	)

	pageW := 595.0

	// ════════════════════════════════════════════════════
	// HEADER — gradient simulé indigo->violet
	// (On simule avec plusieurs rectangles de largeur dégressive)
	// ════════════════════════════════════════════════════
	headerH := 145.0

	// Fond de base indigo
	pdf.SetFillColor(cIndR, cIndG, cIndB)
	pdf.RectFromUpperLeftWithStyle(0, 0, pageW, headerH, "F")

	// Couches violet pour simuler gradient (droite)
	steps := 12
	for i := 0; i < steps; i++ {
		frac := float64(i) / float64(steps)
		x := pageW * frac
		w := pageW * (1 - frac)
		alpha := uint8(float64(cVioR-cIndR)*frac) + cIndR
		alphag := uint8(float64(cVioG-cIndG)*frac) + cIndG
		alphab := uint8(float64(cVioB-cIndB)*frac) + cIndB
		pdf.SetFillColor(alpha, alphag, alphab)
		pdf.RectFromUpperLeftWithStyle(x, 0, w, headerH, "F")
	}

	// Accent orange en bas du header
	pdf.SetFillColor(cOrgR, cOrgG, cOrgB)
	pdf.RectFromUpperLeftWithStyle(0, headerH-4, pageW, 4, "F")

	// ── Photo ───────────────────────────────────────────
	photoSize := 88.0
	photoX := 28.0
	renderPhotoExt(pdf, req.Photo, photoX, (headerH-photoSize)/2, photoSize, photoSize,
		60, 50, 120, cOrgR, cOrgG, cOrgB)

	// ── Nom ─────────────────────────────────────────────
	nameX := photoX + photoSize + 20
	pdf.SetFont("bold", "", 28)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetX(nameX)
	pdf.SetY(20)
	pdf.Cell(nil, req.FullName)

	// ── Job Title ───────────────────────────────────────
	titleY := 52.0
	if req.JobTitle != "" {
		pdf.SetFont("regular", "", 12)
		pdf.SetTextColor(cOrgR, cOrgG, cOrgB)
		pdf.SetX(nameX)
		pdf.SetY(titleY)
		pdf.Cell(nil, req.JobTitle)
		titleY = 70.0
	}

	// ── Ligne orange ─────────────────────────────────────
	pdf.SetLineWidth(1.5)
	pdf.SetStrokeColor(cOrgR, cOrgG, cOrgB)
	pdf.Line(nameX, titleY, nameX+160, titleY)

	// ── Contacts ─────────────────────────────────────────
	contacts := []struct{ val string }{
		{req.Phone}, {req.Email}, {req.Location},
	}
	if req.LinkedIn != "" {
		contacts = append(contacts, struct{ val string }{truncate(req.LinkedIn, 28)})
	}
	if req.Github != "" {
		contacts = append(contacts, struct{ val string }{req.Github})
	}

	cx := nameX
	ctY := titleY + 10
	for i, c := range contacts {
		if c.val == "" {
			continue
		}
		if i > 0 && cx > nameX {
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(200, 200, 255)
			pdf.SetX(cx)
			pdf.SetY(ctY)
			pdf.Cell(nil, "  ·  ")
			cx += 20
		}
		pdf.SetFont("regular", "", 8)
		pdf.SetTextColor(220, 220, 255)
		pdf.SetX(cx)
		pdf.SetY(ctY)
		pdf.Cell(nil, c.val)
		cx += float64(len(c.val))*4.0 + 2
		if cx > pageW-80 {
			cx = nameX
			ctY += 13
		}
	}

	y := headerH + 16

	// ════════════════════════════════════════════════════
	// LAYOUT : sidebar gauche (38%) + main droite (58%)
	// ════════════════════════════════════════════════════
	sideW := 205.0
	mainX := sideW + 18.0
	sideM := 18.0

	sideSection := func(title string, yPos float64) float64 {
		pdf.SetFillColor(cLightR, cLightG, cLightB)
		pdf.RectFromUpperLeftWithStyle(0, yPos, sideW, 16, "F")
		pdf.SetFillColor(cOrgR, cOrgG, cOrgB)
		pdf.RectFromUpperLeftWithStyle(0, yPos, 4, 16, "F")
		pdf.SetFont("bold", "", 8)
		pdf.SetTextColor(cDark2R, cDark2G, cDark2B)
		pdf.SetX(sideM)
		pdf.SetY(yPos + 3)
		pdf.Cell(nil, strings.ToUpper(title))
		return yPos + 22
	}

	mainSection := func(title string, yPos float64) float64 {
		pdf.SetFillColor(cIndR, cIndG, cIndB)
		pdf.RectFromUpperLeftWithStyle(mainX-4, yPos, 4, 14, "F")
		pdf.SetFont("bold", "", 10)
		pdf.SetTextColor(cIndR, cIndG, cIndB)
		pdf.SetX(mainX + 4)
		pdf.SetY(yPos)
		pdf.Cell(nil, strings.ToUpper(title))
		pdf.SetLineWidth(0.5)
		pdf.SetStrokeColor(cLightR, cLightG, cLightB)
		pdf.Line(mainX+4, yPos+15, pageW-20, yPos+15)
		return yPos + 23
	}

	// ════════════════════════════════════════════════════
	// SIDEBAR GAUCHE
	// ════════════════════════════════════════════════════
	sideY := y

	// Skills
	if len(req.Skills) > 0 {
		sideY = sideSection("Competences", sideY)
		percents := []int{92, 85, 88, 75, 80, 72, 94, 68, 82, 78, 90, 76}
		for i, skill := range req.Skills {
			pct := percents[i%len(percents)]
			pdf.SetFont("regular", "", 8)
			pdf.SetTextColor(cTextR, cTextG, cTextB)
			pdf.SetX(sideM)
			pdf.SetY(sideY)
			pdf.Cell(nil, truncate(skill, 22))

			barY := sideY + 10
			pdf.SetFillColor(cLR, cLG, cLB)
			pdf.RectFromUpperLeftWithStyle(sideM, barY, sideW-sideM*2, 4, "F")
			pdf.SetFillColor(cIndR, cIndG, cIndB)
			pdf.RectFromUpperLeftWithStyle(sideM, barY, (sideW-sideM*2)*float64(pct)/100, 4, "F")
			sideY += 20
		}
		sideY += 6
	}

	// Langues
	if len(req.Languages) > 0 {
		sideY = sideSection("Langues", sideY)
		for _, lang := range req.Languages {
			// Badge indigo
			lw := float64(len(lang))*4.5 + 14
			pdf.SetFillColor(cIndR, cIndG, cIndB)
			pdf.RectFromUpperLeftWithStyle(sideM, sideY-1, lw, 13, "F")
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetX(sideM + 7)
			pdf.SetY(sideY)
			pdf.Cell(nil, lang)
			sideY += 17
		}
		sideY += 4
	}

	// Certifications
	if len(req.Certifications) > 0 {
		sideY = sideSection("Certifications", sideY)
		for _, cert := range req.Certifications {
			for _, line := range wrapStr(cert, 24) {
				pdf.SetFillColor(cOrgR, cOrgG, cOrgB)
				pdf.RectFromUpperLeftWithStyle(sideM, sideY+5, 5, 5, "F")
				pdf.SetFont("regular", "", 7.5)
				pdf.SetTextColor(cTextR, cTextG, cTextB)
				pdf.SetX(sideM + 10)
				pdf.SetY(sideY)
				pdf.Cell(nil, line)
				sideY += 14
			}
		}
	}

	// ════════════════════════════════════════════════════
	// COLONNE PRINCIPALE
	// ════════════════════════════════════════════════════
	mainY := y

	// About
	if req.Summary != "" {
		mainY = mainSection("Profil", mainY)
		for _, line := range wrapStr(req.Summary, 58) {
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(cGrayR, cGrayG, cGrayB)
			pdf.SetX(mainX + 4)
			pdf.SetY(mainY)
			pdf.Cell(nil, line)
			mainY += 12
		}
		mainY += 10
	}

	// Expériences
	if len(req.Experiences) > 0 {
		mainY = mainSection("Experience", mainY)
		for _, exp := range req.Experiences {
			// Dot orange
			pdf.SetFillColor(cOrgR, cOrgG, cOrgB)
			pdf.RectFromUpperLeftWithStyle(mainX-4, mainY+4, 8, 8, "F")

			pdf.SetFont("bold", "", 9.5)
			pdf.SetTextColor(cDark2R, cDark2G, cDark2B)
			pdf.SetX(mainX + 9)
			pdf.SetY(mainY)
			pdf.Cell(nil, strings.ToUpper(exp.Company))

			dates := exp.StartDate
			if exp.EndDate != "" {
				dates += " — " + exp.EndDate
			} else {
				dates += " — Présent"
			}
			dw := float64(len(dates)) * 4.0
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(cIndR, cIndG, cIndB)
			pdf.SetX(pageW - 20 - dw)
			pdf.SetY(mainY)
			pdf.Cell(nil, dates)

			mainY += 13
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(cIndR, cIndG, cIndB)
			pdf.SetX(mainX + 9)
			pdf.SetY(mainY)
			pdf.Cell(nil, exp.Position)

			if exp.Description != "" {
				mainY += 12
				for _, line := range wrapStr(exp.Description, 56) {
					pdf.SetFont("regular", "", 8)
					pdf.SetTextColor(cGrayR, cGrayG, cGrayB)
					pdf.SetX(mainX + 9)
					pdf.SetY(mainY)
					pdf.Cell(nil, "· "+line)
					mainY += 11
				}
			} else {
				mainY += 8
			}
			pdf.SetLineWidth(0.3)
			pdf.SetStrokeColor(cLightR, cLightG, cLightB)
			pdf.Line(mainX+4, mainY+4, pageW-20, mainY+4)
			mainY += 12
		}
	}

	// Formation
	if len(req.Education) > 0 {
		mainY = mainSection("Formation", mainY)
		for _, edu := range req.Education {
			pdf.SetFillColor(cOrgR, cOrgG, cOrgB)
			pdf.RectFromUpperLeftWithStyle(mainX-4, mainY+4, 8, 8, "F")

			pdf.SetFont("bold", "", 9.5)
			pdf.SetTextColor(cDark2R, cDark2G, cDark2B)
			pdf.SetX(mainX + 9)
			pdf.SetY(mainY)
			pdf.Cell(nil, edu.Degree)

			if edu.Year > 0 {
				yr := fmt.Sprintf("%d", edu.Year)
				yw := float64(len(yr)) * 4.5
				pdf.SetFont("bold", "", 8)
				pdf.SetTextColor(cOrgR, cOrgG, cOrgB)
				pdf.SetX(pageW - 20 - yw)
				pdf.SetY(mainY)
				pdf.Cell(nil, yr)
			}
			mainY += 13
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(cGrayR, cGrayG, cGrayB)
			pdf.SetX(mainX + 9)
			pdf.SetY(mainY)
			pdf.Cell(nil, edu.School)
			mainY += 18
		}
	}

	// Séparateur vertical sidebar/main
	maxH := sideY
	if mainY > maxH {
		maxH = mainY
	}
	pdf.SetLineWidth(0.5)
	pdf.SetStrokeColor(cLR, cLG, cLB)
	pdf.Line(sideW+1, y-16, sideW+1, maxH+10)

	return pdf.GetBytesPdf(), nil
}
