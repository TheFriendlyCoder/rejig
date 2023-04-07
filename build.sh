#!/usr/bin/env bash -e
# Compile a release for the current platform
#goreleaser build --single-target --clean --snapshot

# Run linter
golangci-lint run

# Run tests
./test.sh