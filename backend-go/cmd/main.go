// backend-go/cmd/main.go
package main

import (
	"fmt"
	"loan-backend/config"
	"loan-backend/handlers"
	"loan-backend/model"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.Load()

	// Load model
	mod, err := model.GetModel(cfg.ModelPath)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	log.Println("Model loaded successfully.")

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/predict", handlers.PredictHandler(mod))
	mux.HandleFunc("/features", handlers.FeaturesHandler(mod))
	mux.HandleFunc("/metrics", handlers.MetricsHandler(mod))

	// Wrap with CORS
	handler := enableCORS(mux)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
