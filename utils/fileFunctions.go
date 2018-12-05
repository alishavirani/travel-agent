package utils

import (
	"encoding/json"
	"log"
	"os"
	"travel-agent-backend/models"
)

func LoadConfig(path string) models.Config {
	file, err := os.Open(path)
	if err != nil {
		log.Panic("Error in opening config file:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := models.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Panic("Error in decoding config:", err)
	}
	return config
}
