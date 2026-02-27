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

# Create logs directory
echo "📁 Creating logs directory..."
mkdir -p logs

# Download Go dependencies
echo "📥 Downloading Go dependencies..."
go mod download
go mod tidy

echo ""
echo "✅ Setup complete!"
echo ""
echo "📖 Quick Start Guide:"
echo ""
echo "1. Start Vector in one terminal:"
echo "   $ vector --config vector.yaml"
echo ""
echo "2. Run the sender in another terminal:"
echo "   $ go run sender.go"
echo ""
echo "3. Or use the convenience script:"
echo "   $ ./run.sh"
echo ""
echo "For more information, see README.md"
