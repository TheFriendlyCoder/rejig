#!/usr/bin/env sh
# Compile a release for the current platform
goreleaser build --single-target --clean --snapshot
# Run all tests recursively
go test ./...