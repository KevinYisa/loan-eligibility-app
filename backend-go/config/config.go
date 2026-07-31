// backend-go/config/config.go
package config

import (
	"os"
)

type Config struct {
	Port      string
	ModelPath string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	modelPath := os.Getenv("MODEL_PATH")
	if modelPath == "" {
		// Default: assuming we are in backend-go/, and models are inside a 'models' folder
		modelPath = "./models/"
	}

	return &Config{
		Port:      port,
		ModelPath: modelPath,
	}
}
