#!/bin/bash

APP_NAME="reminder-cli"

# Clean previous builds
rm -rf dist
mkdir dist

GOOS=linux GOARCH=amd64 go build -o dist/${APP_NAME}-linux-amd64 ./cmd/root.go
GOOS=linux GOARCH=arm64 go build -o dist/${APP_NAME}-linux-arm64 ./cmd/root.go

GOOS=darwin GOARCH=amd64 go build -o dist/${APP_NAME}-darwin-amd64 ./cmd/root.go
GOOS=darwin GOARCH=arm64 go build -o dist/${APP_NAME}-darwin-arm64 ./cmd/root.go

GOOS=windows GOARCH=amd64 go build -o dist/${APP_NAME}-windows-amd64.exe ./cmd/root.go
GOOS=windows GOARCH=arm64 go build -o dist/${APP_NAME}-windows-arm64.exe ./cmd/root.go
