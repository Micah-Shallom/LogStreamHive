package services

import (
	"context"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type LogFileHandler struct {
	FilePath         string
	LastPosition     int64
	Mu               sync.Mutex
	Logger           *log.Logger
	CentrifugoClient *CentrifugoClient
	ChannelID        string
	NatsClient       *NatsClient
	Subject          string
}

type LogMessage struct {
	Timestamp string `json:"timestamp"`
	FilePath  string `json:"file_path"`
	Line      string `json:"line"`
}

type LogCollectorService struct {
	config           Config
	handlers         map[string]*LogFileHandler
	watcher          *fsnotify.Watcher
	logger           *log.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
	centrifugoClient CentrifugoClient
	natsClient       *NatsClient
	subject          string
}

type Config struct {
	LogFiles      []string         `yaml:"log_files" json:"log_files"`
	CheckInterval float64          `yaml:"check_interval" json:"check_interval"`
	Centrifugo    CentrifugoConfig `yaml:"centrifugo" json:"centrifugo"`
	ChannelID     string           `yaml:"channel_id" json:"channel_id"`
	Nats          NatsConfig       `yaml:"nats"`
}
