package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v3"
)


func NewLogFileHandler(filePath string, logger *log.Logger, centrifugoClient *CentrifugoClient, channelID string) (*LogFileHandler, error) {
	handler := &LogFileHandler{
		filePath:         filePath,
		lastPosition:     0,
		logger:           logger,
		channelID:        channelID,
		centrifugoClient: centrifugoClient,
		natsConn:         &nats.Conn{},
	}

	// Initialize position at end of file
	if err := handler.initializePosition(); err != nil {
		return nil, fmt.Errorf("failed to initialize position: %w", err)
	}

	return handler, nil
}

func (h *LogFileHandler) initializePosition() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	file, err := os.Open(h.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			h.logger.Printf("file doesnt exist yet: %s", h.filePath)
			return nil
		}
	}
	defer file.Close()

	pos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	h.lastPosition = pos
	h.logger.Printf("Initialized tracking for %s at position %d", h.filePath, pos)
	return nil
}

func (h *LogFileHandler) collectNewLogs() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	file, err := os.Open(h.filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	_, err = file.Seek(h.lastPosition, io.SeekStart) ///seek to last read position
	if err != nil {
		return fmt.Errorf("error seeking file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	hasNewContent := false

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			fmt.Printf("Collected: %s\n", line)
			if h.centrifugoClient != nil {
				logMsg := LogMessage{
					Timestamp: time.Now().Format(time.RFC3339),
					FilePath:  h.filePath,
					Line:      line,
				}
				if err := h.centrifugoClient.PublishLog(h.channelID, logMsg); err != nil {
					h.logger.Printf("Failed to publish log to Centrifugo channel '%s': %v", h.channelID, err)
				}

				if h.natsConn != nil {
					if err := h.natsConn.Publish("logs.raw", []byte(logMsg)); err != nil {
						h.logger.Printf("Failed to publish log to NATS: %v", err)
					}
				}
			}

			hasNewContent = true
		}
	}

	if hasNewContent {
		fmt.Println("---------------------------------------")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	newPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("error getting file position: %w", err)
	}

	h.lastPosition = newPos
	return nil
}

func loadConfig(configPath string, logger *log.Logger) (Config, error) {
	// Default configuration
	// config := Config{
	// 	LogFiles:      []string{"./logs/service.log"},
	// 	CheckInterval: 0.5,
	// 	Centrifugo: CentrifugoConfig{
	// 		APIKey: "",
	// 		URL:    "http://centrifugo:8080",
	// 	},
	// 	ChannelID: "logs",
	// }

	var config Config

	// Try to load from file
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			logger.Printf("Error reading config file: %v", err)
			logger.Println("Using default configuration")
			return config, nil
		}

		if err := yaml.Unmarshal(data, &config); err != nil {
			logger.Printf("Error parsing config file: %v", err)
			logger.Println("Using default configuration")
			return config, nil
		}

		logger.Println("Loaded configuration from file")
	} else {
		logger.Println("Config file not found, using default configuration")
	}

	logger.Printf("Monitoring log files: %v", config.LogFiles)
	logger.Printf("Publishing to channel: %s", config.ChannelID)
	return config, nil
}
