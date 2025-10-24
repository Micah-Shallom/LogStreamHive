package main

import (
	"encoding/json"
	"regexp"
	"strings"
)

func NewLogParser() *LogParser {
	return &LogParser{
		// Apache: 203.0.113.45 - - [24/Oct/2025:11:02:52 +0000] "POST /users HTTP/1.1" 201 20480 "-" "PostmanRuntime/7.28.4"
		apacheRegex: regexp.MustCompile(`(?P<ip>\S+) - - \[(?P<date>\d+/\w+/\d+):(?P<time>\d+:\d+:\d+) (?P<timezone>[\+\-]\d+)\] "(?P<method>\S+) (?P<path>\S+) HTTP/(?P<http_version>[\d\.]+)" (?P<status>\d+) (?P<size>\d+) "(?P<referrer>[^"]*)" "(?P<user_agent>[^"]*)"`),

		// Nginx: 2025/10/24 11:02:52 [warn] scheduler: User logged in successfully
		nginxRegex: regexp.MustCompile(`(?P<timestamp>\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) \[(?P<level>\w+)\] (?P<process>\w+): (?P<message>.+)`),

		// App: [2025-10-24T11:02:52Z] INFO [inventory-service] User logged in successfully
		appRegex: regexp.MustCompile(`\[(?P<timestamp>[^\]]+)\] (?P<level>\w+) \[(?P<service>[^\]]+)\] (?P<message>.+)`),
	}
}

func (p *LogParser) DetectFormat(logLine string) string {
	logLine = strings.TrimSpace(logLine)

	//try json first
	if len(logLine) > 0 && logLine[0] == '{' {
		var jsonTest map[string]any
		if err := json.Unmarshal([]byte(logLine), &jsonTest); err == nil {
			return "json"
		}
	}

	if p.appRegex.MatchString(logLine) {
		return "app"
	}

	if p.apacheRegex.MatchString(logLine) {
		return "apache"
	}

	if p.nginxRegex.MatchString(logLine) {
		return "nginx"
	}

	return "unknown"
}

func (p *LogParser) Parse(logline string) ParsedLog {
	format := p.DetectFormat(logline)

	switch format {
	case "apache":
		return p.parseApacheLog(logline)
	case "nginx":
		return p.parseNginxLog(logline)
	case "json":
		return p.parseJSONLog(logline)
	case "app":
		return p.parseAppLog(logline)
	default:
		return ParsedLog{
			Raw:    logline,
			Format: "unknown",
		}
	}
}

func (p *LogParser) parseApacheLog(logline string) ParsedLog {
	matches := p.apacheRegex.FindStringSubmatch(logline)
	if len(matches) == 0 {
		return ParsedLog{
			Raw:    logline,
			Format: "unknown",
		}
	}

	result := ParsedLog{Format: "apache"}
	for i, name := range p.apacheRegex.SubexpNames() {
		if i > 0 && i < len(matches) {
			value := matches[i]
			switch name {
			case "ip":
				result.SourceIP = value
			case "date":
				result.dateTemp = value
			case "time":
				result.timeTemp = value
			case "timezone":
				result.timezoneTemp = value
			case "method":
				result.Method = value
			case "path":
				result.Path = value
			case "http_version":
				result.Protocol = "HTTP/" + value
			case "status":
				result.StatusCode, _ = parseIntSafe(value)
			case "size":
				result.Size, _ = parseIntSafe(value)
			case "referrer":
				if value != "-" {
					result.Referrer = value
				}
			case "user_agent":
				result.UserAgent = value
			}
		}
	}

	if result.dateTemp != "" && result.timeTemp != "" {
		result.Timestamp = p.convertApacheDate(result.dateTemp, result.timeTemp, result.timezoneTemp)
	}

	return result
}

func (p *LogParser) parseNginxLog(logLine string) ParsedLog {
	matches := p.nginxRegex.FindStringSubmatch(logLine)
	if len(matches) == 0 {
		return ParsedLog{
			Raw:    logLine,
			Format: "unknown",
		}
	}

	result := ParsedLog{Format: "nginx"}
	for i, name := range p.nginxRegex.SubexpNames() {
		if i > 0 && i < len(matches) {
			value := matches[i]
			switch name {
			case "timestamp":
				// Convert "2025/10/24 11:02:52" to RFC3339
				result.Timestamp = p.convertNginxDate(value)
			case "level":
				result.LogLevel = strings.ToUpper(value)
			case "process":
				result.Process = value
			case "message":
				result.Message = value
			}
		}
	}

	return result
}

func (p *LogParser) parseAppLog(logLine string) ParsedLog {
	matches := p.appRegex.FindStringSubmatch(logLine)
	if len(matches) == 0 {
		return ParsedLog{
			Raw:    logLine,
			Format: "unknown",
		}
	}

	result := ParsedLog{Format: "app"}
	for i, name := range p.appRegex.SubexpNames() {
		if i > 0 && i < len(matches) {
			value := matches[i]
			switch name {
			case "timestamp":
				result.Timestamp = value // Already in RFC3339 format
			case "level":
				result.LogLevel = value
			case "service":
				result.Service = value
			case "message":
				result.Message = value
			}
		}
	}

	return result
}

func (p *LogParser) parseJSONLog(logLine string) ParsedLog {
	var data map[string]any
	if err := json.Unmarshal([]byte(logLine), &data); err != nil {
		return ParsedLog{
			Raw:    logLine,
			Format: "unknown",
		}
	}

	result := ParsedLog{
		Format: "json",
		Extra:  make(map[string]any),
	}

	fieldMappings := map[string]string{
		"ip":         "source_ip",
		"level":      "log_level",
		"log_type":   "log_level",
		"msg":        "message",
		"time":       "timestamp",
		"user_id":    "user_id",
		"request_id": "request_id",
		"duration":   "duration",
	}

	for key, value := range data {
		// Check if this field has a mapping
		if mappedKey, exists := fieldMappings[key]; exists {
			p.setField(&result, mappedKey, value)
		} else {
			// Try to set directly on the struct
			if !p.setField(&result, key, value) {
				// If not a standard field, add to Extra
				result.Extra[key] = value
			}
		}
	}

	return result
}
