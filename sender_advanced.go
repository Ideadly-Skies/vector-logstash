//go:build advanced
// +build advanced
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	v2 "github.com/elastic/go-lumber/client/v2"
)

// Example showing different log patterns and use cases

// Different log message types
type ErrorLog struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Service    string                 `json:"service"`
	Message    string                 `json:"message"`
	Stacktrace string                 `json:"stacktrace,omitempty"`
	ErrorCode  string                 `json:"error_code,omitempty"`
	Metadata   map[string]interface{} `json:"metadata"`
}

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

func main() {
	// Connect to Vector
	client, err := v2.SyncDial("localhost:5044",
		v2.Timeout(30*time.Second),
		v2.CompressionLevel(3),
	)
	if err != nil {
		log.Fatalf("Failed to create lumber client: %v", err)
	}
	defer client.Close()

	log.Println("🚀 Starting advanced log sender...")
	log.Println("📊 Sending various log patterns...")

	// Simulate realistic log patterns
	for i := 0; i < 30; i++ {
		pattern := i % 3

		switch pattern {
		case 0:
			sendAccessLog(client, i)
		case 1:
			sendMetricLog(client, i)
		case 2:
			sendErrorLog(client, i)
		}

		time.Sleep(500 * time.Millisecond)
	}

	log.Println("✅ All logs sent!")
}

func sendAccessLog(client *v2.SyncClient, seq int) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	paths := []string{"/api/users", "/api/orders", "/api/products", "/health", "/metrics"}
	statusCodes := []int{200, 201, 204, 400, 404, 500}

	accessLog := AccessLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      "INFO",
		Service:    "api-gateway",
		Method:     methods[rand.Intn(len(methods))],
		Path:       paths[rand.Intn(len(paths))],
		StatusCode: statusCodes[rand.Intn(len(statusCodes))],
		Duration:   rand.Float64() * 1000,
		ClientIP:   fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
		UserAgent:  "Mozilla/5.0 (compatible; Test/1.0)",
		Metadata: map[string]interface{}{
			"sequence":    seq,
			"request_id":  fmt.Sprintf("req-%d", seq),
			"environment": "testing",
		},
	}

	sendLog(client, accessLog, "access")
}

func sendMetricLog(client *v2.SyncClient, seq int) {
	metrics := []struct {
		name  string
		unit  string
		value float64
	}{
		{"cpu_usage", "percent", rand.Float64() * 100},
		{"memory_usage", "MB", rand.Float64() * 2048},
		{"request_rate", "req/s", rand.Float64() * 1000},
		{"error_rate", "errors/s", rand.Float64() * 10},
	}

	metric := metrics[rand.Intn(len(metrics))]

	metricLog := MetricLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      "INFO",
		Service:    "monitoring",
		MetricName: metric.name,
		Value:      metric.value,
		Unit:       metric.unit,
		Tags: map[string]string{
			"host":        "server-01",
			"environment": "testing",
			"region":      "us-east-1",
		},
		Metadata: map[string]interface{}{
			"sequence": seq,
		},
	}

	sendLog(client, metricLog, "metric")
}

func sendErrorLog(client *v2.SyncClient, seq int) {
	errors := []struct {
		message string
		code    string
	}{
		{"Database connection timeout", "DB_TIMEOUT"},
		{"Unable to parse request body", "PARSE_ERROR"},
		{"Authentication failed", "AUTH_FAILED"},
		{"Service unavailable", "SERVICE_DOWN"},
	}

	selectedError := errors[rand.Intn(len(errors))]

	errorLog := ErrorLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      "ERROR",
		Service:    "backend-service",
		Message:    selectedError.message,
		ErrorCode:  selectedError.code,
		Stacktrace: fmt.Sprintf("at handler.go:123\nat middleware.go:45\nat main.go:%d", seq),
		Metadata: map[string]interface{}{
			"sequence":    seq,
			"environment": "testing",
			"trace_id":    fmt.Sprintf("trace-%d", seq),
		},
	}

	sendLog(client, errorLog, "error")
}

func sendLog(client *v2.SyncClient, logData interface{}, logType string) {
	msgBytes, err := json.Marshal(logData)
	if err != nil {
		log.Printf("❌ Error marshaling %s log: %v", logType, err)
		return
	}

	batch := []interface{}{
		map[string]interface{}{
			"message":    string(msgBytes),
			"@timestamp": time.Now().Format(time.RFC3339),
		},
	}

	n, err := client.Send(batch)
	if err != nil {
		log.Printf("❌ Error sending %s log: %v", logType, err)
	} else {
		log.Printf("✓ Sent %s log (acked: %d)", logType, n)
	}
}
