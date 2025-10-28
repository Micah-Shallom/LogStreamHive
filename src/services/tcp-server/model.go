package main

import (
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/nats-io/nats.go"
)

type ConnectionHandler struct {
	conn        net.Conn
	bufferSize  int
	logger      *log.Logger
	running     atomic.Bool
	natsConn    *nats.Conn
	natsSubject string
}

type ServerConfig struct {
	Host        string
	Port        int
	BufferSize  int
	NatsURL     string
	NatsSubject string
}

type TCPServer struct {
	config     ServerConfig
	listener   net.Listener
	logger     *log.Logger
	running    atomic.Bool
	wg         sync.WaitGroup
	handlers   []*ConnectionHandler
	handlersMu sync.Mutex
	natsConn   *nats.Conn
}
