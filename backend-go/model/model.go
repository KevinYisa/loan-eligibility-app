package model

// backend-go/model/model.go

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
)

type Weights struct {
	Bias            float64            `json:"bias"`
	Coefficients    map[string]float64 `json:"coefficients"`
	CoefficientsArr []float64          // not in JSON, will be built from map
}

type ScalerParams struct {
	Mean  []float64 `json:"mean"`
	Scale []float64 `json:"scale"`
}

type FeatureImportance struct {
	Feature    string  `json:"feature"`
	Importance float64 `json:"importance"`
}

type Metrics struct {
	Accuracy        float64 `json:"accuracy"`
	Precision       float64 `json:"precision"`
	Recall          float64 `json:"recall"`
	F1              float64 `json:"f1"`
	RocAuc          float64 `json:"roc_auc"`
	ConfusionMatrix [][]int `json:"confusion_matrix"`
}

var (
	instance *Model
	once     sync.Once
)

type Model struct {
	FeatureNames      []string
	Weights           *Weights
	Scaler            *ScalerParams
	FeatureImportance []FeatureImportance
	Metrics           *Metrics
	mu                sync.RWMutex
}

// GetModel loads everything once (singleton)
func GetModel(modelPath string) (*Model, error) {
	var err error
	once.Do(func() {
		instance, err = loadModel(modelPath)
	})
	return instance, err
}

func loadModel(modelPath string) (*Model, error) {
	model := &Model{}

	// 1. Load Feature Names
	data, err := os.ReadFile(filepath.Join(modelPath, "feature_names.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read feature_names.json: %v", err)
	}
	if err := json.Unmarshal(data, &model.FeatureNames); err != nil {
		return nil, fmt.Errorf("failed to parse feature_names.json: %v", err)
	}

	// 2. Load Weights
	data, err = os.ReadFile(filepath.Join(modelPath, "weights.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read weights.json: %v", err)
	}
	var weights Weights
	if err := json.Unmarshal(data, &weights); err != nil {
		return nil, fmt.Errorf("failed to parse weights.json: %v", err)
	}
	// Build coefficient array in the exact same order as FeatureNames
	weights.CoefficientsArr = make([]float64, len(model.FeatureNames))
	for i, name := range model.FeatureNames {
		if val, ok := weights.Coefficients[name]; ok {
			weights.CoefficientsArr[i] = val
		} else {
			return nil, fmt.Errorf("weight for feature %s not found in weights.json", name)
		}
	}
	model.Weights = &weights

	// 3. Load Scaler Params
	data, err = os.ReadFile(filepath.Join(modelPath, "scaler_params.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read scaler_params.json: %v", err)
	}
	var scaler ScalerParams
	if err := json.Unmarshal(data, &scaler); err != nil {
		return nil, fmt.Errorf("failed to parse scaler_params.json: %v", err)
	}
	model.Scaler = &scaler

	// 4. Load Feature Importance (for frontend)
	data, err = os.ReadFile(filepath.Join(modelPath, "feature_importance.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read feature_importance.json: %v", err)
	}
	var importance []FeatureImportance
	if err := json.Unmarshal(data, &importance); err != nil {
		return nil, fmt.Errorf("failed to parse feature_importance.json: %v", err)
	}
	model.FeatureImportance = importance

	// 5. Load Metrics (for frontend)
	data, err = os.ReadFile(filepath.Join(modelPath, "metrics.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read metrics.json: %v", err)
	}
	var metrics Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, fmt.Errorf("failed to parse metrics.json: %v", err)
	}
	model.Metrics = &metrics

	return model, nil
}

// Predict takes raw input features (in the same order as FeatureNames) and returns probability.
func (m *Model) Predict(input []float64) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(input) != len(m.FeatureNames) {
		return 0, fmt.Errorf("expected %d features, got %d", len(m.FeatureNames), len(input))
	}

	// 1. Apply Standard Scaler: (x - mean) / scale
	scaled := make([]float64, len(input))
	for i := 0; i < len(input); i++ {
		scaled[i] = (input[i] - m.Scaler.Mean[i]) / m.Scaler.Scale[i]
	}

	// 2. Linear combination: bias + dot(weights, scaled_features)
	var linear float64 = m.Weights.Bias
	for i := 0; i < len(scaled); i++ {
		linear += m.Weights.CoefficientsArr[i] * scaled[i]
	}

	// 3. Sigmoid: 1 / (1 + exp(-linear))
	probability := 1.0 / (1.0 + math.Exp(-linear))
	return probability, nil
}

// GetFeatures returns the feature importance list
func (m *Model) GetFeatures() []FeatureImportance {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.FeatureImportance
}

// GetMetrics returns the model metrics
func (m *Model) GetMetrics() *Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Metrics
}
