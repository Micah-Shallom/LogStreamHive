package main

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type ErrorSequences []ErrorSequence
type Anomalys []Anomaly

type Statistics struct {
	LogTypeCounts     map[string]int     `json:"logTypeCounts"`
	ServiceDurations  map[string]float64 `json:"serviceDurations"`
	ServiceCallCounts map[string]int     `json:"serviceCallCounts"`
	ErrorSequences    ErrorSequences     `json:"errorSequences"`
	AnomalyDetections Anomalys           `json:"anomaylDetections"`
	UpdatedAt         time.Time          `json:"updatedAt"`
}

type ErrorSequence struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Count     int       `json:"count"`
	Service   string    `json:"service"`
}

type Anomaly struct {
	Timestamp  time.Time `json:"timestamp"`
	Service    string    `json:"service"`
	MetricName string    `json:"metricName"`
	Value      float64   `json:"value"`
	Threshold  float64   `json:"threshold"`
}

func analyzeLogFiles(cfg *App) Statistics {
	stats := Statistics{
		LogTypeCounts:     make(map[string]int),
		ServiceDurations:  make(map[string]float64),
		ServiceCallCounts: make(map[string]int),
		UpdatedAt:         time.Now(),
	}

	content, err := os.ReadFile(cfg.Config.OutputFile)
	if err != nil {
		return stats
	}

	var entries []LogEntry
	for _, line := range strings.Split(string(content), "\n") {
		if line == "" {
			continue
		}
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	analyzeEntries(&stats, entries)
	return stats
}

func analyzeEntries(stats *Statistics, entries []LogEntry) {
	for _, entry := range entries {
		stats.LogTypeCounts[entry.LogType]++

		stats.ServiceDurations[entry.Service] += float64(entry.Duration)
		stats.ServiceCallCounts[entry.Service]++

	}

	for service := range stats.ServiceDurations {
		if count := stats.ServiceCallCounts[service]; count > 0 {
			stats.ServiceDurations[service] /= float64(count)
		}
	}

	stats.ErrorSequences = detectErrorSequences(entries)
	stats.AnomalyDetections = detectAnomalies(entries)
}

func detectErrorSequences(entries []LogEntry) ErrorSequences {
	sequences := ErrorSequences{}
	currentSequence := ErrorSequence{}
	errorCount := 0

	for i, entry := range entries {
		if entry.LogType == "ERROR" {
			errorCount++
			if errorCount == 1 {
				currentSequence.StartTime, _ = time.Parse(time.RFC3339, entry.Timestamp)
				currentSequence.Service = entry.Service
			}
		} else {
			if errorCount >= 3 {
				currentSequence.EndTime, _ = time.Parse(time.RFC3339, entries[i-1].Timestamp)
				currentSequence.Count = errorCount
				sequences = append(sequences, currentSequence)
			}
			errorCount = 0
		}
	}

	return sequences
}

func detectAnomalies(entries []LogEntry) []Anomaly {
	anomalies := []Anomaly{}

	// Calculate baseline metrics
	baselineDuration := calculateBaselineDuration(entries)

	// Detect duration anomalies
	for _, entry := range entries {
		if float64(entry.Duration) > baselineDuration*2 {
			timestamp, _ := time.Parse(time.RFC3339, entry.Timestamp)
			anomalies = append(anomalies, Anomaly{
				Timestamp:  timestamp,
				Service:    entry.Service,
				MetricName: "duration",
				Value:      float64(entry.Duration),
				Threshold:  baselineDuration * 2,
			})
		}
	}

	return anomalies
}

func calculateBaselineDuration(entries []LogEntry) float64 {
	if len(entries) == 0 {
		return 0
	}

	var total float64
	for _, entry := range entries {
		total += float64(entry.Duration)
	}

	return total / float64(len(entries))
}
