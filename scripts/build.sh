#!/usr/bin/env bash
#
# This script builds the application from source for multiple platforms.
set -e

echo "==> Building..."
GIT_COMMIT="$(git rev-parse HEAD)"
go build -o ./bin/semaphore -ldflags "-X main.version=unreleased -X main.build=${GIT_COMMIT} -X main.label=development" ./cmd/semaphore

echo "==> Results:"
ls -hl bin/