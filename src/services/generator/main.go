package main

import (
	"log"
)


type App struct {
	Config
}

func NewApp(cfg Config) *App {
	return &App{Config: cfg}
}

func main() {
	cfg := LoadConfig()

	app := NewApp(cfg)

	generator := NewLogGenerator(cfg)
	go generator.Run(0)

	router := app.setupRouter()
	log.Println("Starting Gin server on :8000")
	if err := router.Run(":8000"); err != nil {
		log.Fatalf("failed to start Gin server: %v", err)
	}
}
