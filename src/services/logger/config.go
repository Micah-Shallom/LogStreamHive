package main

import (
	"strings"
)

var defaultLogDistribution = map[string]int{
	"INFO":    70,
	"WARNING": 20,
	"ERROR":   5,
	"DEBUG":   5,
}

const (
	DEBUG    = "DEBUG"
	INFO     = "INFO"
	WARNING  = "WARNING"
	ERROR    = "ERROR"
	CRITICAL = "CRITICAL"
)

// const DefaultOutputFile = "/logs/generated_logs.log"

type Config struct {
	LogRate         int            `json:"LOG_RATE"`
	LogTypes        []string       `json:"LOG_TYPES"`
	LogDistribution map[string]int `json:"LOG_DISTRIBUTION"`
	OutputFile      string         `json:"OUTPUT_FILE"`
	ConsoleOutput   bool           `json:"CONSOLE_OUTPUT"`
}

func LoadConfig() Config {

	//setting up my defaults
	defaultRate := 5
	defaultTypes := []string{INFO, WARNING, ERROR, DEBUG}
	defaultOutputFile := "./logs/service.log"
	defaultConsole := true

	//reading from env
	rate := getEnvAsInt("LOG_RATE", defaultRate)
	types := getEnvAsSlice("LOG_TYPES", defaultTypes, ",")
	outputFile := getEnv("OUTPUT_FILE", defaultOutputFile)
	console := getEnvAsBool("CONSOLE_OUTPUT", defaultConsole)

	distribution := make(map[string]int)
	for _, t := range types {
		key := "LOG_DIST" + strings.ToUpper(t)
		defaultVal := defaultLogDistribution[strings.ToUpper(t)]
		distribution[strings.ToUpper(t)] = getEnvAsInt(key, defaultVal)
	}

	return Config{
		LogRate:         rate,
		LogTypes:        types,
		LogDistribution: distribution,
		OutputFile:      outputFile,
		ConsoleOutput:   console,
	}
}
