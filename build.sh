#!/usr/bin/env sh -e
# Compile a release for the current platform
#goreleaser build --single-target --clean --snapshot
golangci-lint run -p bugs -p error
# Run all tests recursively
#go test ./...
go test -coverprofile cover.out ./...
go tool cover -html=cover.out