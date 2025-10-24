package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	NatsConn *nats.Conn
}

type NatsConfig struct {
	URL string `yaml:"url"`
}

func (c *NatsClient) PublishLogToNats(subject string, logMsg string) error {
	if subject == "" {
		return fmt.Errorf("empty subject supplied")
	}

	payload, err := json.Marshal(logMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal log message: %w", err)
	}

	err = c.NatsConn.Publish(subject, payload)
	if err != nil {
		return fmt.Errorf("failed to publish to subject %s: %w", subject, err)
	}

	log.Printf("Published log to NATS subject %s", subject)

	return nil
}
