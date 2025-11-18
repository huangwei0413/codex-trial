#!/bin/bash

# Test script for student API

set -e

echo "Running unit tests..."
go test -v -race -coverprofile=coverage.out ./test/unit/...

echo "Running integration tests..."
go test -v -race -tags=integration ./test/integration/...

echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo "Test coverage:"
go tool cover -func=coverage.out | tail -1

echo "All tests completed successfully!"
