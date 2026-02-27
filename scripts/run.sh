#!/bin/bash
# Run the complete experiment: start Vector and send logs

echo "🚀 Starting Logstash/Vector Experiment"
echo "======================================"
echo ""

# Check if Vector is running
if lsof -Pi :5044 -sTCP:LISTEN -t >/dev/null ; then
    echo "✓ Vector is already running on port 5044"
    VECTOR_RUNNING=true
else
    echo "▶️  Starting Vector..."
    vector --config configs/vector.yaml &
    VECTOR_PID=$!
    VECTOR_RUNNING=false
    
    # Wait for Vector to start
    echo "⏳ Waiting for Vector to be ready..."
    sleep 3
    
    if ! kill -0 $VECTOR_PID 2>/dev/null; then
        echo "❌ Vector failed to start"
        exit 1
    fi
    
    echo "✓ Vector started (PID: $VECTOR_PID)"
fi

echo ""
echo "📤 Sending logs..."
echo ""

# Build if not exists
if [ ! -f "bin/sender" ]; then
    echo "🔨 Building sender..."
    go build -o bin/sender ./cmd/sender
fi

# Run the sender
./bin/sender "$@"

# Clean up if we started Vector
if [ "$VECTOR_RUNNING" = false ]; then
    echo ""
    echo "🛑 Stopping Vector..."
    kill $VECTOR_PID
    wait $VECTOR_PID 2>/dev/null
fi

echo ""
echo "✅ Experiment complete!"
echo ""
echo "📊 Check results in:"
echo "   - logs/vector-output-*.log"
