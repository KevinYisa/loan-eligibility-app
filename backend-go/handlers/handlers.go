package handlers

// backend-go/handlers/handlers.go

import (
	"encoding/json"
	"loan-backend/model"
	"net/http"
)

type PredictRequest struct {
	NoOfDependents int     `json:"no_of_dependents"`
	Education      int     `json:"education"`     // 1 = Graduate, 0 = Not Graduate
	SelfEmployed   int     `json:"self_employed"` // 1 = Yes, 0 = No
	IncomeAnnum    float64 `json:"income_annum"`
	LoanAmount     float64 `json:"loan_amount"`
	LoanTerm       int     `json:"loan_term"`
	CibilScore     float64 `json:"cibil_score"`
	TotalAssets    float64 `json:"total_assets"`
}

type PredictResponse struct {
	Eligible    bool    `json:"eligible"`
	Probability float64 `json:"probability"`
	Risk        string  `json:"risk"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func MetricsHandler(m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m.GetMetrics())
	}
}

func FeaturesHandler(m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m.GetFeatures())
	}
}

func PredictHandler(m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req PredictRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Build input slice in the exact order of model.FeatureNames
		// Order: no_of_dependents, education, self_employed, income_annum, loan_amount, loan_term, cibil_score, total_assets
		input := []float64{
			float64(req.NoOfDependents),
			float64(req.Education),
			float64(req.SelfEmployed),
			req.IncomeAnnum,
			req.LoanAmount,
			float64(req.LoanTerm),
			req.CibilScore,
			req.TotalAssets,
		}

		prob, err := m.Predict(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		eligible := prob >= 0.5
		risk := "High"
		if prob >= 0.7 {
			risk = "Low"
		} else if prob >= 0.5 {
			risk = "Medium"
		}

		resp := PredictResponse{
			Eligible:    eligible,
			Probability: prob,
			Risk:        risk,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
