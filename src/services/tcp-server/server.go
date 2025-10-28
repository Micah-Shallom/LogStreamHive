package main

import (
	"fmt"
	"log"
	"net"

	"github.com/nats-io/nats.go"
)

func NewTCPServer(config ServerConfig, logger *log.Logger) (*TCPServer, error) {
	nc, err := nats.Connect(config.NatsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	logger.Printf("connected to NATS at %s", config.NatsURL)

	return &TCPServer{
		config:   config,
		logger:   logger,
		handlers: make([]*ConnectionHandler, 0),
		natsConn: nc,
	}, nil
}

func (s *TCPServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	s.listener = listener
	s.running.Store(true)

	s.logger.Printf("server started on %s", addr)
	s.logger.Printf("publishing logs to NAT subject: %s", s.config.NatsSubject)

	return s.acceptConnections()
}

func (s *TCPServer) acceptConnections() error {
	for s.running.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.running.Load() {
				fmt.Println("Server is shutting down")
				return nil
			}

			s.logger.Printf("error accepting connection: %v", err)
			continue
		}
		s.logger.Printf("new connection from %s", conn.RemoteAddr())

		handler := NewConnectionHandler(conn, s.config.BufferSize, s.logger, s.natsConn, s.config.NatsSubject)

		s.handlersMu.Lock()
		s.handlers = append(s.handlers, handler)
		s.handlersMu.Unlock()

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			handler.Handle()
			s.removeHandler(handler)
		}()
	}

	return nil
}

func (s *TCPServer) removeHandler(handler *ConnectionHandler) {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	for i, h := range s.handlers {
		if h == handler {
			s.handlers = append(s.handlers[:i], s.handlers[i+1:]...)
			break
		}
	}
}

func (s *TCPServer) Stop() {
	s.logger.Println("stopping server...")
	s.running.Store(false)

	if s.listener != nil {
		s.listener.Close()
	}

	s.handlersMu.Lock()
	for _, handler := range s.handlers {
		handler.Stop()
	}
	s.handlersMu.Unlock()

	s.wg.Wait()

	if s.natsConn != nil {
		s.natsConn.Close()
		s.logger.Println("closed nats connection")
	}

	s.logger.Println("server stopped")
}
