package main

import (
	"encoding/json"
	"regexp"
)


func NewLogParser() *LogParser {
	return &LogParser{
		apacheRegex: regexp.MustCompile(`(\S+) - - \[(\d+/\w+/\d+):(\d+:\d+:\d+) ([\+\-]\d+)\] "(\S+) (\S+) ([^"]+)" (\d+) (\d+|-)`),
		nginxRegex:  regexp.MustCompile(`(\S+) - (\S+) \[(\d+/\w+/\d+):(\d+:\d+:\d+) ([\+\-]\d+)\] "(\S+) (\S+) ([^"]+)" (\d+) (\d+) "([^"]*)" "([^"]*)"`),
	}
}

func (p *LogParser) DetectFormat(logLine string) string {
	//try json first
	var jsonTest map[string]any
	if err := json.Unmarshal([]byte(logLine), &jsonTest); err == nil {
		return "json"
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
	default:
		return ParsedLog{
			Raw:    logline,
			Format: "unknown",
		}
	}
}

func (p *LogParser) parseApacheLog(logline string) ParsedLog {
	matches := p.apacheRegex.FindStringSubmatch(logline)
	if len(matches) < 10 {
		return ParsedLog{
			Raw:    logline,
			Format: "unknown",
		}
	}

	ip := matches[1]
	date := matches[2]
	timeStr := matches[3]
	timezone := matches[4]
	method := matches[5]
	path := matches[6]
	protocol := matches[7]
	statusStr := matches[8]
	sizeStr := matches[9]

	var statusCode int
	if _, err := parseIntSafe(statusStr); err == nil {
		statusCode, _ = parseIntSafe(statusStr)
	}

	var size int
	if sizeStr != "-" {
		size, _ = parseIntSafe(sizeStr)
	}

	timestamp := p.convertApacheDate(date, timeStr, timezone)

	return ParsedLog{
		Timestamp:  timestamp,
		SourceIP:   ip,
		Method:     method,
		Path:       path,
		Protocol:   protocol,
		StatusCode: statusCode,
		Size:       size,
		Format:     "apache",
	}
}

func (p *LogParser) parseNginxLog(logLine string) ParsedLog {
	matches := p.nginxRegex.FindStringSubmatch(logLine)
	if len(matches) < 13 {
		return ParsedLog{
			Raw:    logLine,
			Format: "unknown",
		}
	}

	ip := matches[1]
	user := matches[2]
	date := matches[3]
	timeStr := matches[4]
	timezone := matches[5]
	method := matches[6]
	path := matches[7]
	protocol := matches[8]
	statusStr := matches[9]
	sizeStr := matches[10]
	referrer := matches[11]
	userAgent := matches[12]

	statusCode, _ := parseIntSafe(statusStr)
	size, _ := parseIntSafe(sizeStr)

	timestamp := p.convertApacheDate(date, timeStr, timezone)

	result := ParsedLog{
		Timestamp:  timestamp,
		SourceIP:   ip,
		Method:     method,
		Path:       path,
		Protocol:   protocol,
		StatusCode: statusCode,
		Size:       size,
		UserAgent:  userAgent,
		Format:     "nginx",
	}

	if user != "-" {
		result.User = user
	}
	if referrer != "-" {
		result.Referrer = referrer
	}

	return result
}


func (p *LogParser) parseJSONLog(logLine string) ParsedLog {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(logLine), &data); err != nil {
		return ParsedLog{
			Raw:    logLine,
			Format: "unknown",
		}
	}

	result := ParsedLog{
		Format: "json",
		Extra:  make(map[string]interface{}),
	}

	// Map common fields to our schema
	fieldMappings := map[string]string{
		"ip":    "source_ip",
		"level": "log_level",
		"msg":   "message",
		"time":  "timestamp",
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