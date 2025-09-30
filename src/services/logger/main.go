package main

import (
	"log"
)

type LogEntry struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Source    string `json:"source"`
}

type App struct {
	Config
}

func NewApp(cfg Config) *App{
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
