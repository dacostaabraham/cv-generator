package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// tokenSecret est chargé depuis $TOKEN_SECRET (env var)
// Fallback en dev uniquement — NE PAS utiliser en production sans le définir
func tokenSecret() []byte {
	s := os.Getenv("TOKEN_SECRET")
	if s == "" {
		s = "cv-generator-dev-secret-change-me-in-prod"
	}
	return []byte(s)
}

// generateNonce produit 8 octets aléatoires en hex
func generateNonce() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// CreatePaymentToken crée un token signé HMAC-SHA256 après vérification du paiement.
// Format : base64url(saleID:timestamp:nonce).HMAC
// Valide pendant 1 heure, usage unique côté client (stocké en sessionStorage).
func CreatePaymentToken(saleID string) string {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := generateNonce()
	payload := fmt.Sprintf("%s:%s:%s", saleID, ts, nonce)
	encoded := base64.URLEncoding.EncodeToString([]byte(payload))

	mac := hmac.New(sha256.New, tokenSecret())
	mac.Write([]byte(encoded))
	sig := hex.EncodeToString(mac.Sum(nil))

	return encoded + "." + sig
}

// VerifyPaymentToken valide le token et retourne le saleID si OK.
// Retourne une erreur si expiré (>1h) ou signature invalide.
func VerifyPaymentToken(token string) (string, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return "", errors.New("token malformé")
	}
	encoded, sig := parts[0], parts[1]

	// Vérifier la signature
	mac := hmac.New(sha256.New, tokenSecret())
	mac.Write([]byte(encoded))
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return "", errors.New("signature invalide")
	}

	// Décoder le payload
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.New("payload illisible")
	}
	fields := strings.SplitN(string(decoded), ":", 3)
	if len(fields) != 3 {
		return "", errors.New("payload malformé")
	}
	saleID, tsStr := fields[0], fields[1]

	// Vérifier l'expiration (1 heure)
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return "", errors.New("timestamp invalide")
	}
	if time.Now().Unix()-ts > 3600 {
		return "", errors.New("token expiré")
	}

	return saleID, nil
}
