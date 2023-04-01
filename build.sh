#!/usr/bin/env sh -e
# Compile a release for the current platform
#goreleaser build --single-target --clean --snapshot
golangci-lint run -p bugs -p error
# Run all tests recursively
#go test ./...

# With coverage
#go test --cover ./...

# With HTML report
go test -coverprofile cover.out ./...
go tool cover -html=cover.out