package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type LogEntry struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Source    string `json:"source"`
}

type AppConfig struct {
	Config
}

var logFilePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	// Navigate up three levels: logger -> services -> src -> root
	rootDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "..")
	logFilePath = filepath.Join(rootDir, "log-output", "service.log")

	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}
}

func main() {
	cfg := LoadConfig()

	generator := NewLogGenerator(cfg)
	go generator.Run(0)

	router := setupRouter()
	log.Println("Starting Gin server on :8000")
	if err := router.Run(":8000"); err != nil {
		log.Fatalf("failed to start Gin server: %v", err)
	}
}
