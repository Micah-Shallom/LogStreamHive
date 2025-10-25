package main

import (
	"encoding/json"
	"strconv"
)

func getStringValue(m map[string]any, key string) (string, bool) {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal, true
		}
	}
	return "", false
}

// Helper function to safely extract int values from nested maps
func getIntValue(m map[string]any, key string) (int, bool) {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		case string:
			strconvVal, err := strconv.Atoi(v)
			if err == nil {
				return strconvVal, true
			}
			return 0, false
		}
	}
	return 0, false
}

func getIndexKeys(parsedData []LogEntry) []string {
	keySet := make(map[string]bool)

	for _, entry := range parsedData {
		if _, ok := getStringValue(entry, "timestamp"); ok {
			keySet["date"] = true
		}
		if _, ok := getStringValue(entry, "format"); ok {
			keySet["format"] = true
		}
		if _, ok := getStringValue(entry, "log_type"); ok {
			keySet["level"] = true
		}
		if _, ok := getStringValue(entry, "log_level"); ok {
			keySet["level"] = true
		}
		if _, ok := getStringValue(entry, "service"); ok {
			keySet["service"] = true
		}
		if _, ok := getIntValue(entry, "status_code"); ok {
			keySet["status"] = true
		}
		if _, ok := getStringValue(entry, "method"); ok {
			keySet["method"] = true
		}
		if _, ok := getStringValue(entry, "source_ip"); ok {
			keySet["ip"] = true
		}
		if _, ok := getStringValue(entry, "user_id"); ok {
			keySet["user"] = true
		}

		// Check nested in extra
		if extra, ok := entry["extra"].(map[string]any); ok {
			if lineStr, ok := extra["line"].(string); ok {
				var nestedLog map[string]any
				if err := json.Unmarshal([]byte(lineStr), &nestedLog); err == nil {
					if _, ok := getStringValue(nestedLog, "log_type"); ok {
						keySet["level"] = true
					}
					if _, ok := getStringValue(nestedLog, "service"); ok {
						keySet["service"] = true
					}
					if _, ok := getIntValue(nestedLog, "status_code"); ok {
						keySet["status"] = true
					}
					if _, ok := getStringValue(nestedLog, "user_id"); ok {
						keySet["user"] = true
					}
				}
			}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	return keys
}
