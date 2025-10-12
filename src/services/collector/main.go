package main

import (
	"log"
	"os"
)

type App struct {
	Config
}

func NewApp(cfg Config) *App {
	return &App{Config: cfg}
}

func main() {
	logger := log.New(os.Stdout, "[LogCollector] ", log.LstdFlags)

	//Load configuration
	config, err := loadConfig("/app/config.yml", logger)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return
	}

	app := NewApp(config)

	service, err := NewLogCollectorService(config, logger)
	if err != nil {
		log.Fatalf("failed to create collector service: %v", err)
	}

	go func() {
		if err := service.Run(); err != nil {
			logger.Fatalf("collector service error: %v", err)
		}
	}()

	server := app.setupRouter()
	logger.Println("Starting Gin server on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("failed to start Gin server: %v", err)
	}
}
