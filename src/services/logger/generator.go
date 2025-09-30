package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogGenerator struct {
	cfg    Config
	logger *log.Logger
}

func NewLogGenerator(cfg Config) *LogGenerator {
	var writers []io.Writer

	if cfg.OutputFile != "" {
		logfile := &lumberjack.Logger{
			Filename:   cfg.OutputFile,
			MaxSize:    1,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		}

		writers = append(writers, logfile)
	}

	if cfg.ConsoleOutput {
		writers = append(writers, os.Stdout)
	}

	multi := io.MultiWriter(writers...)
	return &LogGenerator{
		cfg:    cfg,
		logger: log.New(multi, "", 0),
	}
}

func (lg *LogGenerator) selectLogType() string {
	types := []string{}
	weights := []int{}
	total := 0

	for t, w := range lg.cfg.LogDistribution {
		types = append(types, t)
		weights = append(weights, w)
		total += w
	}

	r := rand.Intn(total)
	for i, w := range weights {
		if r < w {
			return types[i]
		}
		r -= w
	}
	return "INFO"
}

func (lg *LogGenerator) generateLogMessage() string {
	logType := lg.selectLogType()
	messages := map[string][]string{
		"INFO": {
			"User logged in successfully",
			"Page loaded in 0.2 seconds",
			"Database connection established",
			"Cache refreshed successfully",
			"API request completed",
		},
		"WARNING": {
			"High memory usage detected",
			"API response time exceeding threshold",
			"Database connection pool running low",
			"Retry attempt for failed operation",
			"Cache miss rate increasing",
		},
		"ERROR": {
			"Failed to connect to database",
			"API request timeout",
			"Invalid user credentials",
			"Processing error in data pipeline",
			"Out of memory error",
		},
		"DEBUG": {
			"Function X called with parameters Y",
			"SQL query execution details",
			"Cache lookup performed",
			"Request headers processed",
			"Internal state transition",
		},
	}

	msgs := messages[logType]
	message := msgs[rand.Intn(len(msgs))]
	timestamp := time.Now().Format(time.RFC3339)
	logID := fmt.Sprintf("LOG-%d-%d", time.Now().Unix(), rand.Intn(9000)+1000)

	log_entry := LogEntry{
		ID:        logID,
		Message:   message,
		Timestamp: timestamp,
		Level:     logType,
		Source:    "log-generator",
	}

	log, ok := json.Marshal(log_entry)
	if ok != nil {
		panic("unable to marshal log into logentry")
	}

	return string(log)
}

func (lg *LogGenerator) Run(duration float64) {
	sleep := time.Second
	if lg.cfg.LogRate > 0 {
		sleep = time.Duration(1e9 / float64(lg.cfg.LogRate))
	}

	start := time.Now()
	count := 0

	for {
		LogEntry := lg.generateLogMessage()
		lg.logger.Println(LogEntry)
		count++

		time.Sleep(sleep)

		if duration > 0 && time.Since(start).Seconds() >= duration {
			fmt.Printf("Generated %d log entries\n", count)
			break
		}
	}
}
