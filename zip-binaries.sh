#!/bin/bash
cd dist || exit 1

zip reminder-cli-linux-amd64.zip reminder-cli-linux-amd64
zip reminder-cli-linux-arm64.zip reminder-cli-linux-arm64
zip reminder-cli-darwin-amd64.zip reminder-cli-darwin-amd64
zip reminder-cli-darwin-arm64.zip reminder-cli-darwin-arm64
zip reminder-cli-windows-amd64.zip reminder-cli-windows-amd64.exe
zip reminder-cli-windows-arm64.zip reminder-cli-windows-arm64.exe
