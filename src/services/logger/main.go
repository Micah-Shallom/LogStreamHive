package main

import (
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// simple logger service to test environment is working
	cfg := LoadConfig()

	logfile := &lumberjack.Logger{
		Filename:   "/var/log/logger/service.log",
		MaxSize:    1, //megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   true,
	}

	multiwriter := io.MultiWriter(os.Stdout, logfile)
	logger := log.New(multiwriter, "", 0)

	logger.Println(formatLog(cfg.LogFormat, cfg.LogLevel, "Starting distributed logger service..."))

	for {
		logger.Println(formatLog(cfg.LogFormat, cfg.LogLevel, "Logger service is running. This will be part of our distributed system!"))
		time.Sleep(cfg.LogInterval)
	}
}
