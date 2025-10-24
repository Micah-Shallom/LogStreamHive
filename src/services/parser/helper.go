package main

import (
	"encoding/json"
	"strings"
	"time"
)

func parseIntSafe(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0, nil
	}

	var result int
	err := json.Unmarshal([]byte(s), &result)
	return result, err
}

func (p *LogParser) convertApacheDate(date, timeStr, timezone string) string {
	// Format: 10/Oct/2000:13:55:36 +0000
	dateTimeStr := date + ":" + timeStr + " " + timezone
	layout := "02/Jan/2006:15:04:05 -0700"

	t, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		// Return original format if parsing fails
		return dateTimeStr
	}

	return t.Format(time.RFC3339)
}

func (p *LogParser) setField(log *ParsedLog, fieldName string, value interface{}) bool {
	strValue, _ := value.(string)
	intValue, _ := value.(float64)

	switch fieldName {
	case "timestamp":
		log.Timestamp = strValue
	case "source_ip":
		log.SourceIP = strValue
	case "method":
		log.Method = strValue
	case "path":
		log.Path = strValue
	case "protocol":
		log.Protocol = strValue
	case "status_code":
		log.StatusCode = int(intValue)
	case "size":
		log.Size = int(intValue)
	case "user":
		log.User = strValue
	case "referrer":
		log.Referrer = strValue
	case "user_agent":
		log.UserAgent = strValue
	case "message":
		log.Message = strValue
	case "log_level":
		log.LogLevel = strValue
	default:
		return false
	}
	return true
}
