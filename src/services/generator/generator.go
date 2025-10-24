package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

var sampleData = map[string][]string{
	"ip":       {"192.168.1.1", "10.0.0.2", "172.16.254.1", "8.8.8.8", "1.1.1.1", "203.0.113.45", "198.51.100.23"},
	"endpoint": {"users", "products", "orders", "auth", "search", "health", "api/v1", "dashboard", "settings"},
	"process":  {"worker1", "worker2", "worker3", "main", "background", "scheduler"},
	"method":   {"GET", "POST", "PUT", "DELETE", "PATCH"},
	"protocol": {"HTTP/1.1", "HTTP/2.0"},
	"useragent": {
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"curl/7.68.0",
		"PostmanRuntime/7.28.4",
	},
}

var statusCodes = []int{200, 200, 200, 201, 204, 301, 304, 400, 401, 403, 404, 500, 502, 503}
var responseSizes = []int{100, 250, 500, 1024, 2048, 4096, 8192, 10240, 15360, 20480}

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

func (lg *LogGenerator) generateApacheLog() string {
	ip := sampleData["ip"][rand.Intn(len(sampleData["ip"]))]
	timestamp := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	method := sampleData["method"][rand.Intn(len(sampleData["method"]))]
	endpoint := sampleData["endpoint"][rand.Intn(len(sampleData["endpoint"]))]
	protocol := sampleData["protocol"][rand.Intn(len(sampleData["protocol"]))]
	status := statusCodes[rand.Intn(len(statusCodes))]
	size := responseSizes[rand.Intn(len(responseSizes))]
	referer := "-"
	userAgent := sampleData["useragent"][rand.Intn(len(sampleData["useragent"]))]

	return fmt.Sprintf(`%s - - [%s] "%s /%s %s" %d %d "%s" "%s"`,
		ip, timestamp, method, endpoint, protocol, status, size, referer, userAgent)
}

func (lg *LogGenerator) generateNginxLog() string {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	level := lg.selectLogType()
	process := sampleData["process"][rand.Intn(len(sampleData["process"]))]

	// Map log types to nginx-style levels
	nginxLevel := map[string]string{
		"DEBUG":    "debug",
		"INFO":     "info",
		"WARNING":  "warn",
		"ERROR":    "error",
		"CRITICAL": "crit",
	}[level]

	userID := fmt.Sprintf("user-%d", rand.Intn(9000)+1000)
	lg.updateUserSession(userID)
	message := lg.createMessageFromPattern(level, userID)

	return fmt.Sprintf("%s [%s] %s: %s",
		timestamp, nginxLevel, process, message)
}

func (lg *LogGenerator) generateLogMessage() string {
	generators := []func() string{
		lg.generateApacheLog,
		lg.generateAppLog,
		lg.generateNginxLog,
		lg.generateJSONLog,
	}

	index := rand.Intn(len(generators))

	return generators[index]()
}

func (lg *LogGenerator) generateAppLog() string {
	logType := lg.selectLogType()
	serviceName := lg.cfg.Services[rand.Intn(len(lg.cfg.Services))]
	userID := fmt.Sprintf("user-%d", rand.Intn(9000)+1000)
	timestamp := time.Now().Format(time.RFC3339)

	lg.updateUserSession(userID)
	message := lg.createMessageFromPattern(logType, userID)

	return fmt.Sprintf("[%s] %s [%s] %s",
		timestamp, logType, serviceName, message)
}

func (lg *LogGenerator) generateJSONLog() string {
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

	data, err := json.Marshal(log_entry)
	if err != nil {
		panic("unable to marshal log entry")
	}
	return string(data)
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
		lg.logger.Println(LogEntry) //this is where the log is actually written
		count++

		if duration > 0 && time.Since(start).Seconds() >= duration {
			fmt.Printf("Generated %d log entries\n", count)
			break
		}

		time.Sleep(sleep)
	}
}
