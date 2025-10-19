package main

//custom log management functionlaity....replaced with lumberjack for efficiency

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type LogStorage struct {
	LogDirectory    string
	BaseFilename    string
	RotationPolicy  RotationPolicy
	RetentionPolicy RetentionPolicy
	CompressRotated bool
	CurrentLogPath  string
	Mu              sync.Mutex
}

func NewLogStorage(
	logDirectory string,
	baseFilename string,
	rotationPolicy RotationPolicy,
	retentionPolicy RetentionPolicy,
	compressRotated bool,
) (*LogStorage, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	return &LogStorage{
		LogDirectory:    logDirectory,
		BaseFilename:    baseFilename,
		RotationPolicy:  rotationPolicy,
		RetentionPolicy: retentionPolicy,
		CompressRotated: compressRotated,
		CurrentLogPath:  filepath.Join(logDirectory, baseFilename),
	}, nil
}

func (ls *LogStorage) WriteLog(logMessage string) error {
	ls.Mu.Lock()
	defer ls.Mu.Unlock()

	if ls.RotationPolicy != nil && ls.RotationPolicy.ShouldRotate(ls.CurrentLogPath) {
		if err := ls.RotateLog(); err != nil {
			return fmt.Errorf("failed to rotate log: %w", err)
		}
	}

	if !strings.HasSuffix(logMessage, "\n") {
		logMessage += "\n"
	}

	file, err := os.OpenFile(ls.CurrentLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(logMessage); err != nil {
		return fmt.Errorf("failed to write log message: %w", err)
	}

	return nil
}

func (ls *LogStorage) RotateLog() error {
	if _, err := os.Stat(ls.CurrentLogPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	//generate timestamp for rotated filename
	timestamp := time.Now().Format("20060102-150405")
	ext := filepath.Ext(ls.BaseFilename)
	baseName := strings.TrimSuffix(ls.BaseFilename, ext)
	rotatedFilename := fmt.Sprintf("%s.%s.log", baseName, timestamp)
	rotatedPath := filepath.Join(ls.LogDirectory, rotatedFilename)

	//rename current log to rotated filename
	if err := os.Rename(ls.CurrentLogPath, rotatedPath); err != nil {
		return fmt.Errorf("failed to remove log file: %w", err)
	}

	if ls.CompressRotated {
		compressedPath := rotatedPath + ".gz"
		if err := ls.CompressFile(rotatedPath, compressedPath); err != nil {
			return fmt.Errorf("failed to compress log file: %w", err)
		}
		if err := os.Remove(rotatedPath); err != nil {
			return fmt.Errorf("failed to remove uncompressed log file: %w", err)
		}
		fmt.Printf("Rotated and compressed log to %s\n", compressedPath)
	} else {
		fmt.Printf("Rotated log to %s\n", rotatedPath)
	}

	if ls.RetentionPolicy != nil {
		pattern := fmt.Sprintf("%s.*", baseName)
		if err := ls.RetentionPolicy.Apply(ls.LogDirectory, pattern); err != nil {
			return fmt.Errorf("failed to apply retention policy: %w", err)
		}
	}

	return nil
}

// compressFile compresses a file using gzip
func (ls *LogStorage) CompressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	gzipWriter := gzip.NewWriter(dstFile)
	defer gzipWriter.Close()

	if _, err := io.Copy(gzipWriter, srcFile); err != nil {
		return fmt.Errorf("failed to compress file: %w", err)
	}

	return nil
}

func (ls *LogStorage) Close() error {
	return nil
}
