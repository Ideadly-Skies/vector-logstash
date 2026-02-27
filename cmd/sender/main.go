package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"logstash-experimentation/pkg/client"
	"logstash-experimentation/pkg/models"
)

func main() {
	// Command line flags
	host := flag.String("host", "localhost:5044", "Vector logstash host:port")
	count := flag.Int("count", 10, "Number of log messages to send")
	interval := flag.Duration("interval", 1*time.Second, "Interval between messages")
	flag.Parse()

	// Create client configuration
	config := &client.Config{
		Address:          *host,
		Timeout:          30 * time.Second,
		CompressionLevel: 3,
	}

	// Create lumber client
	lumberClient, err := client.NewLumberClient(config)
	if err != nil {
		log.Fatalf("Failed to create lumber client: %v", err)
	}
	defer lumberClient.Close()

	log.Printf("Connected to Vector at %s", lumberClient.Address())
	log.Printf("Sending %d messages with %v interval", *count, *interval)

	// Send log messages
	for i := 0; i < *count; i++ {
		logMsg := models.NewLogMessage(
			getLevelForIndex(i),
			"go-lumber-test",
			fmt.Sprintf("Test message number %d from go-lumber client", i+1),
			map[string]interface{}{
				"sequence":    i + 1,
				"environment": "testing",
				"host":        "experimentation-host",
			},
		)

		n, err := lumberClient.Send(logMsg)
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

// getLevelForIndex returns a log level based on the index to vary log levels
func getLevelForIndex(i int) string {
	levels := []string{"INFO", "DEBUG", "WARN", "ERROR", "INFO", "INFO", "DEBUG", "INFO"}
	return levels[i%len(levels)]
}
