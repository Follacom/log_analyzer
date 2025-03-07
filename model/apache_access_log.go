package model

import (
	"fmt"
	"reflect"
	"regexp"
)

// ApacheAccessLog represents an Apache access log entry.
type ApacheAccessLog struct {
	ApacheLog
	ResponseSize      int    `log:"B" json:"response_size"`        // %B
	TimeToServe       int    `log:"D" json:"time_to_serve"`        // %D
	Filename          string `log:"f" json:"filename"`             // %f
	Protocol          string `log:"H" json:"protocol"`             // %H
	NumberOfKeepAlive int    `log:"k" json:"number_of_keep_alive"` // %k
	RequestLogID      string `log:"L" json:"request_log_id"`       // %L
	RequestMethod     string `log:"m" json:"request_method"`       // %m
	Port              int    `log:"p" json:"port"`                 // %p
	ProcessID         int    `log:"P" json:"process_id"`           // %P
	Query             string `log:"q" json:"query"`                // %q
	ResponseHandler   string `log:"R" json:"response_handler"`     // %R
	Status            int    `log:"s" json:"status"`               // %>s
	RemoteUser        string `log:"u" json:"remote_user"`          // %u
	RequestURL        string `log:"U" json:"request_url"`          // %U
	ConnectionStatus  string `log:"X" json:"connection_status"`    // %X
}

func (accessLog *ApacheAccessLog) Parse(logLine string) error {
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
