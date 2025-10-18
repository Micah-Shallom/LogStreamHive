package models

import (
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

type LogFileHandler struct {
	filePath         string
	lastPosition     int64
	mu               sync.Mutex
	logger           *log.Logger
	centrifugoClient *CentrifugoClient
	natsConn         *nats.Conn
	channelID        string
}

type LogMessage struct {
	Timestamp string `json:"timestamp"`
	FilePath  string `json:"file_path"`
	Line      string `json:"line"`
}
