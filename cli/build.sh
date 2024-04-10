#!/bin/bash

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o binaries/enbuild-linux main.go

# Build for Mac
echo "Building for Mac..."
GOOS=darwin GOARCH=amd64 go build -o binaries/enbuild-mac main.go

# Build for Mac M1
echo "Building for Mac M1..."
GOOS=darwin GOARCH=arm64 go build -o binaries/enbuild main.go

echo "Build complete!"
