package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/nats-io/nats.go"
)

func NewConnectionHandler(conn net.Conn, bufferSize int, logger *log.Logger, natsConn *nats.Conn, natsSubject string) *ConnectionHandler {
	return &ConnectionHandler{
		conn:        conn,
		bufferSize:  bufferSize,
		logger:      logger,
		natsConn:    natsConn,
		natsSubject: natsSubject,
	}
}

func (h *ConnectionHandler) Handle() {
	h.running.Store(true)
	defer h.Stop()

	remoteAddr := h.conn.RemoteAddr().String()
	h.logger.Printf("handling connection from %s", remoteAddr)

	reader := bufio.NewReader(h.conn)

	for h.running.Load() {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				h.logger.Printf("error reading from %s: %v", remoteAddr, err)
			} else {
				h.logger.Printf("connection closed by %s", remoteAddr)
			}
			break
		}

		if err := h.processData(line); err != nil {
			h.logger.Printf("error processing data from %s: %v", remoteAddr, err)
			continue
		}

		//send acknowledgment
		if _, err := h.conn.Write([]byte("ACK\n")); err != nil {
			h.logger.Printf("error sending ACK to %s: %v", remoteAddr, err)
			break
		}
	}
}

func (h *ConnectionHandler) processData(data string) error {
	h.logger.Printf("Received log: %s", data)

	var logJSON map[string]any
	if err := json.Unmarshal([]byte(data), &logJSON); err == nil {
		h.logger.Printf("parsed JSON log: %v", logJSON)
	} else {
		h.logger.Printf("received plain text log: %s", data)
	}

	if h.natsConn != nil {
		if err := h.natsConn.Publish(h.natsSubject, []byte(data)); err != nil {
			return fmt.Errorf("failed to pusblish to NATS: %w", err)
		}
		h.logger.Printf("published log to NATS subject: %s", h.natsSubject)
	}

	return nil
}

func (h *ConnectionHandler) Stop() {
	if !h.running.Load(){
		return
	}

	h.running.Store(false)

	if h.conn != nil {
		h.conn.Close()
		h.logger.Printf("closed connection to %s", h.conn.RemoteAddr())
	}
}