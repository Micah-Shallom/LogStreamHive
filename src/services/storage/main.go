package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	logStorage, err := NewLogStorage(
		"./logs",
		"application.log",
		NewSizeBasedRotationPolicy(1024), // 1 KB
		NewCountBasedRetentionPolicy(5),
		true,
	)
	if err != nil {
		fmt.Printf("Error creating log storage: %v\n", err)
		return
	}
	defer logStorage.Close()

	logLevels := []string{"INFO", "DEBUG", "WARNING", "ERROR"}
	messages := map[string][]string{
		"ERROR": {
			"Connection refused",
			"Database query timeout",
			"Failed to process request",
			"Authentication failed",
		},
		"WARNING": {
			"High memory usage detected",
			"Slow query performance",
			"Rate limit approaching",
			"Retry attempt #3",
		},
	}

	fmt.Println("Starting log generation...")
	for i := 0; i < 1000; i++ {
		level := logLevels[rand.Intn(len(logLevels))]
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		message := fmt.Sprintf("Test message %d", i)
		if msgs, ok := messages[level]; ok {
			message = msgs[rand.Intn(len(msgs))]
		}

		logEntry := fmt.Sprintf("%s [%s] %s", timestamp, level, message)
		if err := logStorage.WriteLog(logEntry); err != nil {
			fmt.Printf("Error writing log: %v\n", err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("Generated %d log entries...\n", i+1)
		}

		time.Sleep(time.Duration(10+rand.Intn(40)) * time.Millisecond)
	}
	fmt.Println("Log generation complete!")
}