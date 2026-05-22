#!/bin/bash

# View Backend Logs Script
# Usage: ./scripts/logs.sh

if [ "$1" == "api" ]; then
    echo "📊 Backend API Logs (live updates)..."
    docker logs vmconnect-api -f
elif [ "$1" == "db" ]; then
    echo "📊 Database Logs (live updates)..."
    docker logs vmconnect-db -f
elif [ "$1" == "status" ]; then
    echo "📊 Backend Status..."
    docker ps | grep vmconnect
else
    echo "Backend Logs Helper"
    echo ""
    echo "Usage: ./scripts/logs.sh [option]"
    echo ""
    echo "Options:"
    echo "  api      - View API logs (live)"
    echo "  db       - View Database logs (live)"
    echo "  status   - Show running containers"
    echo ""
    echo "Examples:"
    echo "  ./scripts/logs.sh api"
    echo "  ./scripts/logs.sh status"
fi
