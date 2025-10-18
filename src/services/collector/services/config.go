package services

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogFiles      []string         `yaml:"log_files" json:"log_files"`
	CheckInterval float64          `yaml:"check_interval" json:"check_interval"`
	Centrifugo    CentrifugoConfig `yaml:"centrifugo" json:"centrifugo"`
	ChannelID     string           `yaml:"channel_id" json:"channel_id"`
	Nats          NatsConfig       `yaml:"nats"`
}

func NewLogFileHandler(filePath string, logger *log.Logger, centrifugoClient *CentrifugoClient, natsClient *NatsClient, channelID string) (*LogFileHandler, error) {
	handler := &LogFileHandler{
		FilePath:         filePath,
		LastPosition:     0,
		Logger:           logger,
		ChannelID:        channelID,
		CentrifugoClient: centrifugoClient,
		NatsClient:       natsClient,
	}

	if err := handler.InitializePosition(); err != nil {
		return nil, fmt.Errorf("failed to initialize position: %w", err)
	}

	return handler, nil
}

func (h *LogFileHandler) InitializePosition() error {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	file, err := os.Open(h.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}
	defer file.Close()

	pos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	h.LastPosition = pos
	return nil
}

func (h *LogFileHandler) CollectNewLogs(logger *log.Logger) error {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	file, err := os.Open(h.FilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	_, err = file.Seek(h.LastPosition, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	hasNewContent := false

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			fmt.Printf("Collected: %s\n", line)
			if h.CentrifugoClient != nil {
				logMsg := LogMessage{
					Timestamp: time.Now().Format(time.RFC3339),
					FilePath:  h.FilePath,
					Line:      line,
				}
				if err := h.CentrifugoClient.PublishLog(h.ChannelID, logMsg); err != nil {
					logger.Printf("Failed to publish log: %v", err)
				}

			}

			if h.NatsClient.NatsConn != nil {
				logMsg := LogMessage{
					Timestamp: time.Now().Format(time.RFC3339),
					FilePath:  h.FilePath,
					Line:      line,
				}
				if err := h.NatsClient.PublishLogToNats("log.raw", logMsg); err != nil {
					logger.Printf("Failed to publish log to NATS: %v", err)
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

	h.LastPosition = newPos
	return nil
}

func LoadConfig(configPath string, logger *log.Logger) (Config, error) {
	var config Config

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			logger.Printf("Error reading config file: %v", err)
			return config, nil
		}

		if err := yaml.Unmarshal(data, &config); err != nil {
			logger.Printf("Error parsing config file: %v", err)
			return config, nil
		}

		logger.Println("Loaded configuration from file")
	}

	logger.Printf("Monitoring log files: %v", config.LogFiles)
	logger.Printf("Publishing to channel: %s", config.ChannelID)
	return config, nil
}
