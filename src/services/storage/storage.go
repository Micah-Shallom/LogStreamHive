package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type LogEntry map[string]any

type StorageMetadata struct {
	ProcessedFiles []string `json:"processed_files"`
	LastUpdated    string   `json:"last_updated"`
}

type LogStorage struct {
	inputDir       string
	storageDir     string
	rotationSizeMB int
	rotationHours  int
	interval       time.Duration
	processedFiles map[string]bool
	processedMutex sync.RWMutex
	indexDir       string
	activeDir      string
	archiveDir     string
	trackingFile   string
}

func NewLogStorage(inputDir, storageDir string, rotationSizeMB, rotationHours int, interval time.Duration) (*LogStorage, error) {
	ls := &LogStorage{
		inputDir:       inputDir,
		storageDir:     storageDir,
		rotationSizeMB: rotationSizeMB,
		rotationHours:  rotationHours,
		interval:       interval,
		processedFiles: make(map[string]bool),
	}

	ls.indexDir = filepath.Join(storageDir, "index")
	ls.activeDir = filepath.Join(storageDir, "active")
	ls.archiveDir = filepath.Join(storageDir, "archive")
	ls.trackingFile = filepath.Join(storageDir, "storage_tracking.json")

	// Create directories
	dirs := []string{ls.storageDir, ls.indexDir, ls.activeDir, ls.archiveDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Load tracking data
	if err := ls.loadTrackingData(); err != nil {
		fmt.Printf("Warning: failed to load tracking data: %v\n", err)
	}

	return ls, nil
}

func (ls *LogStorage) loadTrackingData() error {
	data, err := os.ReadFile(ls.trackingFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var metadata StorageMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return err
	}

	ls.processedMutex.Lock()
	defer ls.processedMutex.Unlock()

	for _, file := range metadata.ProcessedFiles {
		ls.processedFiles[file] = true
	}

	return nil
}

func (ls *LogStorage) saveTrackingData() error {
	ls.processedMutex.RLock()
	defer ls.processedMutex.RUnlock()

	files := make([]string, 0, len(ls.processedFiles))
	for file := range ls.processedFiles {
		files = append(files, file)
	}

	metadata := StorageMetadata{
		ProcessedFiles: files,
		LastUpdated:    time.Now().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(metadata, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(ls.trackingFile, data, 0644)
}

func (ls *LogStorage) updateIndex(parsedData []LogEntry, storageFile string) error {
	for i, entry := range parsedData {
		indexKeys := make(map[string]string)

		//index by timestamp (day)
		if ts, ok := entry["timestamp"].(string); ok {
			parts := strings.Split(ts, ":")
			if len(parts) > 0 {
				indexKeys["date"] = parts[0]
			}
		}

		//index by level
		if level, ok := entry["level"].(string); ok {
			indexKeys["level"] = level
		}

		//index by status
		if status, ok := entry["status"]; ok {
			indexKeys["status"] = fmt.Sprint(status)
		}

		//idnex by service
		if service, ok := entry["service"].(string); ok {
			indexKeys["service"] = service
		}

		//create index entries
		for keyType, keyValue := range indexKeys {
			indexTypeDir := filepath.Join(ls.indexDir, keyType)
			if err := os.MkdirAll(indexTypeDir, 0755); err != nil {
				fmt.Printf("Error creating index directory: %v\n", err)
				continue
			}

			indexFile := filepath.Join(indexTypeDir, keyValue+".idx")
			indexEntry := map[string]any{
				"file": storageFile,
				"line": 1,
			}

			data, err := json.Marshal(indexEntry)
			if err != nil {
				fmt.Printf("Error marshaling index entry: %v\n", err)
				continue
			}

			f, err := os.OpenFile(indexFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error opening index file: %v\n", err)
				continue
			}

			if _, err := f.Write(append(data, '\n')); err != nil {
				f.Close()
				fmt.Printf("Error writing to index file: %v\n", err)
				continue
			}

			f.Close()

		}
	}
}
