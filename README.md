# Logstash/Vector Experimentation with go-lumber

This project demonstrates how to use **Vector's logstash source** to receive logs via the **Lumberjack protocol** using the **go-lumber** client library from Elastic.

## 📋 Overview

The setup consists of:

- **Vector** - Acting as a log receiver with a logstash source (listening on port 5044)
- **Go Lumber Client** - Sending structured log messages via the Lumberjack protocol
- **Output Sinks** - Console and file outputs to verify received logs

## 🎯 What is This Testing?

### Lumberjack Protocol

The Lumberjack protocol is a binary protocol originally created for Logstash Forwarder (now superseded by Filebeat). It's designed for reliable log shipping with:

- Compression support
- Acknowledgment of received data
- Connection multiplexing
- TLS support (optional)

### Vector's Logstash Source

Vector can act as a Lumberjack protocol server, allowing it to receive logs from:

- Filebeat
- Logstash forwarders
- Custom clients (like our go-lumber implementation)

## 🛠️ Prerequisites

1. **Vector** - Install from https://vector.dev/docs/setup/installation/

   ```bash
   # macOS
   brew install vector

   # Or download binary
   curl --proto '=https' --tlsv1.2 -sSfL https://sh.vector.dev | bash
   ```

2. **Go** - Version 1.21 or later

   ```bash
   # macOS
   brew install go

   # Verify
   go version
   ```

## 🚀 Quick Start

### Step 1: Initialize the Go Module

```bash
# Download dependencies
go mod download
go mod tidy
```

### Step 2: Create Logs Directory

```bash
mkdir -p logs
```

### Step 3: Start Vector

In one terminal window:

```bash
vector --config vector.yaml
```

You should see output indicating Vector is listening:

```
INFO vector::internal_events::vector_started: Vector has started.
INFO vector::sources::logstash: Listening for connections. address=0.0.0.0:5044
```

### Step 4: Run the Go Lumber Sender

In another terminal window:

```bash
# Send 10 messages with 1 second interval (default)
go run sender.go

# Or customize:
go run sender.go -count 20 -interval 500ms
go run sender.go -host localhost:5044 -count 5
```

### Step 5: Observe the Results

You should see:

1. **Sender terminal**: Confirmation of sent messages
2. **Vector terminal**: Real-time JSON logs being processed
3. **logs/ directory**: Daily log files with JSON records

## 📝 Configuration Details

### Vector Configuration (`vector.yaml`)

The configuration defines three components:

#### 1. Source: Logstash Input

```toml
[sources.logstash_input]
  type = "logstash"
  address = "0.0.0.0:5044"
```

- Listens on all interfaces, port 5044 (standard Lumberjack port)
- Accepts unencrypted connections for testing

#### 2. Transform: Parse and Enrich

```toml
[transforms.parse_and_enrich]
  type = "remap"
  inputs = ["logstash_input"]
```

- Parses JSON from the message field
- Adds processing timestamp
- Tags with source type

#### 3. Sinks: Console and File

```toml
[sinks.console_output]
  type = "console"

[sinks.file_output]
  type = "file"
  path = "./logs/vector-output-%Y-%m-%d.log"
```

- Outputs to both stdout and daily rotating log files
- Uses JSON encoding for structured logs

### Go Lumber Client (`sender.go`)

Key features:

- **Compression**: Level 3 compression for efficiency
- **Batching**: Sends logs in batches
- **Timeout**: 30-second timeout for connections
- **Structured Logging**: Sends JSON-formatted log messages

## 🧪 Experimentation Ideas

### 1. Test High Volume

```bash
# Send 1000 messages as fast as possible
go run sender.go -count 1000 -interval 10ms
```

### 2. Test Connection Resilience

```bash
# Start sender, then stop/start Vector mid-stream
go run sender.go -count 100 -interval 2s
# Stop Vector (Ctrl+C), wait a few seconds, restart it
```

### 3. TLS Configuration

Add to `vector.yaml`:

```toml
[sources.logstash_input]
  type = "logstash"
  address = "0.0.0.0:5044"

  [sources.logstash_input.tls]
    enabled = true
    crt_file = "/path/to/server.crt"
    key_file = "/path/to/server.key"
    ca_file = "/path/to/ca.crt"  # Optional: for client cert verification
```

Update Go client to use TLS (requires modifying sender.go to use TLS dialer).

### 4. Add Custom Fields

Modify the `LogMessage` struct in `sender.go`:

```go
type LogMessage struct {
    Timestamp string                 `json:"timestamp"`
    Level     string                 `json:"level"`
    Service   string                 `json:"service"`
    Message   string                 `json:"message"`
    TraceID   string                 `json:"trace_id"`  // New field
    UserID    string                 `json:"user_id"`    // New field
    Metadata  map[string]interface{} `json:"metadata"`
}
```

### 5. Test Different Sinks

Add additional sinks to `vector.yaml`:

```toml
# Elasticsearch sink
[sinks.elasticsearch]
  type = "elasticsearch"
  inputs = ["parse_and_enrich"]
  endpoint = "http://localhost:9200"

# HTTP sink
[sinks.http]
  type = "http"
  inputs = ["parse_and_enrich"]
  uri = "http://localhost:8080/logs"
  encoding.codec = "json"
```

## 📊 Monitoring and Debugging

### Check Vector Metrics

```bash
# Vector exposes metrics on port 9598 by default
curl http://localhost:9598/metrics
```

### Review Log Files

```bash
# Watch logs in real-time
tail -f logs/vector-output-*.log | jq .

# Count log levels
jq -r '.level' logs/vector-output-*.log | sort | uniq -c
```

### Vector Validation

```bash
# Validate config before running
vector validate vector.yaml

# Run with debug logging
vector --config vector.yaml --verbose
```

## 🔍 Troubleshooting

### Connection Refused

- Ensure Vector is running before starting the sender
- Check the port (5044) is not in use: `lsof -i :5044`

### No Logs Appearing

- Verify Vector's console output shows incoming connections
- Check firewall settings
- Ensure logs directory exists and is writable

### Go Module Issues

```bash
# Clean and re-download
go clean -modcache
go mod download
```

## 📚 References

- [Vector Logstash Source Documentation](https://vector.dev/docs/reference/configuration/sources/logstash/)
- [go-lumber GitHub Repository](https://github.com/elastic/go-lumber)
- [Lumberjack Protocol Specification](https://github.com/elastic/logstash-forwarder/blob/master/PROTOCOL.md)
- [Vector Configuration Documentation](https://vector.dev/docs/reference/configuration/)

## 🎓 Learning Points

1. **Protocol Compatibility**: Vector's logstash source is fully compatible with the Lumberjack v2 protocol
2. **Structured Logging**: The protocol works best with structured (JSON) log data
3. **Reliability**: Lumberjack provides acknowledgments, making it suitable for critical logs
4. **Performance**: Compression significantly reduces network bandwidth
5. **Flexibility**: Vector can receive via Lumberjack and output to dozens of different sinks

## 🔄 Next Steps

- Experiment with different Vector transforms (filtering, aggregation)
- Set up TLS for secure log transmission
- Integrate with real log sources (Filebeat, Logstash)
- Test with production-scale log volumes
- Explore Vector's metrics and observability features

---

Happy experimenting! 🚀
