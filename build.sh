#!/usr/bin/env bash
set -e

# Compile a release for the current platform
#goreleaser build --single-target --clean --snapshot

# Run tests
./test.sh "$@"

# Run linter
golangci-lint run

# Generate docs
# hugo server -D -s docs
# Clean hugo cache - fixes build problems with themes
# hugo mod clean
