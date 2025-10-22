package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogFileHandler struct {
	parser     *LogParser
	outputDir  string
	outputFile string
	subject    string
	logger     *log.Logger
}

func NewLogFileHandler(outputDir, outputFile, subject string) *LogFileHandler {
	var writers []io.Writer

	if outputFile != "" {
		dir := filepath.Dir(outputFile)

		if err := os.MkdirAll(dir, 0777); err != nil {
			log.Fatalf("error creating output directory: %v", err)
		}

		if err := os.Chmod(dir, 0777); err != nil {
			log.Fatalf("error setting permissions on output directory: %v", err)
		}

		if _, err := os.Stat(outputFile); err == nil {
			_ = os.Chmod(outputFile, 0777)
		}

		logfile := &lumberjack.Logger{
			Filename:   outputFile,
			MaxSize:    1, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}

		writers = append(writers, logfile)
	}

	multi := io.MultiWriter(writers...)
	return &LogFileHandler{
		parser:     NewLogParser(),
		outputDir:  outputDir,
		outputFile: outputFile,
		subject:    subject,
		logger:     log.New(multi, "", 0),
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
	logEntryBytes, err := json.Marshal(parsedLog)
	if err != nil {
		return fmt.Errorf("error marshaling parsed log to JSON: %v", err)
	}

	h.logger.Println(string(logEntryBytes))
	return nil
}
