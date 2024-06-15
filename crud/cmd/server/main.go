package main

import (
	"log"

	"github.com/drawiin/go-expert/crud/config"
)

func main() {
	// LoadConfig function is called to load the configuration values from the .env file
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.Default().Print("Config loaded successfully %v", config)
}
