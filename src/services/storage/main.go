package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	inputDir := flag.String("input-dir", "/data/parsed", "Input directory with parsed logs")
	storageDir := flag.String("storage-dir", "/data/storage", "Storage directory for log data")
	rotationSize := flag.Int("rotation-size", 10, "Size-based rotation trigger in MB")
	rotationHours := flag.Int("rotation-hours", 24, "Time-based rotation trigger in hours")
	intervalSecs := flag.Float64("interval", 10.0, "Polling interval in seconds")

	flag.Parse()

	interval := time.Duration(float64(time.Second) * (*intervalSecs))

	storage, err := NewLogStorage(
		*inputDir,
		*storageDir,
		*rotationSize,
		*rotationHours,
		interval,)
	if err != nil {
		fmt.Printf("Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	storage.Run()
}
