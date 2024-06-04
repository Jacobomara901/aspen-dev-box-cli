#!/bin/bash

# Exit if any command fails
set -e

# Get the ASPEN_DOCKER environment variable
ASPEN_DOCKER=${ASPEN_DOCKER}

# Check if ASPEN_DOCKER is set
if [[ -z "$ASPEN_DOCKER" ]]; then
    echo "Error: ASPEN_DOCKER environment variable not set."
    exit 1
fi

# Create the bin directory if it doesn't exist
mkdir -p bin/linux
mkdir -p bin/darwin

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/linux/adb

# Build for Darwin
GOOS=darwin GOARCH=amd64 go build -o bin/darwin/adb

# Copy the bin directory to ASPEN_DOCKER
cp -r bin "$ASPEN_DOCKER"