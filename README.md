# Logstash/Vector Experimentation

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A production-ready Go application for testing Vector's Logstash source with the Lumberjack protocol using Elastic's go-lumber client.

## 📋 Overview

This project demonstrates industry best practices for:

- **Vector Integration** - Receive logs via the Lumberjack protocol
- **Structured Logging** - Type-safe log models with JSON serialization
- **Client Abstraction** - Clean go-lumber client wrapper
- **Docker Support** - Full containerization with docker-compose
- **Production Patterns** - Proper error handling, configuration management, and project structure

## 🏗️ Project Structure

```
.
├── cmd/                        # Application entrypoints
│   ├── sender/                # Basic log sender
│   └── sender-advanced/       # Advanced patterns sender
├── pkg/                       # Reusable packages
│   ├── client/               # Lumber client wrapper
│   └── models/               # Log data models
├── configs/                   # Configuration files
│   └── vector.yaml           # Vector configuration
├── scripts/                   # Shell scripts
│   ├── setup.sh             # Setup script
│   └── run.sh               # Run script
├── docs/                      # Documentation
│   ├── README.md            # Original detailed guide
│   ├── TROUBLESHOOTING.md   # Debugging guide
│   └── examples/            # Example code
├── Dockerfile                 # Multi-stage Docker build
├── docker-compose.yaml        # Full stack orchestration
├── Makefile                   # Build automation
└── go.mod                     # Go dependencies
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Vector (or Docker)
- Make (optional)

### Using Make (Recommended)

```bash
# Install dependencies
make deps

# Build binaries
make build

# Start Vector in one terminal
make vector-start

# Run basic sender in another terminal
make run

# Or run advanced sender
make run-advanced
```

### Using Docker Compose

```bash
# Start everything (Vector + Sender)
docker-compose up

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### Manual Build

```bash
# Build
go build -o bin/sender ./cmd/sender
go build -o bin/sender-advanced ./cmd/sender-advanced

# Run Vector
vector --config configs/vector.yaml

# Run sender
./bin/sender -host localhost:5044 -count 10 -interval 1s
```

## 📦 Installation

```bash
# Clone repository
git clone <repository-url>
cd logstash-experimentation

# Install dependencies
make deps

# Build
make build
```

## 💻 Usage

### Basic Sender

Sends simple structured log messages:

```bash
# Default: 10 messages, 1s interval
./bin/sender

# Custom parameters
./bin/sender -host localhost:5044 -count 20 -interval 500ms
```

### Advanced Sender

Sends realistic log patterns (access logs, metrics, errors):

```bash
./bin/sender-advanced
```

## 🔧 Configuration

### Vector Configuration

Located at `configs/vector.yaml`:

```yaml
sources:
  logstash_input:
    type: logstash
    address: "0.0.0.0:5044"

transforms:
  parse_and_enrich:
    type: remap
    inputs: [logstash_input]
    source: |
      . = parse_json!(.message)
      .processed_at = now()
      .source_type = "lumberjack"

sinks:
  console_output:
    type: console
    inputs: [parse_and_enrich]
    encoding:
      codec: json
```

### Client Configuration

```go
config := &client.Config{
    Address:          "localhost:5044",
    Timeout:          30 * time.Second,
    CompressionLevel: 3,
}
```

## 🧪 Testing

```bash
# Run tests
make test

# With coverage
make coverage

# View coverage report
open coverage.html
```

## 🎨 Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run all checks
make all
```

## 📊 Monitoring

Vector exposes metrics on port 9598:

```bash
# Check metrics
curl http://localhost:9598/metrics

# Health check
curl http://localhost:9598/health
```

## 🐳 Docker

### Build Image

```bash
make docker-build
```

### Run with Docker Compose

```bash
# Start services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

## 📚 API Documentation

### Models

```go
// Basic log message
log := models.NewLogMessage("INFO", "my-service", "message", metadata)

// Error log with stack trace
errorLog := models.NewErrorLog("service", "error msg", "ERR_CODE", stacktrace, metadata)

// Access log
accessLog := models.NewAccessLog("service", "GET", "/api/users", "127.0.0.1", 200, 45.2, metadata)

// Metric log
metricLog := models.NewMetricLog("service", "cpu_usage", "percent", 75.5, tags, metadata)
```

### Client

```go
// Create client
client, err := client.NewLumberClient(config)
defer client.Close()

// Send single log
n, err := client.Send(logData)

// Send batch
n, err := client.SendBatch([]interface{}{log1, log2, log3})
```

## 🛠️ Development

### Adding New Log Types

1. Define model in `pkg/models/log.go`
2. Add constructor function
3. Use in your sender application

### Extending Client

Add methods to `pkg/client/lumber.go`:

```go
func (lc *LumberClient) SendWithRetry(logData interface{}, retries int) error {
    // Implementation
}
```

## 📖 Additional Documentation

- [Detailed Guide](docs/README.md) - Comprehensive documentation
- [Troubleshooting](docs/TROUBLESHOOTING.md) - Common issues and solutions
- [Vector Documentation](https://vector.dev/docs/)
- [go-lumber Repository](https://github.com/elastic/go-lumber)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details

## 🔗 References

- [Vector Logstash Source](https://vector.dev/docs/reference/configuration/sources/logstash/)
- [Lumberjack Protocol](https://github.com/elastic/logstash-forwarder/blob/master/PROTOCOL.md)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

## 🎯 Makefile Commands

Run `make help` to see all available commands:

```
  build                Build all binaries
  build-sender         Build basic sender
  build-sender-advanced Build advanced sender
  run                  Run basic sender
  run-advanced         Run advanced sender
  test                 Run tests
  coverage             Generate coverage report
  fmt                  Format code
  lint                 Run linter
  clean                Clean build artifacts
  vector-start         Start Vector
  vector-validate      Validate Vector configuration
  docker-build         Build Docker image
  docker-up            Start services with docker-compose
  docker-down          Stop services
  all                  Run all checks and build
```

---

Made with ❤️ for production-ready logging experimentation
