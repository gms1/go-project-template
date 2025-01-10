#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

echo "lint..."

set -e
set -o pipefail
cd "${DN}/.."

go mod tidy
golangci-lint run

echo "lint: SUCCEEDED"
