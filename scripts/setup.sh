#!/bin/bash
# Quick start script for the logstash experimentation

set -e

echo "🚀 Logstash/Vector Experimentation Setup"
echo "========================================"
echo ""

# Check if Vector is installed
if ! command -v vector &> /dev/null; then
    echo "❌ Vector is not installed"
    echo "📦 Install with: brew install vector"
    echo "   Or visit: https://vector.dev/docs/setup/installation/"
    exit 1
fi

echo "✓ Vector found: $(vector --version | head -n1)"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    echo "📦 Install with: brew install go"
    exit 1
fi

echo "✓ Go found: $(go version)"
echo ""

# Create logs and bin directories
echo "📁 Creating required directories..."
mkdir -p logs bin

# Download Go dependencies
echo "📥 Downloading Go dependencies..."
go mod download
go mod tidy

# Build binaries
echo "🔨 Building binaries..."
go build -o bin/sender ./cmd/sender
go build -o bin/sender-advanced ./cmd/sender-advanced

echo ""
echo "✅ Setup complete!"
echo ""
echo "📖 Quick Start Guide:"
echo ""
echo "1. Start Vector in one terminal:"
echo "   $ vector --config configs/vector.yaml"
echo ""
echo "2. Run the sender in another terminal:"
echo "   $ ./bin/sender"
echo "   $ ./bin/sender-advanced"
echo ""
echo "3. Or use Make commands:"
echo "   $ make vector-start"
echo "   $ make run"
echo ""
echo "4. Or use docker-compose:"
echo "   $ docker-compose up"
echo ""
echo "For more information, see README.md"
