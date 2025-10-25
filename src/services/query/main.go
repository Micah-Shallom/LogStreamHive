package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	storageDir := flag.String("storage-dir", "/data/storage", "Storage directory for log data")
	indexType := flag.String("index-type", "", "Type of index to search (date, level, status, service, etc.)")
	indexValue := flag.String("index-value", "", "Value to search for in the index")
	pattern := flag.String("pattern", "", "Regular expression pattern to search for in logs")
	format := flag.String("format", "text", "Output format (text or json)")

	flag.Parse()

	query := NewLogQuery(*storageDir)

	var results []map[string]any
	var err error

	if *indexType != "" && *indexValue != "" {
		// validIndexTypes := map[string]bool{"date": true, "level": true, "status": true, "service": true, "format": true, "ip": true, "method": true, "user": true}
		index_path := fmt.Sprintf("%s/index", *storageDir)
		entries, err := os.ReadDir(index_path)
		if err != nil {
			fmt.Println("Error reading storage logs directory:", err)
			fmt.Println("No indexes available.")
			os.Exit(1)
		}
		validIndexTypes := make(map[string]bool)
		for _, entry := range entries {
			name := entry.Name()
			validIndexTypes[name] = true
		}
		if !validIndexTypes[*indexType] {
			fmt.Printf("Invalid index type. Must be one of: %v\n", func() []string {
				keys := make([]string, 0, len(validIndexTypes))
				for k := range validIndexTypes {
					keys = append(keys, k)
				}
				return keys
			}())
			os.Exit(1)
		}

		fmt.Printf("Searching for logs with %s=%s\n", *indexType, *indexValue)
		results, err = query.FindLogsByIndex(*indexType, *indexValue)
		if err != nil {
			fmt.Printf("Error searching index: %v\n", err)
			os.Exit(1)
		}
	} else if *pattern != "" {
		fmt.Printf("Searching for logs matching pattern: %s\n", *pattern)
		results, err = query.SearchAllLogs(*pattern)
		if err != nil {
			fmt.Printf("Error searching logs: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Please specify either --index-type and --index-value, or --pattern")
		flag.Usage()
		os.Exit(1)
	}

	query.DisplayResults(results, *format)
}
