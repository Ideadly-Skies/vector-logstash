package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"logstash-experimentation/pkg/client"
	"logstash-experimentation/pkg/models"
)

func main() {
	// Create client configuration
	config := &client.Config{
		Address:          "localhost:5044",
		Timeout:          30 * time.Second,
		CompressionLevel: 3,
	}

	// Create lumber client
	lumberClient, err := client.NewLumberClient(config)
	if err != nil {
		log.Fatalf("Failed to create lumber client: %v", err)
	}
	defer lumberClient.Close()

	log.Println("🚀 Starting advanced log sender...")
	log.Printf("📊 Connected to %s", lumberClient.Address())
	log.Println("Sending various log patterns...")

	// Simulate realistic log patterns
	for i := 0; i < 30; i++ {
		pattern := i % 3

		switch pattern {
		case 0:
			sendAccessLog(lumberClient, i)
		case 1:
			sendMetricLog(lumberClient, i)
		case 2:
			sendErrorLog(lumberClient, i)
		}

		time.Sleep(500 * time.Millisecond)
	}

	log.Println("✅ All logs sent!")
}

func sendAccessLog(lumberClient *client.LumberClient, seq int) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	paths := []string{"/api/users", "/api/orders", "/api/products", "/health", "/metrics"}
	statusCodes := []int{200, 201, 204, 400, 404, 500}

	accessLog := models.NewAccessLog(
		"api-gateway",
		methods[rand.Intn(len(methods))],
		paths[rand.Intn(len(paths))],
		fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
		statusCodes[rand.Intn(len(statusCodes))],
		rand.Float64()*1000,
		map[string]interface{}{
			"sequence":    seq,
			"request_id":  fmt.Sprintf("req-%d", seq),
			"environment": "testing",
		},
	)
	accessLog.UserAgent = "Mozilla/5.0 (compatible; Test/1.0)"

	if n, err := lumberClient.Send(accessLog); err != nil {
		log.Printf("❌ Error sending access log: %v", err)
	} else {
		log.Printf("✓ Sent access log (acked: %d)", n)
	}
}

func sendMetricLog(lumberClient *client.LumberClient, seq int) {
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

	metricLog := models.NewMetricLog(
		"monitoring",
		metric.name,
		metric.unit,
		metric.value,
		map[string]string{
			"host":        "server-01",
			"environment": "testing",
			"region":      "us-east-1",
		},
		map[string]interface{}{
			"sequence": seq,
		},
	)

	if n, err := lumberClient.Send(metricLog); err != nil {
		log.Printf("❌ Error sending metric log: %v", err)
	} else {
		log.Printf("✓ Sent metric log (acked: %d)", n)
	}
}

func sendErrorLog(lumberClient *client.LumberClient, seq int) {
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

	errorLog := models.NewErrorLog(
		"backend-service",
		selectedError.message,
		selectedError.code,
		fmt.Sprintf("at handler.go:123\nat middleware.go:45\nat main.go:%d", seq),
		map[string]interface{}{
			"sequence":    seq,
			"environment": "testing",
			"trace_id":    fmt.Sprintf("trace-%d", seq),
		},
	)

	if n, err := lumberClient.Send(errorLog); err != nil {
		log.Printf("❌ Error sending error log: %v", err)
	} else {
		log.Printf("✓ Sent error log (acked: %d)", n)
	}
}
