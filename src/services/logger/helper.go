package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
)

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}

	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}

	log.Printf("Invalid int for %s: %s, defaulting to %d", key, valStr, defaultVal)
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valStr := strings.ToLower(getEnv(key, ""))
	if valStr == "" {
		return defaultVal
	}
	if valStr == "true" || valStr == "1" || valStr == "yes" {
		return true
	}
	if valStr == "false" || valStr == "0" || valStr == "no" {
		return false
	}
	log.Printf("Invalid bool for %s: %s, defaulting to %v", key, valStr, defaultVal)
	return defaultVal
}

func getEnvAsSlice(key string, defaultVal []string, sep string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	parts := strings.Split(valStr, sep)
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

func getEnvAsFloat(key string, defaultVal float64) float64 {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	log.Printf("Invalid float for %s: %s, defaulting to %f", key, valStr, defaultVal)
	return defaultVal
}

func tailLogs(filePath string, numLines int) ([]string, error) {
	if numLines <= 0 {
		return nil, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var (
		size    = stat.Size()
		chunk   = int64(1024) // read 1KB chunks from the end
		buffer  []byte
	)

	for pos := size; pos > 0 && len(lines) < numLines; pos -= chunk {
		if pos < chunk {
			chunk = pos
		}
		buf := make([]byte, chunk)
		_, err := file.ReadAt(buf, pos-chunk)
		if err != nil {
			return nil, err
		}

		// prepend buffer
		buffer = append(buf, buffer...)

		// split by lines
		for {
			idx := bytes.LastIndexByte(buffer, '\n')
			if idx == -1 {
				break
			}
			line := append([]byte(nil), buffer[idx+1:]...)
			buffer = buffer[:idx]
			if len(line) > 0 {
				lines = append([]string{string(line)}, lines...)
				if len(lines) == numLines {
					return lines, nil
				}
			}
		}
	}

	if len(buffer) > 0 && len(lines) < numLines {
		lines = append([]string{string(buffer)}, lines...)
	}

	return lines, nil
}
