package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type RetentionPolicy interface {
	Apply(logDirectory, logFilePattern string) error
}

// keeps a maximum number of log files
type CountBasedRetentionPolicy struct {
	MaxFiles int
}

func NewCountBasedRetentionPolicy(maxFiles int) *CountBasedRetentionPolicy {
	return &CountBasedRetentionPolicy{MaxFiles: maxFiles}
}

// deletes oldest files when count exceeds threshold
func (p *CountBasedRetentionPolicy) Apply(logDirectory, logFilePattern string) error {
	pattern := filepath.Join(logDirectory, logFilePattern)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	//sort files by creation time
	sort.Slice(matches, func(i, j int) bool {
		infoI, errI := os.Stat(matches[i])
		infoJ, errJ := os.Stat(matches[j])
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	//remove excess files
	if len(matches) > p.MaxFiles {
		filesToRemove := matches[:len(matches)-p.MaxFiles]
		for _, filePath := range filesToRemove {
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("error deleting %s: %v\n", filePath, err)
			} else {
				fmt.Printf("deleted old log files: %s\n", filePath)
			}
		}
	}
	return nil
}

// deletes log files older than a specified date
type AgeBasedRetentionPolicy struct {
	MaxAgeSeconds int64
}

func NewAgeBasedRetentionPolicy(maxAgeDays int) *AgeBasedRetentionPolicy {
	return &AgeBasedRetentionPolicy{
		MaxAgeSeconds: int64(maxAgeDays) * 24 * 60 * 60,
	}
}

func (p *AgeBasedRetentionPolicy) Apply(logDirectory, logFilePattern string) error {
	pattern := filepath.Join(logDirectory, logFilePattern)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	for _, logFile := range matches {
		info, err := os.Stat(logFile)
		if err != nil {
			continue
		}

		fileAge := currentTime.Unix() - info.ModTime().Unix()
		if fileAge > p.MaxAgeSeconds {
			if err := os.Remove(logFile); err != nil {
				fmt.Printf("error deleting %s: %v\n", logFile, err)
			} else {
				fmt.Printf("deleted expired log file: %s\n", logFile)
			}
		}
	}
	
	return nil
}
