package utils

import (
	"encoding/json"
	"log"
	"os"
	"travel-agent-backend/models"
)

//LoadConfig loads config.json file from the given path
func LoadConfig(path string) models.Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error in opening config file: %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := models.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error in decoding config: %v", err)
	}
	return config
}

//WriteToFile writes data to a given filepath
func WriteToFile(filePath string, data []byte) error {
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}
