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
		if level, ok := entry["log_type"].(string); ok {
			indexKeys["level"] = level
		}

		//idnex by service
		if service, ok := entry["service"].(string); ok {
			indexKeys["service"] = service
		}
		
		//index by status
		if status, ok := entry["status"]; ok {
			indexKeys["status"] = fmt.Sprint(status)
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
				"line": i,
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
	return nil
}

func (ls *LogStorage) Run() {
	fmt.Printf("Starting log storage system\n")
	fmt.Printf("Watching for parsed logs in %s\n", ls.inputDir)
	fmt.Printf("Storing logs in %s\n", ls.storageDir)
	fmt.Printf("Rotation policy: %dMB or %d hours\n", ls.rotationSizeMB, ls.rotationHours)

	ticker := time.NewTicker(ls.interval)
	defer ticker.Stop()


	for range ticker.C {
		// Check for new parsed log files
		fmt.Printf("[DEBUG] Ticker triggered at %s\n", time.Now().Format(time.RFC3339))

		entries, err := os.ReadDir(ls.inputDir)
		if err != nil {
			fmt.Printf("Error reading input directory: %v\n", err)
			continue
		}

		fmt.Printf("[DEBUG] Found %d entries in input directory\n", len(entries))

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if strings.HasSuffix(entry.Name(), ".json") || strings.HasSuffix(entry.Name(), ".log") {
				fmt.Println("Skipping non-compressed file:", entry.Name())
				continue
			}

			filePath := filepath.Join(ls.inputDir, entry.Name())

			ls.processedMutex.RLock()
			processed := ls.processedFiles[filePath]
			ls.processedMutex.RUnlock()

			if processed {
				fmt.Println("Skipping already processed file:", filePath)
				continue
			}

			fmt.Printf("Processing parsed file: %s\n", filePath)
			if err := ls.processCompressedFile(filePath); err != nil {
				fmt.Printf("Error processing compressed file: %v\n", err)
				continue
			}

			ls.processedMutex.Lock()
			ls.processedFiles[filePath] = true
			ls.processedMutex.Unlock()

			if err := ls.saveTrackingData(); err != nil {
				fmt.Printf("Error saving tracking data: %v\n", err)
			}
		}

		// Check if rotation is needed
		ls.checkRotation()
	}
}
