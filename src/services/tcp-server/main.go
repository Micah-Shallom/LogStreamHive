package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Server host address")
	port := flag.Int("port", 9000, "Server port")
	bufferSize := flag.Int("buffer-size", 4096, "Buffer size for reading")
	natsURL := flag.String("nats-url", "nats://localhost:4222", "NATS server URL")
	natsSubject := flag.String("nats-subject", "logs.raw", "NATS subject to publish to")

	flag.Parse()
	logger := log.New(os.Stdout, "[TCPServer] ", log.LstdFlags)

	config := ServerConfig{
		Host:        *host,
		Port:        *port,
		BufferSize:  *bufferSize,
		NatsURL:     *natsURL,
		NatsSubject: *natsSubject,
	}

	server, err := NewTCPServer(config, logger)
	if err != nil {
		logger.Fatalf("Failed to create server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Println("received shutdown signal")
		server.Stop()
	}()

	if err := server.Start(); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
