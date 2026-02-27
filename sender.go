//go:build !advanced
// +build !advanced

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	v2 "github.com/elastic/go-lumber/client/v2"
)

type LogMessage struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func main() {
	// Command line flags
	host := flag.String("host", "localhost:5044", "Vector logstash host:port")
	count := flag.Int("count", 10, "Number of log messages to send")
	interval := flag.Duration("interval", 1*time.Second, "Interval between messages")
	flag.Parse()

	// Create lumber client (v2)
	// For testing without TLS, we use SyncDial which creates a sync client with TCP connection
	client, err := v2.SyncDial(*host,
		v2.Timeout(30*time.Second),
		v2.CompressionLevel(3),
	)
	if err != nil {
		log.Fatalf("Failed to create lumber client: %v", err)
	}
	defer client.Close()

	log.Printf("Connected to Vector at %s", *host)
	log.Printf("Sending %d messages with %v interval", *count, *interval)

	// Send log messages
	for i := 0; i < *count; i++ {
		logMsg := LogMessage{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     getLevelForIndex(i),
			Service:   "go-lumber-test",
			Message:   fmt.Sprintf("Test message number %d from go-lumber client", i+1),
			Metadata: map[string]interface{}{
				"sequence":    i + 1,
				"environment": "testing",
				"host":        "experimentation-host",
			},
		}

		// Convert to JSON
		msgBytes, err := json.Marshal(logMsg)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			continue
		}

		// Send via lumber protocol
		// The message is sent as a map with a "message" field containing our JSON
		batch := []interface{}{
			map[string]interface{}{
				"message":    string(msgBytes),
				"@timestamp": time.Now().Format(time.RFC3339),
			},
		}

		n, err := client.Send(batch)
		if err != nil {
			log.Printf("Error sending message %d: %v", i+1, err)
		} else {
			log.Printf("✓ Sent message %d: level=%s (acked: %d)", i+1, logMsg.Level, n)
		}

		if i < *count-1 {
			time.Sleep(*interval)
		}
	}

	log.Println("All messages sent successfully!")
}

// Helper function to vary log levels
func getLevelForIndex(i int) string {
	levels := []string{"INFO", "DEBUG", "WARN", "ERROR", "INFO", "INFO", "DEBUG", "INFO"}
	return levels[i%len(levels)]
}
