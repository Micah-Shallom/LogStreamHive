package nats

import (
	"collector/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (c *models.CentrifugoClient) PublishLogToNats(channelID string, logMsg models.LogMessage) error {
	if channelID == "" {
		return fmt.Errorf("empty channel_id supplied")
	}

	payload, err := json.Marshal(logMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal log message: %w", err)
	}

	err = c.Client.Publish(context.Background(), channelID, payload)
	if err != nil {
		c.Logger.Printf("Failed to publish to channel %s: %v", channelID, err)
		return err
	}

	return nil
}
