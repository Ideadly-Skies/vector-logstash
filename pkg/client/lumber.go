package client

import (
	"encoding/json"
	"fmt"
	"time"

	v2 "github.com/elastic/go-lumber/client/v2"
)

// LumberClient wraps the go-lumber client with additional functionality
type LumberClient struct {
	client  *v2.SyncClient
	address string
}

// Config holds configuration for the Lumber client
type Config struct {
	Address          string
	Timeout          time.Duration
	CompressionLevel int
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Address:          "localhost:5044",
		Timeout:          30 * time.Second,
		CompressionLevel: 3,
	}
}

// NewLumberClient creates a new Lumber client with the given configuration
func NewLumberClient(config *Config) (*LumberClient, error) {
	if config == nil {
		config = DefaultConfig()
	}

	client, err := v2.SyncDial(config.Address,
		v2.Timeout(config.Timeout),
		v2.CompressionLevel(config.CompressionLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lumber client: %w", err)
	}

	return &LumberClient{
		client:  client,
		address: config.Address,
	}, nil
}

// Send sends a log entry to the server
func (lc *LumberClient) Send(logData interface{}) (int, error) {
	msgBytes, err := json.Marshal(logData)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal log data: %w", err)
	}

	batch := []interface{}{
		map[string]interface{}{
			"message":    string(msgBytes),
			"@timestamp": time.Now().Format(time.RFC3339),
		},
	}

	n, err := lc.client.Send(batch)
	if err != nil {
		return 0, fmt.Errorf("failed to send batch: %w", err)
	}

	return n, nil
}

// SendBatch sends multiple log entries in a single batch
func (lc *LumberClient) SendBatch(logEntries []interface{}) (int, error) {
	if len(logEntries) == 0 {
		return 0, nil
	}

	batch := make([]interface{}, 0, len(logEntries))
	timestamp := time.Now().Format(time.RFC3339)

	for _, entry := range logEntries {
		msgBytes, err := json.Marshal(entry)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal log entry: %w", err)
		}

		batch = append(batch, map[string]interface{}{
			"message":    string(msgBytes),
			"@timestamp": timestamp,
		})
	}

	n, err := lc.client.Send(batch)
	if err != nil {
		return 0, fmt.Errorf("failed to send batch: %w", err)
	}

	return n, nil
}

// Close closes the lumber client connection
func (lc *LumberClient) Close() error {
	if lc.client != nil {
		return lc.client.Close()
	}
	return nil
}

// Address returns the server address the client is connected to
func (lc *LumberClient) Address() string {
	return lc.address
}
