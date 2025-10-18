package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type LogFileHandler struct {
	parser         *LogParser
	processedLines map[uint64]bool
	maxTrackedSize int
	outputDir      string
}

func NewLogFileHandler(outputDir string) *LogFileHandler {
	return &LogFileHandler{
		parser:         NewLogParser(),
		processedLines: make(map[uint64]bool),
		maxTrackedSize: 100000,
		outputDir:      outputDir,
	}
}

func (h *LogFileHandler) ProcessFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineHash := h.hashLine(line)

		if h.processedLines[lineHash] {
			continue
		}

		parsedLog := h.parser.Parse(strings.TrimSpace(line))

		if err := h.outputParsedLog(parsedLog, filePath); err != nil {
			log.Printf("Error outputting parsed log: %v", err)
		}

		h.processedLines[lineHash] = true

		if len(h.processedLines) > h.maxTrackedSize {
			h.cleanupProcessedLines()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return nil
}

func (h *LogFileHandler) cleanupProcessedLines() {
	//keep only the most recent 50% of entries
	//in production systems, a mosre sophisticated approach can be used like an LRU cache or time-based expiration

	newMap := make(map[uint64]bool)
	count := 0
	targetSize := h.maxTrackedSize / 2

	for hash := range h.processedLines {
		if count >= targetSize {
			break
		}
		newMap[hash] = true
		count++
	}

	h.processedLines = newMap
	log.Printf("cleaned up processed lines map. New size: %d", len(h.processedLines))
}

func (h *LogFileHandler) hashLine(line string) uint64 {
	hasher := fnv.New64a()
	hasher.Write([]byte(line))
	return hasher.Sum64()
}

func (h *LogFileHandler) outputParsedLog(parsedLog ParsedLog, sourceFile string) error {
	parsedLog.SourceFile = filepath.Base(sourceFile)

	outputDir := h.outputDir
	if outputDir == "" {
		outputDir = filepath.Join(filepath.Dir(sourceFile), "parsed")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %ww", err)
	}

	//create output file path
	outputFile := filepath.Join(outputDir, fmt.Sprintf("parsed_%s.json", filepath.Base(sourceFile)))

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
	fmt.Printf("parsed log from %s: %s\n", sourceFile, jsonStr)

	return nil
}

// watchdirectory monitors a directory for log file changes
func (h *LogFileHandler) WatchDirectory(logDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(logDir); err != nil {
		return fmt.Errorf("error adding directory to watcher: %w", err)
	}

	log.Printf("starting log processing service, monitoring directory: %s", logDir)

	if err := h.processExistingFiles(logDir); err != nil {
		log.Printf("error processing existing files: %v", err)
	}

	//watch for file changes
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//only process write events for .log files
				if event.Op&fsnotify.Write == fsnotify.Write {
					if strings.HasSuffix(event.Name, ".log") {
						if err := h.ProcessFile(event.Name); err != nil {
							log.Printf("error processing file %s: %v", event.Name, err)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("watcher error: %v", err)
			}
		}
	}()

	<-done
	return nil
}

func (h *LogFileHandler) processExistingFiles(logDir string) error {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".log") {
			filePath := filepath.Join(logDir, entry.Name())
			log.Printf("Processing existing file: %s", filePath)
			if err := h.ProcessFile(filePath); err != nil {
				log.Printf("Error processing file %s: %v", filePath, err)
			}
		}
	}

	return nil
}
