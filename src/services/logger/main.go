package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
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

const logFilePath = "/var/log/logger/service.log"

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

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/config", configHandler)
	router.GET("/logs", logsHandler)

	return router
}

func configHandler(c *gin.Context) {
	cfg := AppConfig{
		Config: LoadConfig(),
	}
	c.JSON(http.StatusOK, cfg)
}

func logsHandler(c *gin.Context) {
	file, err := os.ReadFile(logFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read log file"})
		return
	}

	logStrings := strings.Split(strings.TrimSpace(string(file)), "\n")
	var logs []LogEntry

	for _, logStr := range logStrings {
		var logEntry LogEntry
		if err := json.Unmarshal([]byte(logStr), &logEntry); err == nil {
			logs = append(logs, logEntry)
		}
	}

	c.JSON(http.StatusOK, logs)
}
