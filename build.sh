#!/bin/sh
set -e

go mod tidy

LDFLAGS='-s -w'

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o nocat_windows_amd64.exe .
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="$LDFLAGS" -o nocat_windows_arm64.exe .

# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -o nocat_darwin_amd64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -o nocat_darwin_arm64 .

# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o nocat_linux_amd64 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$LDFLAGS" -o nocat_linux_arm64 .

echo "Build completed."