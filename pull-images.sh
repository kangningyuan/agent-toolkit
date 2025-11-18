#!/bin/bash

echo "==========================================="
echo "Checking Docker availability..."
echo "==========================================="

# Check if Docker is installed and running
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed. Please install Docker first."
    exit 1
fi

if ! docker info &> /dev/null; then
    echo "Error: Docker daemon is not running. Please start Docker service."
    exit 1
fi

echo "Docker is available and running."

echo "==========================================="
echo "Pulling required Docker images for secure execution..."
echo "==========================================="

# Pull Alpine Linux for bash commands
echo "Pulling alpine:latest..."
docker pull alpine:latest

# Pull Python for code execution  
echo "Pulling python:3-alpine..."
docker pull python:3-alpine

echo "==========================================="
echo "Docker images ready for secure execution environment."
echo "Available images:"
docker images | grep -E "(alpine|python)" || echo "No matching images found"
echo "==========================================="