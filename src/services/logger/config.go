package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DEBUG    = "DEBUG"
	INFO     = "INFO"
	WARNING  = "WARNING"
	ERROR    = "ERROR"
	CRITICAL = "CRITICAL"
)

type Config struct {
	LogLevel           string        `json:"LOG_LEVEL"`
	LogFormat          string        `json:"LOG_FORMAT"`
	LogInterval        time.Duration `json:"LOG_INTERVAL"`
	LogIntervalSeconds int        `json:"LOG_INTERVAL_SECONDS"`
}

func LoadConfig() Config {
	logLevel := getEnv("LOG_LEVEL", INFO)
	logFormat := getEnv("LOG_FORMAT", "[{time}] [{level}] {message}")
	logFreqStr := getEnv("LOG_INTERVAL_SECONDS", "5")

	intervalSec, err := strconv.Atoi(logFreqStr)
	if err != nil || intervalSec <= 0 {
		log.Printf("Invalid LOG_INTERVAL_SECONDS: %s, defaulting to 5 seconds", &logFreqStr)
		intervalSec = 5
	}

	return Config{
		LogLevel:    logLevel,
		LogFormat:   logFormat,
		LogIntervalSeconds: intervalSec,
		LogInterval: time.Duration(intervalSec) * time.Second,
	}
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

func formatLog(format, level, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	output := strings.ReplaceAll(format, "{time}", timestamp)
	output = strings.ReplaceAll(output, "{level}", level)
	output = strings.ReplaceAll(output, "{message}", message)
	return output
}
