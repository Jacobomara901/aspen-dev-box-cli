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

# Only build macOS binaries if we're on macOS
if [[ "$(uname)" == "Darwin" ]]; then
    echo "Building macOS binaries..."
    # Build for Intel Macs (amd64)
    GOOS=darwin GOARCH=amd64 go build -o bin/darwin/adb-amd64
    # Build for Apple Silicon (arm64)
    GOOS=darwin GOARCH=arm64 go build -o bin/darwin/adb-arm64
    # Create universal binary
    lipo -create -output bin/darwin/adb bin/darwin/adb-amd64 bin/darwin/adb-arm64
    # Clean up intermediate files
    rm bin/darwin/adb-amd64 bin/darwin/adb-arm64
else
    echo "Not on macOS, skipping macOS binary builds as we cannot make them universal on linux. The actions runner will do this on push to main."
fi

# Copy the bin directory to ASPEN_DOCKER
cp -r bin "$ASPEN_DOCKER"