#!/bin/bash

APP_NAME="reminder-cli"

# Clean previous builds
rm -rf dist
mkdir dist

# Linux
GOOS=linux GOARCH=amd64 go build -o dist/${APP_NAME}-linux-amd64
GOOS=linux GOARCH=arm64 go build -o dist/${APP_NAME}-linux-arm64

# macOS
GOOS=darwin GOARCH=amd64 go build -o dist/${APP_NAME}-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o dist/${APP_NAME}-darwin-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o dist/${APP_NAME}-windows-amd64.exe
GOOS=windows GOARCH=arm64 go build -o dist/${APP_NAME}-windows-arm64.exe
