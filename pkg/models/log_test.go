package models

import (
	"testing"
	"time"
)

func TestNewLogMessage(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		service  string
		message  string
		metadata map[string]interface{}
	}{
		{
			name:    "basic log message",
			level:   "INFO",
			service: "test-service",
			message: "test message",
			metadata: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:     "nil metadata",
			level:    "ERROR",
			service:  "test-service",
			message:  "error occurred",
			metadata: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := NewLogMessage(tt.level, tt.service, tt.message, tt.metadata)

			if log == nil {
				t.Fatal("NewLogMessage returned nil")
			}

			if log.Level != tt.level {
				t.Errorf("Level = %v, want %v", log.Level, tt.level)
			}

			if log.Service != tt.service {
				t.Errorf("Service = %v, want %v", log.Service, tt.service)
			}

			if log.Message != tt.message {
				t.Errorf("Message = %v, want %v", log.Message, tt.message)
			}

			// Check timestamp is valid RFC3339
			_, err := time.Parse(time.RFC3339, log.Timestamp)
			if err != nil {
				t.Errorf("Timestamp is not valid RFC3339: %v", err)
			}
		})
	}
}

func TestNewErrorLog(t *testing.T) {
	log := NewErrorLog("service", "error message", "ERR001", "stacktrace", map[string]interface{}{"key": "value"})

	if log == nil {
		t.Fatal("NewErrorLog returned nil")
	}

	if log.Level != "ERROR" {
		t.Errorf("Level = %v, want ERROR", log.Level)
	}

	if log.ErrorCode != "ERR001" {
		t.Errorf("ErrorCode = %v, want ERR001", log.ErrorCode)
	}

	if log.Stacktrace != "stacktrace" {
		t.Errorf("Stacktrace = %v, want stacktrace", log.Stacktrace)
	}
}

func TestNewAccessLog(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedLevel string
	}{
		{"2xx success", 200, "INFO"},
		{"3xx redirect", 301, "INFO"},
		{"4xx client error", 404, "WARN"},
		{"5xx server error", 500, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := NewAccessLog("service", "GET", "/api/test", "127.0.0.1", tt.statusCode, 100.5, nil)

			if log == nil {
				t.Fatal("NewAccessLog returned nil")
			}

			if log.Level != tt.expectedLevel {
				t.Errorf("Level = %v, want %v for status %d", log.Level, tt.expectedLevel, tt.statusCode)
			}

			if log.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %v, want %v", log.StatusCode, tt.statusCode)
			}

			if log.Method != "GET" {
				t.Errorf("Method = %v, want GET", log.Method)
			}
		})
	}
}

func TestNewMetricLog(t *testing.T) {
	tags := map[string]string{"env": "test"}
	metadata := map[string]interface{}{"region": "us-east-1"}

	log := NewMetricLog("monitoring", "cpu_usage", "percent", 75.5, tags, metadata)

	if log == nil {
		t.Fatal("NewMetricLog returned nil")
	}

	if log.Level != "INFO" {
		t.Errorf("Level = %v, want INFO", log.Level)
	}

	if log.MetricName != "cpu_usage" {
		t.Errorf("MetricName = %v, want cpu_usage", log.MetricName)
	}

	if log.Value != 75.5 {
		t.Errorf("Value = %v, want 75.5", log.Value)
	}

	if log.Unit != "percent" {
		t.Errorf("Unit = %v, want percent", log.Unit)
	}

	if len(log.Tags) != 1 {
		t.Errorf("Tags length = %v, want 1", len(log.Tags))
	}
}

func TestGetAccessLogLevel(t *testing.T) {
	tests := []struct {
		statusCode int
		want       string
	}{
		{100, "INFO"},
		{200, "INFO"},
		{201, "INFO"},
		{301, "INFO"},
		{400, "WARN"},
		{404, "WARN"},
		{500, "ERROR"},
		{503, "ERROR"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := getAccessLogLevel(tt.statusCode)
			if got != tt.want {
				t.Errorf("getAccessLogLevel(%d) = %v, want %v", tt.statusCode, got, tt.want)
			}
		})
	}
}
