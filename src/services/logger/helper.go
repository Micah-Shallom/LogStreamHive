package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}

	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}

	log.Printf("Invalid int for %s: %s, defaulting to %d", key, valStr, defaultVal)
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valStr := strings.ToLower(getEnv(key, ""))
	if valStr == "" {
		return defaultVal
	}
	if valStr == "true" || valStr == "1" || valStr == "yes" {
		return true
	}
	if valStr == "false" || valStr == "0" || valStr == "no" {
		return false
	}
	log.Printf("Invalid bool for %s: %s, defaulting to %v", key, valStr, defaultVal)
	return defaultVal
}

func getEnvAsSlice(key string, defaultVal []string, sep string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	parts := strings.Split(valStr, sep)
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

func getEnvAsFloat(key string, defaultVal float64) float64 {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	log.Printf("Invalid float for %s: %s, defaulting to %f", key, valStr, defaultVal)
	return defaultVal
}

