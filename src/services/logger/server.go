package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

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
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, []LogEntry{})
		return
	}

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