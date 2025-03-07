package model

import (
	"fmt"
	"reflect"
	"regexp"
)

// ApacheAccessLog represents an Apache access log entry.
type ApacheErrorLog struct {
	ApacheLog
	Error             string `gorm:"size:255" log:"E" json:"error"`                 // %E
	Filename          string `gorm:"size:255" log:"F" json:"filename"`              // %F
	NumberOfKeepAlive int    `gorm:"default:0" log:"k" json:"number_of_keep_alive"` // %k
	LogLevel          string `gorm:"size:50" log:"l" json:"log_level"`              // %l
	LogID             string `gorm:"size:100" log:"L" json:"log_id"`                // %L
	ConnectionLogID   string `gorm:"size:100" log:"cL" json:"connection_log_id"`    //%{c}L
	ModuleName        string `gorm:"size:100" log:"m" json:"module_name"`           // %m
	LogMessage        string `gorm:"size:500" log:"M" json:"log_message"`           // %M
	ProcessID         int    `gorm:"default:0" log:"P" json:"process_id"`           // %P
	ThreadID          int    `gorm:"default:0" log:"T" json:"thread_id"`            // %T
	SystemThreadID    int    `gorm:"default:0" log:"gT" json:"system_thread_id"`    //%{g}T
}

func (accessLog *ApacheErrorLog) Parse(logLine string) error {
	// Extract key-value pairs from log format: [key:"value"]
	var logPattern = regexp.MustCompile(`\[([A-Za-z-]+):\"([^\"]*)\"]`)

	// Find all matched substrings
	matches := logPattern.FindAllStringSubmatch(logLine, -1)
	if matches == nil {
		return fmt.Errorf("error: %s", "failed to parse logLine")
	}

	entryValue := reflect.ValueOf(accessLog).Elem()

	for _, match := range matches {
		// Extract key-value pair
		key, value := match[1], match[2]

		LoopThroughReflection(entryValue, "log", key, value)
	}

	return nil
}
