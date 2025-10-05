package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

type LogCollectorService struct {
	config   Config
	handlers map[string]*LogFileHandler
	watcher  *fsnotify.Watcher
	logger   *log.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewLogCollectorService(configPath string) (*LogCollectorService, error) {
	logger := log.New(os.Stdout, "[LogCollector] ", log.LstdFlags)

	//Load configuarionation
	config, err := loadConfig(configPath, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &LogCollectorService{
		config:   config,
		handlers: make(map[string]*LogFileHandler),
		watcher:  watcher,
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
	}

	return service, nil
}

func (s *LogCollectorService) Start() error {
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
		handler, err := NewLogFileHandler(absPath, s.logger)
		if err != nil {
			return fmt.Errorf("failed to create handler for %s: %w", absPath, err)
		}
		s.handlers[absPath] = handler

		//watching the dir containing the log file
		//note: setup single watcher per directory instead of per file in the future
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
					if err := handler.collectNewLogs(); err != nil {
						s.logger.Printf("error collecting logs from %s: %w", event.Name, err)
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
