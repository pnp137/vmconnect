#!/bin/bash

# Backend Stop Script for macOS/Linux
# Usage: ./scripts/stop-backend.sh

echo "🛑 Stopping VMConnect Backend..."
echo ""

# Check if containers are running
if ! docker ps | grep -q vmconnect; then
    echo "⚠️  Backend is not running"
    exit 0
fi

echo "🛑 Stopping vmconnect backend..."
cd "$(dirname "$0")/.."
docker compose -f build/docker-compose.yml down

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Backend stopped successfully!"
    echo ""
    echo "📌 Database data is preserved for next time."
    echo "   To start again: ./scripts/start-backend.sh"
else
    echo "❌ Failed to stop backend"
fi
