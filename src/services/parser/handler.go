package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type LogFileHandler struct {
	parser    *LogParser
	outputDir string
	subject   string
}

func NewLogFileHandler(outputDir, subject string) *LogFileHandler {
	return &LogFileHandler{
		parser:    NewLogParser(),
		outputDir: outputDir,
		subject:   subject,
	}
}

func (h *LogFileHandler) ProcessLog(logLine string) error {

	parsedLog := h.parser.Parse(strings.TrimSpace(logLine))

	if err := h.publishLogToNatSubject(parsedLog); err != nil {
		log.Printf("Error outputing parsed log to nat subject: %v", err)
	}

	return nil
}

func (h *LogFileHandler) publishLogToNatSubject(parsedLog ParsedLog) error {
	jsonData, err := json.Marshal(parsedLog)
	if err != nil {
		return fmt.Errorf("error marshalling log to json: %w", err)
	}

	if h.NatsClient.NatsConn != nil {
		if err := h.NatsClient.PublishLogToNats(h.Subject, logMsg); err != nil {
			logger.Printf("Failed to publish log to NATS: %v", err)
		}
	}

	return nil
}

