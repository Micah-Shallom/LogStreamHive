package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogGenerator struct {
	cfg          Config
	logger       *log.Logger
	userSessions map[string]UserSession
	InBurstMode  bool
	BurstEndTime time.Time
}

func NewLogGenerator(cfg Config) *LogGenerator {
	var writers []io.Writer

	if cfg.OutputFile != "" {
		dir := filepath.Dir(cfg.OutputFile)

		if err := os.MkdirAll(dir, 0777); err != nil {
			log.Fatalf("failed to create log dir %s: %v", dir, err)
		}

		if err := os.Chmod(dir, 0777); err != nil {
			log.Printf("warning: failed to chmod log dir: %v", err)
		}

		if _, err := os.Stat(cfg.OutputFile); err == nil {
			_ = os.Chmod(cfg.OutputFile, 0777)
		}

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
		cfg:          cfg,
		logger:       log.New(multi, "", 0),
		userSessions: make(map[string]UserSession),
		InBurstMode:  false,
		BurstEndTime: time.Now(),
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

	serviceName := lg.cfg.Services[rand.Intn(len(lg.cfg.Services))]
	userID := fmt.Sprintf("user-%d", rand.Intn(9000)+1000)
	requestID := fmt.Sprintf("req-%d-%d", time.Now().Unix(), rand.Intn(9000)+1000)
	duration := rand.Intn(496) + 5
	timestamp := time.Now().Format(time.RFC3339)

	lg.updateUserSession(userID)
	message := lg.createMessageFromPattern(logType, userID)

	log_entry := LogEntry{
		Timestamp: timestamp,
		LogType:   logType,
		Service:   serviceName,
		RequestID: requestID,
		UserID:    userID,
		Duration:  duration,
		Message:   message,
	}

	switch lg.cfg.LogFormat {
	case "json":
		data, err := json.Marshal(log_entry)
		if err != nil {
			panic("unable to marshal log entry")
		}
		return string(data)
	case "csv":
		values := []string{
			timestamp, logType, serviceName, userID, requestID,
			strconv.Itoa(duration), message,
		}
		return strings.Join(values, ",")
	default:
		return fmt.Sprintf("%s [%s] %s [%s] [%s] (%dms): %s",
			timestamp, logType, serviceName, userID, requestID, duration, message)
	}
}

func (lg *LogGenerator) Run(duration float64) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	//running in burst mode
	fmt.Printf("Starting log generator with rate: %d logs/second ", lg.cfg.LogRate)
	fmt.Printf("Log format: %s", lg.cfg.LogFormat)
	fmt.Printf("Burst mode enabled: %v", lg.cfg.EnableBursts)

	start := time.Now()
	count := 0

	for {
		now := time.Now()

		if lg.cfg.EnableBursts {
			if !lg.InBurstMode {
				if rand.Float64() < lg.cfg.BurstFrequency {
					lg.InBurstMode = true
					lg.BurstEndTime = now.Add(time.Duration(lg.cfg.BurstDuration * float64(time.Second)))
					// lg.logger.Printf("⚡⚡ Entering burst mode for %.2f seconds", lg.cfg.BurstDuration)
				}
			} else {
				if now.After(lg.BurstEndTime) {
					lg.InBurstMode = false
					// lg.logger.Println("✓ Exiting burst mode")
				}
			}
		}

		currentRate := float64(lg.cfg.LogRate)
		if currentRate <= 0 {
			currentRate = 1.0
		}
		if lg.InBurstMode {
			currentRate = currentRate * float64(lg.cfg.BurstMultiplier)
		}

		sleep := time.Duration(float64(time.Second) / currentRate)
		if sleep < time.Nanosecond {
			sleep = time.Nanosecond
		}

		LogEntry := lg.generateLogMessage()
		lg.logger.Println(LogEntry)
		count++

		if duration > 0 && time.Since(start).Seconds() >= duration {
			fmt.Printf("Generated %d log entries\n", count)
			break
		}

		time.Sleep(sleep)
	}
}
