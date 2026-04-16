package main

import (
	"log"

	"github.com/savanyv/cartify/config"
	"github.com/savanyv/cartify/internal/app"
)

func main() {
	cfg := config.LoadConfig()

	server := app.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}