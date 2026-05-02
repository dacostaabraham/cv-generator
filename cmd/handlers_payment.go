package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ════════════════════════════════════════════════════════
// HANDLER PAIEMENT — Vérification Chariow + émission token
// ════════════════════════════════════════════════════════

const (
	chariowProductID = "prd_4cjim0" // GENERATEUR DE CV — 1,000 FCFA
	chariowAPIBase   = "https://api.chariow.com/v1"
)

// chariowSaleResponse représente la réponse de l'API Chariow pour une vente
type chariowSaleResponse struct {
	Data struct {
		ID     string `json:"id"`
		Status string `json:"status"` // "completed", "pending", "refunded"...
		Product struct {
			ID string `json:"id"`
		} `json:"product"`
		Customer struct {
			Email string `json:"email"`
		} `json:"customer"`
		Amount struct {
			Value    float64 `json:"value"`
			Currency string  `json:"currency"`
		} `json:"amount"`
		CreatedAt string `json:"created_at"`
	} `json:"data"`
	Message string `json:"message"`
}

// PaymentVerifyRequest : corps JSON envoyé par le frontend
type PaymentVerifyRequest struct {
	SaleID string `json:"sale_id"`
}

// PaymentVerifyResponse : réponse renvoyée au frontend
type PaymentVerifyResponse struct {
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// VerifyPayment — POST /api/verify-payment
// Vérifie qu'un sale_id Chariow est valide puis retourne un token signé
func (h *Handler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Lire le corps
	var req PaymentVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SaleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentVerifyResponse{Error: "sale_id manquant"})
		return
	}

	fmt.Printf("[payment] Vérification sale_id: %s\n", req.SaleID)

	// Appeler l'API Chariow
	sale, err := fetchChariowSale(req.SaleID)
	if err != nil {
		fmt.Printf("[payment] Erreur API Chariow: %v\n", err)
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(PaymentVerifyResponse{Error: "Impossible de vérifier le paiement"})
		return
	}

	// Vérifier que la vente est bien pour notre produit
	if sale.Data.Product.ID != chariowProductID {
		fmt.Printf("[payment] Produit invalide: %s (attendu: %s)\n", sale.Data.Product.ID, chariowProductID)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(PaymentVerifyResponse{Error: "Ce paiement ne correspond pas au produit CV Generator"})
		return
	}

	// Vérifier que le statut est "completed"
	if sale.Data.Status != "completed" {
		fmt.Printf("[payment] Statut invalide: %s\n", sale.Data.Status)
		w.WriteHeader(http.StatusPaymentRequired)
		json.NewEncoder(w).Encode(PaymentVerifyResponse{Error: fmt.Sprintf("Paiement non finalisé (statut: %s)", sale.Data.Status)})
		return
	}

	// ✅ Paiement valide — générer le token
	token := CreatePaymentToken(req.SaleID)
	fmt.Printf("[payment] ✅ Paiement vérifié pour %s — token émis\n", sale.Data.Customer.Email)

	json.NewEncoder(w).Encode(PaymentVerifyResponse{
		Token:   token,
		Message: fmt.Sprintf("Paiement de %s %s confirmé", formatAmount(sale.Data.Amount.Value), sale.Data.Amount.Currency),
	})
}

// fetchChariowSale interroge l'API Chariow pour vérifier une vente
func fetchChariowSale(saleID string) (*chariowSaleResponse, error) {
	apiKey := os.Getenv("CHARIOW_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CHARIOW_API_KEY non définie")
	}

	url := fmt.Sprintf("%s/sales/%s", chariowAPIBase, saleID)

	client := &http.Client{Timeout: 10 * time.Second}
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("création requête: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("appel API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lecture réponse: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Chariow HTTP %d: %s", resp.StatusCode, string(body))
	}

	var sale chariowSaleResponse
	if err := json.Unmarshal(body, &sale); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return &sale, nil
}

func formatAmount(v float64) string {
	if v == float64(int(v)) {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.2f", v)
}
