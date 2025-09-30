package main

import (
	"strings"
)

const (
	DEBUG    = "DEBUG"
	INFO     = "INFO"
	WARNING  = "WARNING"
	ERROR    = "ERROR"
	CRITICAL = "CRITICAL"
)

var defaultLogDistribution = map[string]int{
	"INFO":    70,
	"WARNING": 20,
	"ERROR":   5,
	"DEBUG":   5,
}

type Config struct {
	LogRate         int            `json:"LOG_RATE"`
	LogTypes        []string       `json:"LOG_TYPES"`
	LogDistribution map[string]int `json:"LOG_DISTRIBUTION"`
	OutputFile      string         `json:"OUTPUT_FILE"`
	ConsoleOutput   bool           `json:"CONSOLE_OUTPUT"`

	LogFormat       string   `json:"LOG_FORMAT"`
	Services        []string `json:"SERVICES"`
	EnableBursts    bool     `json:"ENABLE_BURSTS"`
	BurstFrequency  float64  `json:"BURST_FREQUENCY"`
	BurstMultiplier int      `json:"BURST_MULTIPLIER"`
	BurstDuration   float64  `json:"BURST_DURATION"`
}

func LoadConfig() Config {

	//setting up my defaults
	defaultRate := 5
	defaultTypes := []string{INFO, WARNING, ERROR, DEBUG}
	defaultOutputFile := "./logs/service.log"
	defaultConsole := true
	defaultFormat := "text"
	defaultServices := []string{"user-service", "payment-service", "inventory-service", "notification-service"}

	//reading from env
	rate := getEnvAsInt("LOG_RATE", defaultRate)
	types := getEnvAsSlice("LOG_TYPES", defaultTypes, ",")
	console := getEnvAsBool("CONSOLE_OUTPUT", defaultConsole)
	services := getEnvAsSlice("SERVICES", defaultServices, ",")
	logFormat := getEnv("LOG_FORMAT", defaultFormat)
	outputFile := getEnv("OUTPUT_FILE", defaultOutputFile)

	distribution := make(map[string]int)
	for _, t := range types {
		key := "LOG_DIST" + strings.ToUpper(t)
		defaultVal := defaultLogDistribution[strings.ToUpper(t)]
		distribution[strings.ToUpper(t)] = getEnvAsInt(key, defaultVal)
	}

	enableBursts := getEnvAsBool("ENABLE_BURSTS", true)
	burstDuration := getEnvAsFloat("BURST_DURATION", 3.0)
	burstFrequency := getEnvAsFloat("BURST_FREQUENCY", 0.05)
	burstMultiplier := getEnvAsInt("BURST_MULTIPLIER", 5)

	return Config{
		LogRate:         rate,
		LogTypes:        types,
		LogDistribution: distribution,
		OutputFile:      outputFile,
		ConsoleOutput:   console,
		LogFormat:       logFormat,
		Services:        services,
		EnableBursts:    enableBursts,
		BurstFrequency:  burstFrequency,
		BurstMultiplier: burstMultiplier,
		BurstDuration:   burstDuration,
	}
}
