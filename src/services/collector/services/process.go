package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/nats-io/nats.go"
)

func NewLogCollectorService(config Config, logger *log.Logger, subject string) (*LogCollectorService, error) {
	ctx, cancel := context.WithCancel(context.Background())

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	// Initialize Centrifugo client if API key is provided
	var centrifugoClient *CentrifugoClient
	if config.Centrifugo.APIKey != "" {
		centrifugoClient, err = NewCentrifugoClient(config.Centrifugo, logger)
		if err != nil {
			logger.Printf("⚠ Warning: Failed to initialize Centrifugo client: %v", err)
			logger.Println("⚠ Continuing without websocket publishing")
		}
	} else {
		logger.Println("⚠ Centrifugo API key not provided, running without websocket publishing")
	}

	var nc *nats.Conn
	if config.Nats.URL != "" {
		nc, err = nats.Connect(config.Nats.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to nats: %w", err)
		}
	} else {
		logger.Println("⚠ NATS URL not provided, running without NATS publishing")
	}

	service := &LogCollectorService{
		config:           config,
		handlers:         make(map[string]*LogFileHandler),
		watcher:          watcher,
		logger:           logger,
		ctx:              ctx,
		cancel:           cancel,
		centrifugoClient: *centrifugoClient,
		natsClient:       &NatsClient{NatsConn: nc},
		subject:          subject,
	}

	return service, nil
}

func (s *LogCollectorService) Start() error {
	centrifugoClient, err := NewCentrifugoClient(s.config.Centrifugo, s.logger)
	if err != nil {
		s.logger.Printf("⚠ Warning: Failed to initialize Centrifugo client: %v", err)
		s.logger.Println("⚠ Continuing without websocket publishing")
	}
	channelID := "logs"

	for _, logFilePath := range s.config.LogFiles {

		absPath, err := filepath.Abs(logFilePath)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path for %s: %w", logFilePath, err)
		}

		//create a directory if it doesnt exist
		logDir := filepath.Dir(absPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", absPath, err)
		}

		//create log file if it doesnt exist
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			file, err := os.Create(absPath)
			if err != nil {
				return fmt.Errorf("failed to create log file %s: %w", absPath, err)
			}
			file.Close()
			s.logger.Printf("created empty log file: %s", absPath)
		}

		//setup handler for this log file
		handler, err := NewLogFileHandler(absPath, s.logger, centrifugoClient, s.natsClient, channelID, s.subject)
		if err != nil {
			return fmt.Errorf("failed to create handler for %s: %w", absPath, err)
		}
		s.handlers[absPath] = handler

		//watching the dir containing the log file
		if err := s.watcher.Add(logDir); err != nil {
			return fmt.Errorf("failed to watch directory %s: %w", logDir, err)
		}

		s.logger.Printf("Started monitoring: %s", absPath)
	}
	s.wg.Add(1)
	go s.processEvents()

	return nil
}

func (s *LogCollectorService) processEvents() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return

		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			//check if this is a file we are monitoring
			if handler, exists := s.handlers[event.Name]; exists {
				//only process write/create events
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					if err := handler.CollectNewLogs(s.logger); err != nil {
						s.logger.Printf("error collecting logs from %s: %v", event.Name, err)
					}
				}
			}

		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			s.logger.Printf("watcher error: %v", err)
		}
	}
}

func (s *LogCollectorService) Stop() {
	s.logger.Println("stopping log collector service...")
	s.cancel()
	s.watcher.Close()

	if s.natsClient != nil {
		s.natsClient.NatsConn.Close()
		s.logger.Println("closed NATS connection")
	}

	s.wg.Wait()
	s.logger.Println("stopped all log monitoring")
}

func (s *LogCollectorService) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	s.logger.Println("Log collector service is running. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	s.logger.Println("Received stop signal")
	s.Stop()

	return nil
}
