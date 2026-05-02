package templates

import (
	"fmt"
	"os"
	"strings"

	"github.com/signintech/gopdf"
)

// ════════════════════════════════════════════════════════
// HELPERS PARTAGÉS POUR TOUS LES TEMPLATES ÉTENDUS
// ════════════════════════════════════════════════════════

// renderPhotoExt affiche la photo (helper commun)
func renderPhotoExt(pdf *gopdf.GoPdf, photoBytes []byte, x, y, w, h float64,
	bgR, bgG, bgB, acR, acG, acB uint8) {
	if len(photoBytes) > 0 {
		ext := ".jpg"
		if len(photoBytes) >= 4 &&
			photoBytes[0] == 0x89 && photoBytes[1] == 0x50 &&
			photoBytes[2] == 0x4E && photoBytes[3] == 0x47 {
			ext = ".png"
		}
		tmpFile, err := os.CreateTemp("", "cv-photo-ext-*"+ext)
		if err == nil {
			_, writeErr := tmpFile.Write(photoBytes)
			tmpFile.Close()
			if writeErr == nil {
				imgErr := pdf.Image(tmpFile.Name(), x, y, &gopdf.Rect{W: w, H: h})
				os.Remove(tmpFile.Name())
				if imgErr == nil {
					return
				}
				fmt.Printf("[warn] photo ext: %v\n", imgErr)
			} else {
				os.Remove(tmpFile.Name())
			}
		}
	}
	// Placeholder
	pdf.SetFillColor(bgR, bgG, bgB)
	pdf.RectFromUpperLeftWithStyle(x, y, w, h, "F")
	pdf.SetFillColor(acR, acG, acB)
	pdf.RectFromUpperLeftWithStyle(x, y+h-3, w, 3, "F")
	pdf.SetFont("regular", "", 7)
	pdf.SetTextColor(180, 180, 180)
	pdf.SetX(x + w/2 - 10)
	pdf.SetY(y + h/2 - 4)
	pdf.Cell(nil, "PHOTO")
}

// wrapStr coupe le texte en lignes de maxLen caractères
func wrapStr(text string, maxLen int) []string {
	if text == "" {
		return nil
	}
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

// truncate coupe un texte à maxLen caractères
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
