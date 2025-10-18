package main

import (
	"collector/services"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "[LogCollector] ", log.LstdFlags)

	//Load configuration
	config, err := services.LoadConfig("/app/config.yml", logger)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return
	}

	app := services.NewApp(config)

	service, err := services.NewLogCollectorService(config, logger)
	if err != nil {
		log.Fatalf("failed to create collector service: %v", err)
	}

	go func() {
		if err := service.Run(); err != nil {
			logger.Fatalf("collector service error: %v", err)
		}
	}()

	server := app.SetupRouter()
	logger.Println("Starting Gin server on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("failed to start Gin server: %v", err)
	}
}
