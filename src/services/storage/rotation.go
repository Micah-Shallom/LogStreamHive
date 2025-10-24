package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (ls *LogStorage) checkRotation() {
	entries, err := os.ReadDir(ls.activeDir)
	if err != nil || len(entries) == 0 {
		return
	}

	// Check size-based rotation
	var totalSize int64
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				fmt.Println(err)
				continue
			}
			totalSize += info.Size()
		}
	}

	totalSizeMB := totalSize / (1024 * 1024)
	if totalSizeMB >= int64(ls.rotationSizeMB) {
		ls.rotateLogs("size")
		return
	}

	// Check time-based rotation
	if len(entries) > 0 {
		oldestEntry, _ := entries[0].Info()
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if info.ModTime().Before(oldestEntry.ModTime()) {
				oldestEntry = info
			}
		}

		if time.Since(oldestEntry.ModTime()) > time.Duration(ls.rotationHours)*time.Hour {
			ls.rotateLogs("time")
		}
	}
}

func (ls *LogStorage) rotateLogs(reason string) {
	fmt.Printf("Rotating logs due to %s trigger\n", reason)

	timestamp := time.Now().Unix()
	archiveSubDir := filepath.Join(ls.archiveDir, fmt.Sprintf("rotated_%d", timestamp))

	if err := os.MkdirAll(archiveSubDir, 0755); err != nil {
		fmt.Printf("Error creating archive directory: %v\n", err)
		return
	}

	entries, err := os.ReadDir(ls.activeDir)
	if err != nil {
		fmt.Printf("Error reading active directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			src := filepath.Join(ls.activeDir, entry.Name())
			dst := filepath.Join(archiveSubDir, entry.Name())
			if err := os.Rename(src, dst); err != nil {
				fmt.Printf("Error moving file %s: %v\n", entry.Name(), err)
			}
		}
	}

	fmt.Printf("Rotated logs to %s\n", archiveSubDir)
}
