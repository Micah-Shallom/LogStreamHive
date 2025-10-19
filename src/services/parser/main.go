package main

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		outputDir = "logs/parsed"
	}
	subject := "logs.raw"

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("failed to connect to nats: %v", err)
	}
	defer nc.Close()

	handler := NewLogFileHandler(outputDir, subject)

	_, err = nc.Subscribe(subject, func(m *nats.Msg) {
		log.Printf("Received log: %s", string(m.Data))
		if err := handler.ProcessLog(string(m.Data)); err != nil {
			log.Printf("failed to process log: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("failed to subscribe to subject '%s': %v", subject, err)
	}

	log.Printf("Subscribed to '%s' subject", subject)

	// Keep the process running
	select {}
}
