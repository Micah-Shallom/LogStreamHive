package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type IndexEntry struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

type LogQuery struct {
	storageDir string
	indexDir   string
	activeDir  string
	archiveDir string
}

func NewLogQuery(storageDir string) *LogQuery {
	return &LogQuery{
		storageDir: storageDir,
		indexDir:   filepath.Join(storageDir, "index"),
		activeDir:  filepath.Join(storageDir, "active"),
		archiveDir: filepath.Join(storageDir, "archive"),
	}
}

func (lq *LogQuery) FindLogsByIndex(indexType, indexValue string) ([]map[string]any, error) {
	indexFile := filepath.Join(lq.indexDir, indexType, indexValue+".idx")
	fmt.Println(indexFile)

	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		fmt.Printf("No index found for %s=%s\n", indexType, indexValue)
		return []map[string]any{}, nil
	}

	data, err := os.ReadFile(indexFile)
	if err != nil {
		return nil, fmt.Errorf("error reading index file: %w", err)
	}

	results := []map[string]any{}
	lines := strings.SplitSeq(string(data), "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var indexEntry IndexEntry
		if err := json.Unmarshal([]byte(line), &indexEntry); err != nil {
			fmt.Printf("Error parsing index entry: %v\n", err)
			continue
		}

		logFile := filepath.Join(lq.storageDir, indexEntry.File)
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			continue
		}

		logData, err := os.ReadFile(logFile)
		if err != nil {
			fmt.Printf("Error reading log file %s: %v\n", logFile, err)
			continue
		}

		var logs []map[string]any
		if err := json.Unmarshal(logData, &logs); err != nil {
			fmt.Printf("Error parsing log file %s: %v\n", logFile, err)
			continue
		}

		// If line number is specified, try to match it
		if indexEntry.Line >= 0 && indexEntry.Line < len(logs) {
			results = append(results, logs[indexEntry.Line])
		} else {
			// Otherwise add all logs from this file
			results = append(results, logs...)
		}
	}

	return results, nil
}

func (lq *LogQuery) SearchAllLogs(pattern string) ([]map[string]any, error) {
	regex, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	results := []map[string]any{}

	// Search in active logs
	activeFiles, err := os.ReadDir(lq.activeDir)
	if err == nil {
		for _, entry := range activeFiles {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				filePath := filepath.Join(lq.activeDir, entry.Name())
				fileResults := lq.searchFile(filePath, regex)
				results = append(results, fileResults...)
			}
		}
	}

	// Search in archived logs
	archiveDirs, err := os.ReadDir(lq.archiveDir)
	if err == nil {
		for _, archiveDir := range archiveDirs {
			if archiveDir.IsDir() {
				archivePath := filepath.Join(lq.archiveDir, archiveDir.Name())
				archiveFiles, err := os.ReadDir(archivePath)
				if err != nil {
					continue
				}

				for _, entry := range archiveFiles {
					if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
						filePath := filepath.Join(archivePath, entry.Name())
						fileResults := lq.searchFile(filePath, regex)
						results = append(results, fileResults...)
					}
				}
			}
		}
	}

	return results, nil
}

func (lq *LogQuery) searchFile(filePath string, regex *regexp.Regexp) []map[string]any {
	results := []map[string]any{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return results
	}

	var logs []map[string]any
	if err := json.Unmarshal(data, &logs); err != nil {
		fmt.Printf("Error parsing file %s: %v\n", filePath, err)
		return results
	}

	for _, entry := range logs {
		// Search in the raw field
		if raw, ok := entry["raw"].(string); ok && regex.MatchString(raw) {
			results = append(results, entry)
			continue
		}

		// Search in the message field
		if message, ok := entry["message"].(string); ok && regex.MatchString(message) {
			results = append(results, entry)
		}
	}

	return results
}

func (lq *LogQuery) DisplayResults(results []map[string]any, format string) {
	if len(results) == 0 {
		fmt.Println("No matching logs found")
		return
	}

	fmt.Printf("Found %d matching log entries:\n", len(results))

	if format == "json" {
		output, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting results: %v\n", err)
			return
		}
		fmt.Println(string(output))
	} else {
		// Text format
		for i, entry := range results {
			fmt.Printf("\n--- Result %d ---\n", i+1)

			if timestamp, ok := entry["timestamp"].(string); ok {
				fmt.Printf("Time: %s\n", timestamp)
			}

			if level, ok := entry["level"].(string); ok {
				fmt.Printf("Level: %s\n", level)
			}

			if logLevel, ok := entry["log_level"].(string); ok {
				fmt.Printf("Level: %s\n", logLevel)
			}

			if sourceFile, ok := entry["source_file"].(string); ok {
				fmt.Printf("Source: %s\n", sourceFile)
			}

			if service, ok := entry["service"].(string); ok {
				fmt.Printf("Service: %s\n", service)
			}

			if raw, ok := entry["raw"].(string); ok {
				fmt.Printf("Log: %s\n", raw)
			} else if message, ok := entry["message"].(string); ok {
				fmt.Printf("Message: %s\n", message)
			}

			fmt.Println(strings.Repeat("-", 40))
		}
	}
}
