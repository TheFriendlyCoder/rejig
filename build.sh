#!/usr/bin/env sh -e
# Compile a release for the current platform
#goreleaser build --single-target --clean --snapshot

# Run tests
./test.sh

# Run linter
golangci-lint run