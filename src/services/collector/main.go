package main

import (
	"log"
)

func main() {

	service, err := NewLogCollectorService("/app/config.yml")
	if err != nil {
		log.Fatalf("failed to create collector service: %v", err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("service error: %v", err)
	}
}
