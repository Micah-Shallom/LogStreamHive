package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *App) setupRouter() *gin.Engine {
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

	router.GET("/logs", app.logsHandler)
	router.GET("/config", app.configHandler)
	router.GET("/statistics", app.getStatistics)

	return router
}

func (app *App) configHandler(c *gin.Context) {
	rd := BuildSuccessResponse(http.StatusOK, "Configuration retrieved successfully", app.Config)
	c.JSON(http.StatusOK, rd)
}

func (app *App) getStatistics(c *gin.Context) {
	stats := analyzeLogFiles(app)
	rd := BuildSuccessResponse(http.StatusOK, "Statistics retrieved successfully", stats)
	c.JSON(http.StatusOK, rd)
}

func (app *App) logsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 10
	}

	if _, err := os.Stat(app.Config.OutputFile); os.IsNotExist(err) {
		rd := BuildErrorResponse(http.StatusNotFound, "error", "Log file not found", err, []LogEntry{})
		c.JSON(http.StatusNotFound, rd)
		return
	}

	lines, err := tailLogs(app.Config.OutputFile, limit)
	if err != nil {
		rd := BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to read logs", err, []LogEntry{})
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	var logs []LogEntry
	for _, line := range lines {
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		logs = append(logs, entry)
	}

	rd := BuildSuccessResponse(http.StatusOK, "Logs retrieved successfully", logs)
	c.JSON(http.StatusOK, rd)
}
