package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	if err := h.outputParsedLogLine(parsedLog); err != nil {
		log.Printf("Error outputting parsed log: %v", err)
	}

	return nil
}

func (h *LogFileHandler) outputParsedLogLine(parsedLog ParsedLog) error {


	if err := os.MkdirAll(h.outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	fileName := fmt.Sprintf("parsed_logs_%s.json", time.Now().Format("2006-01-02"))
	outputFile := filepath.Join(h.outputDir, fileName)

	//open file in append mode
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening output file: %w", err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(parsedLog)
	if err != nil {
		return fmt.Errorf("error marshalling log to json: %w", err)
	}

	if _, err := file.Write(append(jsonData, '\n')); err != nil {
		return fmt.Errorf("error writing to output file: %w", err)
	}

	jsonStr := string(jsonData)
	if len(jsonStr) > 100 {
		jsonStr = jsonStr[:100] + "..."
	}
	fmt.Printf("parsed log: %s\n", jsonStr)

	return nil
}
