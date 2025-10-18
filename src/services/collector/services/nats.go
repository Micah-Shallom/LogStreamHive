package services

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	NatsConn *nats.Conn
}

type NatsConfig struct {
	URL string `yaml:"url"`
}

func (c *NatsClient) PublishLogToNats(channel string, logMsg LogMessage) error {
	if channel == "" {
		return fmt.Errorf("empty channel supplied")
	}

	payload, err := json.Marshal(logMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal log message: %w", err)
	}

	err = c.NatsConn.Publish(channel, payload)
	if err != nil {
		return fmt.Errorf("failed to publish to channel %s: %w", channel, err)
	}

	return nil
}
