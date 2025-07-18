#!/bin/bash

# Development stack starter script
# This script starts both API and frontend services with proper cleanup

# Function to handle cleanup
cleanup() {
    echo ""
    echo "🛑 Stopping all services..."
    if [ ! -z "$API_PID" ]; then
        kill $API_PID 2>/dev/null || true
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    sleep 2
    go run . cleanup 2>/dev/null || true
    echo "✅ All services stopped"
    exit 0
}

# Set up signal handlers
trap cleanup INT TERM

# Start API server in background
echo "🔧 Starting API server..."
go run . api &
API_PID=$!

# Give API server time to start
sleep 3

# Start frontend in background
echo "🌐 Starting frontend..."
(cd packages/webapp && bun run dev) &
FRONTEND_PID=$!

# Wait for both processes
echo "✅ Both services are running!"
echo "   API: http://localhost:8081"
echo "   Frontend: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop all services gracefully"

# Wait for either process to exit or signal
wait