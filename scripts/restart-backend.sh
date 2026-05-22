#!/bin/bash

# Backend Restart Script
# Stops and then starts the backend services

echo "🔄 Restarting vmconnect backend..."
cd "$(dirname "$0")/.."

docker compose -f build/docker-compose.yml down
sleep 2
docker compose -f build/docker-compose.yml up -d --build

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Backend restarted successfully!"
    echo "📍 API running at: http://localhost:8080"
else
    echo "❌ Failed to restart backend"
    exit 1
fi
