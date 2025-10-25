package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (ls *LogStorage) decompressFile(compressedPath string) ([]byte, error) {
	file, err := os.Open(compressedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open compressed file: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	data, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress file: %w", err)
	}

	return data, nil
}

func (ls *LogStorage) processCompressedFile(filePath string) error {
	data, err := ls.decompressFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to decompress file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var parsedData []LogEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Printf("Warning: failed to parse line: %v\n", err)
			continue
		}
		parsedData = append(parsedData, entry)
	}

	if len(parsedData) == 0 {
		return nil
	}

	timestamp := time.Now().Unix()
	storageFile := filepath.Join(ls.activeDir, fmt.Sprintf("logs_%d.json", timestamp))

	outData, err := json.MarshalIndent(parsedData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal logs: %w", err)
	}

	if err := os.WriteFile(storageFile, outData, 0644); err != nil {
		return fmt.Errorf("failed to write storage file: %w", err)
	}

	relPath, _ := filepath.Rel(ls.storageDir, storageFile)
	if err := ls.updateIndex(parsedData, relPath); err != nil {
		fmt.Printf("Error updating index: %v\n", err)
	}

	fmt.Printf("Stored %d log entries from %s\n", len(parsedData), filepath.Base(filePath))
	return nil
}
