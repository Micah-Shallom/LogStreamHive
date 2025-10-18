package main

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	LogType   string `json:"log_type"`
	UserID    string `json:"user_id"`
	Duration  int    `json:"duration"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Service   string `json:"service"`
}
