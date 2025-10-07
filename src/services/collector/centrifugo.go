package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/centrifugal/gocent"
)

type CentClient struct {
	C *gocent.Client
}

var Client *gocent.Client = &gocent.Client{}

func Connection() *gocent.Client{
	return Client
}

type CentrifugoConfig struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
}

type CentrifugoClient struct {
	client *gocent.Client
	logger *log.Logger
}

func NewCentrifugoClient(config CentrifugoConfig, logger *log.Logger) (*CentrifugoClient, error) {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	client := gocent.New(gocent.Config{
		Addr:       config.URL,
		Key:        config.APIKey,
		HTTPClient: httpClient,
	})

	logger.Printf("Connected to Centrifugo server at %s", config.URL)

	return &CentrifugoClient{
		client: client,
		logger: logger,
	}, nil
}

type LogMessage struct {
	Timestamp string `json:"timestamp"`
	FilePath  string `json:"file_path"`
	Line      string `json:"line"`
}

func (c *CentrifugoClient) PublishLog(channelID string, logMsg LogMessage) error {
	if channelID == "" {
		return fmt.Errorf("empty channel_id supplied")
	}

	payload, err := json.Marshal(logMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal log message: %w", err)
	}

	err = c.client.Publish(context.Background(), channelID, payload)
	if err != nil {
		c.logger.Printf("Failed to publish to channel %s: %v", channelID, err)
		return err
	}

	return nil
}

func (c *CentrifugoClient) PublishBatch(channelID string, logMessages []LogMessage) error {
	if channelID == "" {
		return fmt.Errorf("empty channel_id supplied")
	}

	for _, msg := range logMessages {
		if err := c.PublishLog(channelID, msg); err != nil {
			return err
		}
	}

	return nil
}