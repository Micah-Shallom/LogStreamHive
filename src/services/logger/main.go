package main

import (
	"log"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	LogType   string `json:"log_type"`
	UserID    string `json:"user_id"`
	Duration  int    `json:"duration"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Service   string `json:"service"`
}

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
