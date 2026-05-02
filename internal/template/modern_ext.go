package templates

import (
	"fmt"
	"strings"

	"github.com/dacostaabraham/cv-generator/internal/models"
	"github.com/signintech/gopdf"
)

// ModernTemplateExt — version étendue du template Modern
type ModernTemplateExt struct{}

func (t *ModernTemplateExt) TemplateName() string { return "modern" }

func (t *ModernTemplateExt) RenderExtended(req *models.ExtendedCVRequest) ([]byte, error) {
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
		mSR, mSG, mSB   = 15, 23, 42   // sidebar très sombre
		mAR, mAG, mAB   = 99, 102, 241 // violet accent
		mA2R, mA2G, mA2B = 34, 211, 238 // teal accent 2
		mTR, mTG, mTB   = 248, 250, 252 // texte clair sidebar
		mMR, mMG, mMB   = 148, 163, 184 // texte muté sidebar
		mDR, mDG, mDB   = 15, 23, 42   // texte sombre main
		mGR, mGG, mGB   = 100, 116, 139 // gris
		mLR, mLG, mLB   = 226, 232, 240 // ligne
	)

	sideW := 185.0
	mainX := sideW + 22.0
	sideM := 16.0

	// ── Sidebar ─────────────────────────────────────────
	pdf.SetFillColor(mSR, mSG, mSB)
	pdf.RectFromUpperLeftWithStyle(0, 0, sideW, 842, "F")

	// Bande accent haut
	pdf.SetFillColor(mAR, mAG, mAB)
	pdf.RectFromUpperLeftWithStyle(0, 0, sideW, 5, "F")

	// ── Photo ───────────────────────────────────────────
	photoSize := 78.0
	photoX := (sideW - photoSize) / 2
	renderPhotoExt(pdf, req.Photo, photoX, 18, photoSize, photoSize,
		30, 41, 59, mAR, mAG, mAB)

	// Ligne déco sous photo
	pdf.SetFillColor(mAR, mAG, mAB)
	pdf.RectFromUpperLeftWithStyle(sideM, 102, sideW-sideM*2, 2, "F")

	// ── Nom ─────────────────────────────────────────────
	pdf.SetFont("bold", "", 13)
	pdf.SetTextColor(mTR, mTG, mTB)
	nameLines := wrapStr(req.FullName, 20)
	nameY := 110.0
	for _, line := range nameLines {
		tw := float64(len(line)) * 6.5
		if tw > sideW {
			tw = sideW - sideM
		}
		pdf.SetX((sideW - tw) / 2)
		pdf.SetY(nameY)
		pdf.Cell(nil, line)
		nameY += 17
	}

	// ── Job Title ───────────────────────────────────────
	if req.JobTitle != "" {
		pdf.SetFont("regular", "", 8)
		pdf.SetTextColor(mA2R, mA2G, mA2B)
		jtLines := wrapStr(req.JobTitle, 24)
		for _, line := range jtLines {
			tw := float64(len(line)) * 4.5
			pdf.SetX((sideW - tw) / 2)
			pdf.SetY(nameY)
			pdf.Cell(nil, line)
			nameY += 12
		}
		nameY += 2
	}

	sideY := nameY + 6

	// ── Section helper sidebar ───────────────────────────
	sideSection := func(title string, y float64) float64 {
		pdf.SetFillColor(mA2R, mA2G, mA2B)
		pdf.RectFromUpperLeftWithStyle(sideM, y, 20, 2, "F")
		y += 6
		pdf.SetFont("bold", "", 7.5)
		pdf.SetTextColor(mA2R, mA2G, mA2B)
		pdf.SetX(sideM)
		pdf.SetY(y)
		pdf.Cell(nil, title)
		return y + 13
	}

	// ── Contacts ────────────────────────────────────────
	sideY = sideSection("CONTACT", sideY)
	contacts := []struct{ icon, val string }{
		{"@", req.Email},
		{"T", req.Phone},
		{"L", req.Location},
	}
	if req.LinkedIn != "" {
		contacts = append(contacts, struct{ icon, val string }{"in", truncate(req.LinkedIn, 22)})
	}
	if req.Github != "" {
		contacts = append(contacts, struct{ icon, val string }{"gh", req.Github})
	}
	if req.Website != "" {
		contacts = append(contacts, struct{ icon, val string }{"W", truncate(req.Website, 22)})
	}

	for _, c := range contacts {
		if c.val == "" {
			continue
		}
		pdf.SetFillColor(mAR, mAG, mAB)
		pdf.RectFromUpperLeftWithStyle(sideM, sideY-1, 12, 12, "F")
		pdf.SetFont("bold", "", 6.5)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetX(sideM + 1)
		pdf.SetY(sideY)
		pdf.Cell(nil, c.icon)

		pdf.SetFont("regular", "", 7.5)
		pdf.SetTextColor(mTR, mTG, mTB)
		for i, line := range wrapStr(c.val, 22) {
			pdf.SetX(sideM + 16)
			pdf.SetY(sideY + float64(i)*9)
			pdf.Cell(nil, line)
		}
		sideY += 15
	}
	sideY += 6

	// ── Skills sidebar ───────────────────────────────────
	if len(req.Skills) > 0 {
		sideY = sideSection("COMPETENCES", sideY)
		for _, skill := range req.Skills {
			badgeW := float64(len(skill))*4.8 + 12
			if badgeW > sideW-sideM*2 {
				badgeW = sideW - sideM*2
			}
			pdf.SetFillColor(30, 41, 59)
			pdf.RectFromUpperLeftWithStyle(sideM, sideY-1, badgeW, 13, "F")
			pdf.SetFillColor(mAR, mAG, mAB)
			pdf.RectFromUpperLeftWithStyle(sideM, sideY-1, 3, 13, "F")
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(mTR, mTG, mTB)
			pdf.SetX(sideM + 7)
			pdf.SetY(sideY)
			pdf.Cell(nil, truncate(skill, 20))
			sideY += 16
		}
		sideY += 4
	}

	// ── Langues sidebar ──────────────────────────────────
	if len(req.Languages) > 0 {
		sideY = sideSection("LANGUES", sideY)
		for _, lang := range req.Languages {
			pdf.SetFont("regular", "", 7.5)
			pdf.SetTextColor(mMR, mMG, mMB)
			pdf.SetX(sideM)
			pdf.SetY(sideY)
			pdf.Cell(nil, "▸ "+lang)
			sideY += 13
		}
		sideY += 4
	}

	// ── Certifications sidebar ───────────────────────────
	if len(req.Certifications) > 0 {
		sideY = sideSection("CERTIFICATIONS", sideY)
		for _, cert := range req.Certifications {
			for _, line := range wrapStr(cert, 22) {
				pdf.SetFont("regular", "", 7)
				pdf.SetTextColor(mMR, mMG, mMB)
				pdf.SetX(sideM)
				pdf.SetY(sideY)
				pdf.Cell(nil, "· "+line)
				sideY += 11
			}
		}
	}

	// ════════════════════════════════════════════════════
	// COLONNE PRINCIPALE
	// ════════════════════════════════════════════════════
	mainY := 20.0

	mainSection := func(title string, y float64) float64 {
		pdf.SetFillColor(mAR, mAG, mAB)
		pdf.RectFromUpperLeftWithStyle(mainX, y, 3, 13, "F")
		pdf.SetFont("bold", "", 9.5)
		pdf.SetTextColor(mAR, mAG, mAB)
		pdf.SetX(mainX + 9)
		pdf.SetY(y)
		pdf.Cell(nil, title)
		pdf.SetLineWidth(0.4)
		pdf.SetStrokeColor(mLR, mLG, mLB)
		pdf.Line(mainX, y+15, 575, y+15)
		return y + 23
	}

	// ── Profil ───────────────────────────────────────────
	if req.Summary != "" {
		mainY = mainSection("PROFIL", mainY)
		for _, line := range wrapStr(req.Summary, 52) {
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(mGR, mGG, mGB)
			pdf.SetX(mainX)
			pdf.SetY(mainY)
			pdf.Cell(nil, line)
			mainY += 12
		}
		mainY += 8
	}

	// ── Expériences ─────────────────────────────────────
	if len(req.Experiences) > 0 {
		mainY = mainSection("EXPERIENCE", mainY)
		for _, exp := range req.Experiences {
			pdf.SetFillColor(mAR, mAG, mAB)
			pdf.RectFromUpperLeftWithStyle(mainX, mainY+3, 8, 8, "F")

			pdf.SetFont("bold", "", 9.5)
			pdf.SetTextColor(mDR, mDG, mDB)
			pdf.SetX(mainX + 14)
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
			pdf.SetTextColor(mAR, mAG, mAB)
			pdf.SetX(575 - dw)
			pdf.SetY(mainY + 1)
			pdf.Cell(nil, dates)

			mainY += 13
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(mAR, mAG, mAB)
			pdf.SetX(mainX + 14)
			pdf.SetY(mainY)
			pdf.Cell(nil, exp.Position)

			if exp.Description != "" {
				mainY += 12
				for _, line := range wrapStr(exp.Description, 52) {
					pdf.SetFont("regular", "", 8)
					pdf.SetTextColor(mGR, mGG, mGB)
					pdf.SetX(mainX + 14)
					pdf.SetY(mainY)
					pdf.Cell(nil, "· "+line)
					mainY += 11
				}
			} else {
				mainY += 8
			}
			pdf.SetLineWidth(0.3)
			pdf.SetStrokeColor(mLR, mLG, mLB)
			pdf.Line(mainX+4, mainY+4, 575, mainY+4)
			mainY += 12
		}
	}

	// ── Formation ────────────────────────────────────────
	if len(req.Education) > 0 {
		mainY = mainSection("FORMATION", mainY)
		for _, edu := range req.Education {
			pdf.SetFillColor(mAR, mAG, mAB)
			pdf.RectFromUpperLeftWithStyle(mainX, mainY+3, 8, 8, "F")

			pdf.SetFont("bold", "", 9.5)
			pdf.SetTextColor(mDR, mDG, mDB)
			pdf.SetX(mainX + 14)
			pdf.SetY(mainY)
			pdf.Cell(nil, edu.Degree)

			if edu.Year > 0 {
				yr := fmt.Sprintf("%d", edu.Year)
				yw := float64(len(yr)) * 5.0
				pdf.SetFont("bold", "", 8)
				pdf.SetTextColor(mAR, mAG, mAB)
				pdf.SetX(575 - yw)
				pdf.SetY(mainY + 1)
				pdf.Cell(nil, yr)
			}

			mainY += 13
			pdf.SetFont("regular", "", 8.5)
			pdf.SetTextColor(mGR, mGG, mGB)
			pdf.SetX(mainX + 14)
			pdf.SetY(mainY)
			pdf.Cell(nil, edu.School)
			mainY += 18
		}
	}

	return pdf.GetBytesPdf(), nil
}
