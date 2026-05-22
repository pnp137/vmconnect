#!/bin/bash

# Backend Startup Script for macOS/Linux
# Usage: ./scripts/start-backend.sh

echo "🚀 Starting VMConnect Backend..."
echo ""

# Check if Docker is running
if ! docker ps &> /dev/null; then
    echo "❌ Docker is not running!"
    echo "📦 Please start Docker Desktop and try again."
    exit 1
fi

echo "✅ Docker is running"
echo ""

# Check if already running
if docker ps | grep -q vmconnect-db; then
    echo "⚠️  Backend is already running!"
    echo ""
    echo "API:      http://localhost:8080"
    echo "Database: localhost:3306"
    echo ""
    echo "To stop:  ./scripts/stop-backend.sh"
    exit 0
fi

echo "🚀 Starting vmconnect backend..."
cd "$(dirname "$0")/.."
docker compose -f build/docker-compose.yml up -d --build

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Backend is starting..."
    echo ""
    echo "⏳ Waiting for MySQL to be ready..."
    sleep 10
    
    # Test health endpoint
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Backend is ready!"
        echo ""
        echo "🌐 API URL:      http://localhost:8080"
        echo "🗄️  Database:     localhost:3306"
        echo "👤 DB User:      vmconnect_user"
        echo "🔑 DB Password:  vmconnect_pass"
        echo ""
        echo "📝 Next: Open your frontend and use http://localhost:8080 as the API base URL"
        echo ""
        echo "To view logs:    docker logs vmconnect-api -f"
        echo "To stop:         ./scripts/stop-backend.sh"
    else
        echo "⏳ Still starting... (this can take a minute)"
        echo ""
        echo "Check status with:"
        echo "  docker logs vmconnect-api"
        echo "  docker logs vmconnect-db"
    fi
else
    echo "❌ Failed to start backend"
    echo "Check Docker status: docker ps -a"
fi
