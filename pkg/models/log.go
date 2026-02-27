package models

import "time"

// LogMessage represents a basic log entry
type LogMessage struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ErrorLog represents an error log entry with stack trace
type ErrorLog struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Service    string                 `json:"service"`
	Message    string                 `json:"message"`
	Stacktrace string                 `json:"stacktrace,omitempty"`
	ErrorCode  string                 `json:"error_code,omitempty"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// AccessLog represents an HTTP access log entry
type AccessLog struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Service    string                 `json:"service"`
	Method     string                 `json:"method"`
	Path       string                 `json:"path"`
	StatusCode int                    `json:"status_code"`
	Duration   float64                `json:"duration_ms"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	ClientIP   string                 `json:"client_ip"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// MetricLog represents a metric data point
type MetricLog struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Service    string                 `json:"service"`
	MetricName string                 `json:"metric_name"`
	Value      float64                `json:"value"`
	Unit       string                 `json:"unit"`
	Tags       map[string]string      `json:"tags"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// NewLogMessage creates a new basic log message
func NewLogMessage(level, service, message string, metadata map[string]interface{}) *LogMessage {
	return &LogMessage{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Service:   service,
		Message:   message,
		Metadata:  metadata,
	}
}

// NewErrorLog creates a new error log entry
func NewErrorLog(service, message, errorCode, stacktrace string, metadata map[string]interface{}) *ErrorLog {
	return &ErrorLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      "ERROR",
		Service:    service,
		Message:    message,
		ErrorCode:  errorCode,
		Stacktrace: stacktrace,
		Metadata:   metadata,
	}
}

// NewAccessLog creates a new access log entry
func NewAccessLog(service, method, path, clientIP string, statusCode int, duration float64, metadata map[string]interface{}) *AccessLog {
	return &AccessLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      getAccessLogLevel(statusCode),
		Service:    service,
		Method:     method,
		Path:       path,
		StatusCode: statusCode,
		Duration:   duration,
		ClientIP:   clientIP,
		Metadata:   metadata,
	}
}

// NewMetricLog creates a new metric log entry
func NewMetricLog(service, metricName, unit string, value float64, tags map[string]string, metadata map[string]interface{}) *MetricLog {
	return &MetricLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      "INFO",
		Service:    service,
		MetricName: metricName,
		Value:      value,
		Unit:       unit,
		Tags:       tags,
		Metadata:   metadata,
	}
}

// getAccessLogLevel returns appropriate log level based on HTTP status code
func getAccessLogLevel(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "ERROR"
	case statusCode >= 400:
		return "WARN"
	default:
		return "INFO"
	}
}
