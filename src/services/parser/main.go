package parser

import (
	"log"
	"os"
	"path/filepath"
)



func main() {
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "../logs"
	}

	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		outputDir = filepath.Join(logDir, "parsed")
	}

	handler := NewLogFileHandler(outputDir)

	if err := handler.WatchDirectory(logDir); err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}