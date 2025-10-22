package main

import "regexp"

type ParsedLog struct {
	Timestamp  string         `json:"timestamp,omitempty"`
	SourceIP   string         `json:"source_ip,omitempty"`
	Method     string         `json:"method,omitempty"`
	Path       string         `json:"path,omitempty"`
	Protocol   string         `json:"protocol,omitempty"`
	StatusCode int            `json:"status_code,omitempty"`
	Size       int            `json:"size,omitempty"`
	User       string         `json:"user,omitempty"`
	Referrer   string         `json:"referrer,omitempty"`
	UserAgent  string         `json:"user_agent,omitempty"`
	Format     string         `json:"format"`
	Raw        string         `json:"raw,omitempty"`
	Message    string         `json:"message,omitempty"`
	LogLevel   string         `json:"log_level,omitempty"`
	SourceFile string         `json:"source_file,omitempty"`
	Extra      map[string]any `json:"extra,omitempty"`
}

type LogParser struct {
	apacheRegex *regexp.Regexp
	nginxRegex  *regexp.Regexp
}

