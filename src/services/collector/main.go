package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/toolkits/file"
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
		//create a directory if it doesnt exist
		logDir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", logFilePath, err)
		}

		//create log file if it doesnt exist
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			file, err := os.Create(logFilePath)
			if err != nil {
				return fmt.Errorf("failed to create log file %s: %w", logFilePath, err)
			}
			file.Close()
			s.logger.Printf("created empty log file: %s", logFilePath)
		}

		//setup handler for this log file
		handler, err := NewLogFileHandler(logFilePath, s.logger)
		if err != nil {
			return fmt.Errorf("failed to create handler for %s: %w", logFilePath, err)
		}
		s.handlers[logFilePath] = handler

		//watching the dir containing the log file
		if err := s.watcher.Add(logDir); err != nil {
			return fmt.Errorf("failed to watch directory %s: %w",logDir, err)
		}

		s.logger.Printf("Started monitoring: %s", logFilePath)
	}
	s.wg.Add(1)
	go s.processEvents()

	return nil
}
